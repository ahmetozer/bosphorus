package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/net/websocket"
)

func tcpStart(a arrFlag, wg *sync.WaitGroup) {

	for _, flag := range a {
		connectionString, err := parseConnectionString(flag)
		if err != nil {
			log.Fatalf("tcp flag parse: %s, %s ", err, flag)
		}
		go tcpListener(connectionString, wg)
		wg.Add(1)
	}
}

func tcpListener(c connectionString, wg *sync.WaitGroup) {

	listen, err := net.Listen("tcp", c.localAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// close listener
	defer listen.Close()
	defer wg.Done()
	log.Printf("tcp: %s <=> %s <=> %s", c.localAddr, c.url, c.remmoteAddr)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("unable to accept tcp connection for '%s' :%s", c.localAddr, err)
		}

		go handleTCPrequest(conn, c)
	}
}

func handleTCPrequest(tcpConn net.Conn, c connectionString) {
	defer tcpConn.Close()
	log.Printf("new connection: %s", tcpConn.RemoteAddr())

	remote := newTCPURL(c)
	wsConfig, err := websocket.NewConfig(remote, remote)
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	rwc, err := newWSSocket(wsConfig)
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

func newTCPURL(c connectionString) string {
	return fmt.Sprintf("%s?connType=tcp&remoteAddr=%s", c.url, c.remmoteAddr)
}
