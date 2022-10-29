package client

import (
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/websocket"
)

var (
	Connection     map[string]*net.TCPConn
	ConnectionStat map[string]Status
)

type Status uint8

const (
	ClientNotCreated Status = iota
	ClientCreated
	ClientCreatedAndUsed
	ClientCreatedAndUsedButDied
)

func (d Status) String() string {
	return [...]string{"ClientNotCreated", "ClientCreated", "ClientCreatedAndUsed", "ClientCreatedAndUsedButDied"}[d]
}

func init() {
	Connection = make(map[string]*net.TCPConn)
	ConnectionStat = make(map[string]Status)
}

func HandleClient(wsconn *websocket.Conn) {
	queries := wsconn.Request().URL.Query()

	remoteAddr := queries.Get("remoteAddr")
	connId := queries.Get("connId")
	if len(connId) != 25 {
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		log.Println("ResolveTCPAddr failed:", err.Error())
		return
	}

	if ConnectionStat[connId] == ClientNotCreated {
		log.Println("new connection is creating for", connId, remoteAddr)
		Connection[connId], err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Println("tcp dial failed:", err.Error())
			return
		}
		Connection[connId].SetKeepAlive(true)
		Connection[connId].SetKeepAlivePeriod(time.Second * 5)

		ConnectionStat[connId] = ClientCreated
		log.Println("tcp connected for", connId, remoteAddr)

		defer func(string) {

			for ConnectionStat[connId] == ClientCreated || ConnectionStat[connId] == ClientCreatedAndUsed || ConnectionStat[connId] == ClientCreatedAndUsedButDied {
				time.Sleep(15 * time.Second)
				if ConnectionStat[connId] == ClientCreated {
					break
				}
				// If the connection drop end of the sleep, we need a another lap to ensure at least 15 second is waited
				if ConnectionStat[connId] == ClientCreatedAndUsedButDied {
					ConnectionStat[connId] = ClientNotCreated
				}
			}

			ConnectionStat[connId] = ClientCreated // connection closed
			log.Println("tcp connection closed for", connId, remoteAddr)
			Connection[connId].Close()

			delete(Connection, connId)
			delete(ConnectionStat, connId)

		}(connId)

	} else {
		log.Println("tcp connection reused for", connId, remoteAddr)
	}

	ConnectionStat[connId] = ClientCreatedAndUsed

	// tcp client to web socket data flow
	go io.Copy(Connection[connId], wsconn)

	io.Copy(wsconn, Connection[connId])

	log.Printf("websocket for tcp tunnel end %s %s\n", connId, Connection[connId].RemoteAddr())

	// Connection created and used in the past but died
	ConnectionStat[connId] = ClientCreatedAndUsedButDied
}
