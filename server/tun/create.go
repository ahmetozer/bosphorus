package tun

import (
	"errors"
	"net"
	"runtime"

	"github.com/ahmetozer/wstransit/pkg/conn"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/websocket"
)

func create(cliconn conn.ConnectionString, wsconn *websocket.Conn) (netlink.Link, error) {
	var err error
	config := water.Config{
		DeviceType: water.TUN,
	}

	switch runtime.GOOS {
	case "linux":
		config.PlatformSpecificParams = water.PlatformSpecificParams{
			Name: cliconn.Id,
		}

	default:
		wsconn.Write([]byte("server does not support tun on this operation system"))
		wsconn.Close()
		return &netlink.GenericLink{}, errors.New("unsupported platform")
	}

	networkInterface[cliconn.Id], err = water.New(config)
	if err != nil {
		wsconn.Close()
		return &netlink.GenericLink{}, errors.New("tun create error" + err.Error())
	}

	ConnectionStat[cliconn.Id] = ClientCreated

	link, err := netlink.LinkByName(cliconn.Id)
	if err != nil {
		wsconn.Close()
		return &netlink.GenericLink{}, errors.New("netlink select error" + err.Error())
	}

	netlink.LinkSetUp(link)
	netlink.AddrAdd(link, &netlink.Addr{IPNet: &net.IPNet{IP: net.ParseIP("10.0.3.1"), Mask: net.CIDRMask(24, 32)}})

	return link, nil

}
