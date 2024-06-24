//go:build linux

package tun

import (
	"log"
	"time"

	"github.com/ahmetozer/bosphorus/pkg/conn"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/websocket"
)

func clean(cliconn conn.ConnectionString, wsconn *websocket.Conn, link netlink.Link) {
	for ConnectionStat[cliconn.Id] == ClientCreated || ConnectionStat[cliconn.Id] == ClientCreatedAndUsed || ConnectionStat[cliconn.Id] == ClientCreatedAndUsedButDied {
		time.Sleep(30 * time.Second)
		if ConnectionStat[cliconn.Id] == ClientNotCreated {
			break
		}
		// If the connection drop end of the sleep, we need a another lap to ensure at least 15 second is waited
		if ConnectionStat[cliconn.Id] == ClientCreatedAndUsedButDied {
			ConnectionStat[cliconn.Id] = ClientNotCreated
		}
	}

	ConnectionStat[cliconn.Id] = ClientCreated // connection closed
	networkInterface[cliconn.Id].Close()
	wsconn.Close()

	err := netlink.LinkDel(link)
	if err != nil {
		log.Printf("link is not delted %s %s", cliconn.Id, err)
	}

	wsInterface[cliconn.Id].Close()
	delete(networkInterface, cliconn.Id)
	delete(ConnectionStat, cliconn.Id)
	delete(wsInterface, cliconn.Id)

}
