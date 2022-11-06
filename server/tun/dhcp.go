package tun

import (
	"errors"
	"net"
)

type dhcp struct {
	ipv4 net.IP
	ipv6 net.IP
}

func dhcpRelease(clientNumber uint32) (dhcp, error) {

	var release dhcp
	if clientNumber == 0 {
		return dhcp{}, errors.New("there is no empty ip addresses")
	}

	ipv4 := IpUint64Add(tunIpv4Net.IP.To16(), uint64(clientNumber))

	if tunIpv4Net.Contains(ipv4) {
		release.ipv4 = ipv4.To4()
	} else {
		return dhcp{}, errors.New("created ip is not cidr range, out of ipv4")
	}

	ipv6 := IpUint64Add(tunIpv6Net.IP, uint64(clientNumber))
	if tunIpv4Net.Contains(ipv6) {
		release.ipv6 = ipv6
		return release, nil
	}

	return dhcp{}, errors.New("created ip is not cidr range, out of ipv6")

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
