package client

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/ahmetozer/bosphorus/pkg/conn"
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

	cliconn := conn.GetConnFromParameter(queries)

	if len(cliconn.Id) != 12 {
		log.Println("lenght is not 12:", cliconn)
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", cliconn.RemmoteAddr)
	if err != nil {
		log.Println("ResolveTCPAddr failed:", err.Error())
		return
	}

	if ConnectionStat[cliconn.Id] == ClientNotCreated {
		log.Println("new connection is creating for", cliconn.Id, cliconn.RemmoteAddr)
		Connection[cliconn.Id], err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Println("tcp dial failed:", err.Error())
			wsconn.Close()
			return
		}
		Connection[cliconn.Id].SetKeepAlive(true)
		Connection[cliconn.Id].SetKeepAlivePeriod(time.Second * 5)

		ConnectionStat[cliconn.Id] = ClientCreated
		log.Println("tcp connected for", cliconn.Id, Connection[cliconn.Id].LocalAddr().String(), cliconn.RemmoteAddr)

		defer func(string) {

			for ConnectionStat[cliconn.Id] == ClientCreated || ConnectionStat[cliconn.Id] == ClientCreatedAndUsed || ConnectionStat[cliconn.Id] == ClientCreatedAndUsedButDied {
				time.Sleep(15 * time.Second)
				if ConnectionStat[cliconn.Id] == ClientCreated {
					break
				}
				// If the connection drop end of the sleep, we need a another lap to ensure at least 15 second is waited
				if ConnectionStat[cliconn.Id] == ClientCreatedAndUsedButDied {
					ConnectionStat[cliconn.Id] = ClientNotCreated
				}
			}

			ConnectionStat[cliconn.Id] = ClientCreated // connection closed
			log.Println("tcp connection closed for", cliconn.Id, Connection[cliconn.Id].LocalAddr().String(), cliconn.RemmoteAddr)
			Connection[cliconn.Id].Close()

			delete(Connection, cliconn.Id)
			delete(ConnectionStat, cliconn.Id)

		}(cliconn.Id)

	} else {
		log.Println("tcp connection reused for", cliconn.Id, Connection[cliconn.Id].LocalAddr().String(), cliconn.RemmoteAddr)
	}

	ConnectionStat[cliconn.Id] = ClientCreatedAndUsed

	// tcp client to web socket data flow
	go io.Copy(Connection[cliconn.Id], wsconn)

	io.Copy(wsconn, Connection[cliconn.Id])

	log.Printf("websocket for tcp tunnel end %s %s %s\n", cliconn.Id, Connection[cliconn.Id].LocalAddr().String(), cliconn.RemmoteAddr)

	// Connection created and used in the past but died
	ConnectionStat[cliconn.Id] = ClientCreatedAndUsedButDied
}
