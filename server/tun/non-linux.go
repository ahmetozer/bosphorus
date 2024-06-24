//go:build !linux

package tun

import "golang.org/x/net/websocket"

func HandleClient(wsconn *websocket.Conn) {
	wsconn.Write([]byte("tun is not suppored at server side"))
	wsconn.Close()
}
