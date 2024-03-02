package main

import (
	"embed"
	"errors"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"roomate/db"
	"roomate/handler"
	"roomate/pkg/sb"
	"roomate/pkg/ws"
)

//go:embed public
var FS embed.FS
var ErrUserNotFound = errors.New("user not found")

func foo() error {
	return ErrUserNotFound
}

func main() {
	err := initEverything()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	err = foo()
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			println("si, user not found")
		}
	}

	// 	createTableSQL := `
	// CREATE TABLE messages (
	//     id   BIGSERIAL PRIMARY KEY,
	//     from_id  BIGINT REFERENCES users(id) NOT NULL,
	//     from_authid  string NOT NULL,
	//     from_name text NOT NULL,
	//     message text NOT NULL,
	//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	// );`
	//
	// 	_, err = db.DB.Exec(createTableSQL)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}

	e := echo.New()
	e.Use(handler.WithUser)

	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(FS)))))
	e.GET("/", handler.Make(func(c echo.Context) error {
		user := handler.GetAuthenticatedUser(c)
		if !user.IsLoggedIn {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		return c.Redirect(http.StatusSeeOther, "/chat")
	}))
	e.GET("/login", handler.Make(handler.LoginShow))
	e.POST("/login", handler.Make(handler.Login))
	e.GET("/register", handler.Make(handler.RegisterShow))
	e.POST("/register", handler.Make(handler.Register))
	e.GET("/login/callback", handler.Make(handler.RegisterCallback))
	e.GET("/chat", handler.Make(handler.ChatShow))
	e.GET("/chatws", handler.Make(handler.ChatWS))
	e.GET("/stream", handler.Make(handler.StreamShow))
	e.GET("/streamws", handler.Make(handler.StreamWS))

	e.Logger.Fatal(e.Start(":3000"))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	ws.Init()
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
