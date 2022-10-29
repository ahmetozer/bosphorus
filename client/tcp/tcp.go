package tcp

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/ahmetozer/wstransit/client/flags"
	"github.com/ahmetozer/wstransit/pkg/conn"
)

const (
	TCPTypeListen = iota
	TCPTypeRawListen
)

func TcpStart(a flags.ArrFlag, wg *sync.WaitGroup, ttype int) {

	for _, flag := range a {
		connectionString, err := parseTCPConnectionString(flag)
		if err != nil {
			log.Fatalf("tcp flag parse: %s, %s ", err, flag)
		}
		switch ttype {
		case TCPTypeListen:
			go tcpListener(connectionString, wg)
		// case TCPTypeRawListen:
		// 	go tcpRawListener(connectionString, wg)
		default:
			log.Fatalf("unexpected tcp type %d\n", ttype)
		}
		wg.Add(1)
	}
}

func parseTCPConnectionString(config string) (conn.ConnectionString, error) {

	// Example connection string per argument
	// 1 localAddr					2 host						3 remoteAddr
	// localhost:8119;http://nyc-sv-nss.ahmet.engineer/ws/tcp;localhost:8118

	s := strings.Split(config, ";")

	if len(s) != 3 {
		return conn.ConnectionString{}, errors.New("property does not meet struct")
	}

	return conn.ConnectionString{
		LocalAddr:   s[0],
		Url:         s[1],
		RemmoteAddr: s[2],
		Type:        conn.TCP,
		Id:          conn.GenerateConnID(),
	}, nil

}
