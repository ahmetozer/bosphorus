package tun

import (
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ahmetozer/bosphorus/pkg/conn"
	"github.com/songgao/water"
	"golang.org/x/net/websocket"
)

var (
	networkInterface map[string]*water.Interface
	wsInterface      map[string]*websocket.Conn
	ConnectionStat   map[string]Status
	RouteTableNumber int = 300
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
	networkInterface = make(map[string]*water.Interface)
	wsInterface = make(map[string]*websocket.Conn)
	ConnectionStat = make(map[string]Status)

	routeTableNumber := os.Getenv("ROUTE_TABLE_NUMBER")
	if routeTableNumber != "" {
		var err error
		RouteTableNumber, err = strconv.Atoi(routeTableNumber)
		if err != nil {
			log.Fatalf("TUN_ROUTE_TABLE_NUMBER is not a number: %s", err)
		}
	}
}

func HandleClient(wsconn *websocket.Conn) {
	queries := wsconn.Request().URL.Query()

	cliconn := conn.GetConnFromParameter(queries)

	if l := len(cliconn.Id); l != 12 {
		log.Println("Client id is not valid", l)
		wsconn.Write([]byte("error: connection lenght is not equal to 12"))
		wsconn.Close()
		return
	}

	log.Printf("new tun connection %s %s", cliconn.Id, ConnectionStat[cliconn.Id].String())

	if ConnectionStat[cliconn.Id] == ClientNotCreated {
		link, err := create(cliconn, wsconn)
		defer clean(cliconn, wsconn, link)

		if err != nil {
			log.Println("error while creating tun interface", err)
			return
		}

	} else {
		log.Println("tun connection reused for", cliconn.Id)
		wsInterface[cliconn.Id].Close()
	}

	ConnectionStat[cliconn.Id] = ClientCreatedAndUsed
	wsInterface[cliconn.Id] = wsconn

	log.Printf("websocket for tun tunnel start %s \n", cliconn.Id)

	go io.Copy(networkInterface[cliconn.Id], wsconn)
	io.Copy(wsconn, networkInterface[cliconn.Id])

	log.Printf("websocket for tun tunnel end %s \n", cliconn.Id)

	// Connection created and used in the past but died
	ConnectionStat[cliconn.Id] = ClientCreatedAndUsedButDied
}
