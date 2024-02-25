package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"

	"roomate/db"
	"roomate/pkg/ws"
	"roomate/sqlc"
	"roomate/types"
	"roomate/view/chat"
)

var wsMutex sync.Mutex

func ChatShow(c echo.Context) error {
	user := GetAuthenticatedUser(c)
	if !user.IsLoggedIn {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	messages, err := db.Queries.GetMessages(context.Background())
	if err != nil {
		return err
	}
	return render(c, chat.Chat(messages, user.AuthID))
}

func ChatWS(c echo.Context) error {
	user := GetAuthenticatedUser(c)
	if !user.IsLoggedIn {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	userDB, err := db.Queries.GetUserByAuthID(context.Background(), user.AuthID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	websocket.Handler(func(w *websocket.Conn) {
		wsMutex.Lock()
		ws.WSConnections[w] = user
		wsMutex.Unlock()

		defer w.Close()
		for {
			msg := ""
			err := websocket.Message.Receive(w, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			var jsonMap map[string]interface{}
			json.Unmarshal([]byte(msg), &jsonMap)
			message, ok := jsonMap["chat_message"].(string)
			if !ok {
				return
			}
			db.Queries.CreateMessage(context.Background(), sqlc.CreateMessageParams{
				FromID:     userDB.ID,
				FromAuthid: userDB.Authid,
				FromName:   userDB.Name,
				Message:    message,
			})
			broadcast(c, userDB, ws.WSConnections, message)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func broadcast(
	c echo.Context,
	from sqlc.User,
	to types.Conns,
	message string,
) {
	freshTextInput := chat.Input("")
	bufferTextInput := &bytes.Buffer{}
	freshTextInput.Render(context.Background(), bufferTextInput)

	for conn := range to {
		if to[conn].AuthID == from.Authid {
			// e' lui'
			c := chat.Message(from.Name, "10:20", message, true)
			buffer := &bytes.Buffer{}
			c.Render(context.Background(), buffer)
			websocket.Message.Send(
				conn,
				buffer.String()+bufferTextInput.String(),
			)
		} else {
			// e' un altro
			c := chat.Message(from.Name, "10:20", message, false)
			buffer := &bytes.Buffer{}
			c.Render(context.Background(), buffer)
			websocket.Message.Send(
				conn,
				buffer.String()+bufferTextInput.String(),
			)
		}
	}
}
