package tun

import (
	"encoding/binary"
	"errors"
	"net"
)

type dhcp struct {
	ipv4 net.IP
	ipv6 net.IP
}

func dhcpRelease(clientNumber uint32) (dhcp, error) {

	b := make([]byte, 4)
	copy(b, tunIpv4Net.IP[:])

	var release dhcp
	if clientNumber == 0 {
		return dhcp{}, errors.New("there is no empty ipv4 address")
	}

	tt := binary.BigEndian.Uint32(b)
	binary.BigEndian.PutUint32(b, tt+clientNumber)

	if tunIpv4Net.Contains(b) {
		release.ipv4 = b
	} else {
		return dhcp{}, errors.New("created ip is not cidr range, out of ipv4")
	}

	return release, nil

}

func findSpaceClient() uint32 {

	var i uint32 = 1
	for ; i < uint32(len(client)); i++ {
		if !client[i] {
			return i
		}
	}
	return 0
}
