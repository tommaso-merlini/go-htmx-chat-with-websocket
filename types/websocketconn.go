package types

import (
	"golang.org/x/net/websocket"
)

type Conns map[*websocket.Conn]AuthenticatedUser
