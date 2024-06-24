package server

import (
	"github.com/ahmetozer/bosphorus/pkg/conn"
	tcpClient "github.com/ahmetozer/bosphorus/server/tcp/client"
	"github.com/ahmetozer/bosphorus/server/tun"
	"golang.org/x/net/websocket"
)

func HandleNewClient(wsconn *websocket.Conn) {

	connType := wsconn.Request().URL.Query().Get("Type")
	switch connType {
	case conn.TCP.String():
		tcpClient.HandleClient(wsconn)
	case conn.TUN.String():
		tun.HandleClient(wsconn)
	default:
		wsconn.Write([]byte("unexpected connType: " + connType))
		wsconn.Close()
	}

}
