//go:build linux

package tun

import (
	"errors"
	"strings"

	"github.com/ahmetozer/bosphorus/pkg/conn"
)

func parseTUNConnectionString(config string) (conn.ConnectionString, error) {

	// Example connection string per argument
	// 1 networkinterfacename					2 host
	// utun;http://nyc-sv-nss.ahmet.engineer/ws/tcp

	s := strings.Split(config, ";")

	if len(s) != 2 {
		return conn.ConnectionString{}, errors.New("property does not meet struct")
	}

	return conn.ConnectionString{
		LocalInterface: s[0],
		Url:            s[1],
		Id:             conn.GenerateConnID(),
		Type:           conn.TUN,
	}, nil

}
