package tun

import (
	"io"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/ahmetozer/wstransit/pkg/conn"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/websocket"
)

var (
	networkInterface map[string]*water.Interface
	wsInterface      map[string]*websocket.Conn
	ConnectionStat   map[string]Status
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
}

func HandleClient(wsconn *websocket.Conn) {
	queries := wsconn.Request().URL.Query()

	cliconn := conn.GetConnFromParameter(queries)

	var err error
	if l := len(cliconn.Id); l != 12 {
		log.Println("Client id is not valid", l)
		wsconn.Write([]byte("error: connection lenght is not equal to 12"))
		wsconn.Close()
		return
	}

	log.Printf("new tun connection %s %s", cliconn.Id, ConnectionStat[cliconn.Id].String())
	if ConnectionStat[cliconn.Id] == ClientNotCreated {
		config := water.Config{
			DeviceType: water.TUN,
		}

		switch runtime.GOOS {
		case "linux":
			config.PlatformSpecificParams = water.PlatformSpecificParams{
				Name: cliconn.Id,
			}

		default:
			log.Printf("unsupported platform: %s", runtime.GOOS)
			wsconn.Write([]byte("server does not support tun on this operation system"))
			wsconn.Close()
			return
		}

		networkInterface[cliconn.Id], err = water.New(config)
		if err != nil {
			log.Printf("tun create error: %s %s", cliconn.Id, err)
			wsconn.Close()
			return
		}

		ConnectionStat[cliconn.Id] = ClientCreated

		link, err := netlink.LinkByName(cliconn.Id)
		if err != nil {
			log.Printf("netlink select error: %s %s", cliconn.Id, err)
			wsconn.Close()
			return
		}

		netlink.LinkSetUp(link)
		netlink.AddrAdd(link, &netlink.Addr{IPNet: &net.IPNet{IP: net.ParseIP("10.0.3.1"), Mask: net.CIDRMask(24, 32)}})

		defer func() {
			log.Printf("defer executed for %s", cliconn.Id)
			for ConnectionStat[cliconn.Id] == ClientCreated || ConnectionStat[cliconn.Id] == ClientCreatedAndUsed || ConnectionStat[cliconn.Id] == ClientCreatedAndUsedButDied {
				log.Printf("defer id %s stat %s", cliconn.Id, ConnectionStat[cliconn.Id])

				time.Sleep(30 * time.Second)
				if ConnectionStat[cliconn.Id] == ClientNotCreated {
					log.Printf("defer id %s stat %s break", cliconn.Id, ConnectionStat[cliconn.Id])
					break
				}
				// If the connection drop end of the sleep, we need a another lap to ensure at least 15 second is waited
				if ConnectionStat[cliconn.Id] == ClientCreatedAndUsedButDied {
					ConnectionStat[cliconn.Id] = ClientNotCreated
				}
			}

			ConnectionStat[cliconn.Id] = ClientCreated // connection closed
			log.Println("tun connection closed for", cliconn.Id)
			networkInterface[cliconn.Id].Close()
			wsconn.Close()

			log.Printf("link is not delted %s %s", cliconn.Id, err)
			netlink.LinkDel(link)
			wsInterface[cliconn.Id].Close()
			delete(networkInterface, cliconn.Id)
			delete(ConnectionStat, cliconn.Id)
			delete(wsInterface, cliconn.Id)

		}()

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
