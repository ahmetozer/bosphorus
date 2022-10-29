package tcp

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/ahmetozer/wstransit/client/ws"
	"github.com/ahmetozer/wstransit/pkg/conn"
	"golang.org/x/net/websocket"
)

func tcpListener(c conn.ConnectionString, wg *sync.WaitGroup) {

	listen, err := net.Listen("tcp", c.LocalAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// close listener
	defer listen.Close()
	defer wg.Done()
	log.Printf("tcp: %s <=> %s <=> %s", c.LocalAddr, c.Url, c.RemmoteAddr)

	remote := conn.NewURL(c)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("unable to accept tcp connection for '%s' :%s", c.LocalAddr, err)
		}

		go handleTCPrequest(conn, c, remote)
	}
}

func handleTCPrequest(tcpConn net.Conn, c conn.ConnectionString, remote string) {
	defer tcpConn.Close()
	log.Printf("new connection: %s", tcpConn.RemoteAddr())

	wsConfig, err := websocket.NewConfig(remote, remote)
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	rwc, err := ws.NewWSSocket(wsConfig)
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	wsClient, err := websocket.NewClient(wsConfig, rwc)
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	go io.Copy(wsClient, tcpConn)

	io.Copy(tcpConn, wsClient)
}
