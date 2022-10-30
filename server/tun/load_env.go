package tun

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ahmetozer/wstransit/pkg/conn"
)

var (
	tunIPv6Cidr string
	ipv6TunNet  *net.IPNet
	ipv6TunIp   net.IP
)

func init() {

	envIpv6Cidr()
}

func envIpv6Cidr() {
	tunIPv6Cidr := os.Getenv("TUN_IPV6_CIDR")
	var err error
	if tunIPv6Cidr == "" {
		id := conn.GenerateConnID()
		p, err := base64.StdEncoding.DecodeString(id)
		if err != nil {
			log.Printf("err %s %b", id, p)
			log.Fatal(err)
		}
		tunIPv6Cidr = hex.EncodeToString(p)
		// fd01:8429:15ee:67aa::/64
		tunIPv6Cidr = "fd" + tunIPv6Cidr[0:2] + ":" + tunIPv6Cidr[2:6] + ":" + tunIPv6Cidr[6:10] + ":" + tunIPv6Cidr[10:14] + "::/64"
		log.Printf("TUN_IPV6_CIDR is not asigned, auto value used\n")
	}

	release := make([]bool, 1024)
	release[0] = true
	ipv6TunIp, ipv6TunNet, err = net.ParseCIDR(tunIPv6Cidr)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	s := strings.Split(tunIPv6Cidr, "/")
	cidrSize, _ := strconv.Atoi(s[1])
	if cidrSize > 80 {
		log.Fatalf("TUN_IPV6_CIDR '%d' size is low. cidr range must bigger than 80\n", cidrSize)
	}

	log.Printf("TUN_IPV6_CIDR: %s/%s", &ipv6TunNet.IP, s[1])

}
