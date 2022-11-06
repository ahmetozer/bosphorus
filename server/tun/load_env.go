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
	client []bool

	tunIpv6Ip net.IP
	tunIpv4Ip net.IP

	tunIPv6Cidr string
	tunIPv4Cidr string

	tunIpv6Net *net.IPNet
	tunIpv4Net *net.IPNet
)

func init() {

	envIpv6Cidr()
	envIpv4Cidr()

}

func envIpv6Cidr() {
	tunIPv6Cidr = os.Getenv("TUN_IPV6_CIDR")
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

	tunIpv6Ip, tunIpv6Net, err = net.ParseCIDR(tunIPv6Cidr)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	s := strings.Split(tunIPv6Cidr, "/")
	cidrSize, err := strconv.Atoi(s[1])
	if err != nil {
		log.Fatalf("Error cidrSize: %s\n", err)
	}
	if cidrSize > 80 {
		log.Fatalf("TUN_IPV6_CIDR '%d' size is low. cidr range must bigger or equal to 80\n", cidrSize)
	}

	log.Printf("TUN_IPV6_CIDR: %s/%s\n", &tunIpv6Net.IP, s[1])

}

func envIpv4Cidr() {
	tunIPv4Cidr = os.Getenv("TUN_IPV4_CIDR")
	var err error
	if tunIPv4Cidr == "" {
		tunIPv4Cidr = "10.90.2.4/24"
		log.Printf("TUN_IPV4_CIDR is not asigned, default value used\n")
	}

	tunIpv4Ip, tunIpv4Net, err = net.ParseCIDR(tunIPv4Cidr)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	s := strings.Split(tunIPv4Cidr, "/")
	cidrSize, _ := strconv.Atoi(s[1])
	if cidrSize >= 32 {
		log.Fatalf("TUN_IPV4_CIDR '%d' size is too low. recommend cidr equal or bigger then 24\n", cidrSize)
	}
	ipAddresses := 2 << (32 - cidrSize - 1)
	client = make([]bool, ipAddresses)
	client[0] = true // Allocate for server

	log.Printf("TUN_IPV4_CIDR: %s/%s %v", &tunIpv4Net.IP, s[1], ipAddresses)
}
