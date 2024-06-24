//go:build linux

package tun

import (
	"io"
	"log"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/ahmetozer/bosphorus/client/ws"
	"github.com/ahmetozer/bosphorus/pkg/conn"
	"github.com/songgao/water"
	"golang.org/x/net/websocket"
)

var (
	macInterfaceRegex *regexp.Regexp
)

func init() {
	macInterfaceRegex, _ = regexp.Compile("utun[0-9]+")
}

func tunListener(c conn.ConnectionString, wg *sync.WaitGroup) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	if runtime.GOOS == "darwin" {
		if !macInterfaceRegex.MatchString(c.LocalInterface) {
			log.Fatalf("interface name does not meet %s for darwin\n", macInterfaceRegex)
		}
	}
	config.Name = c.LocalInterface

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	remote := conn.NewURL(c)
	wsConfig, err := websocket.NewConfig(remote, remote)
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	for {
		rwc, err := ws.NewWSSocket(wsConfig)
		if err != nil {
			log.Printf("rwc %s", err)
			time.Sleep(time.Second * 10)
			continue
		}

		wsClient, err := websocket.NewClient(wsConfig, rwc)
		if err != nil {
			log.Printf("wsClient %s", err)
			time.Sleep(time.Second * 10)
			continue
		}

		//
		//

		go io.Copy(wsClient, ifce)

		io.Copy(ifce, wsClient)
		log.Printf("re connecting tun %s,%s", c.LocalInterface, c.Url)
		time.Sleep(time.Second * 1)
	}

}
