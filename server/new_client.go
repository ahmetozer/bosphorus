package server

import (
	"golang.org/x/net/websocket"
)

func handleNewClient(wsconn *websocket.Conn) {

	connType := wsconn.Request().URL.Query().Get("connType")
	switch connType {
	case "tcp":
		handleTCPClient(wsconn)
	default:
		wsconn.Write([]byte("unexpected connType: " + connType))
		wsconn.Close()
	}

}
