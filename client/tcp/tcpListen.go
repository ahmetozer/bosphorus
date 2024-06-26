package tcp

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/ahmetozer/bosphorus/client/ws"
	"github.com/ahmetozer/bosphorus/pkg/conn"
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

	for {
		tcpconn, err := listen.Accept()

		if err != nil {
			log.Fatalf("unable to accept tcp connection for '%s' :%s", c.LocalAddr, err)
		}

		go handleTCPrequest(tcpconn, c)
	}
}

func handleTCPrequest(tcpConn net.Conn, c conn.ConnectionString) {
	defer tcpConn.Close()
	log.Printf("new connection: %s", tcpConn.RemoteAddr())
	c.Id = conn.GenerateConnID()
	remote := conn.NewURL(c)
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

	done := make(chan struct{})

	defer wsClient.Close()
	defer tcpConn.Close()
	connActive := true
	go func() {
		defer close(done)
		buf := make([]byte, conn.BufferSize)
		_, err := io.CopyBuffer(wsClient, tcpConn, buf)
		if err != nil && connActive {
			log.Printf("error copying from local to remote: %v", err)
		}
		connActive = false
	}()

	buf := make([]byte, conn.BufferSize)
	_, err = io.CopyBuffer(tcpConn, wsClient, buf)
	if err != nil && connActive {
		log.Printf("error copying from remote to local: %v", err)
	}
	connActive = false

	<-done
}
