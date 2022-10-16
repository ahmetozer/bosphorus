package server

import (
	"io"
	"log"
	"net"

	"golang.org/x/net/websocket"
)

func handleTCPClient(wsconn *websocket.Conn) {
	queries := wsconn.Request().URL.Query()

	remoteAddr := queries.Get("remoteAddr")
	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		return
	}

	tcpconn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		return
	}

	println("connected to :", remoteAddr)
	defer tcpconn.Close()

	go io.Copy(tcpconn, wsconn)
	io.Copy(wsconn, tcpconn)
	log.Printf("conn finished %s <=> %s\n", wsconn.RemoteAddr(), tcpconn.RemoteAddr())
}
