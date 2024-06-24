package tun

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ahmetozer/bosphorus/pkg/conn"
)

var (
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
		slog.Debug("TUN_IPV6_CIDR is not asigned, auto value used")
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

	slog.Debug(fmt.Sprintf("TUN_IPV6_CIDR: %s/%s\n", &tunIpv6Net.IP, s[1]))

}

func envIpv4Cidr() {
	tunIPv4Cidr = os.Getenv("TUN_IPV4_CIDR")
	var err error
	if tunIPv4Cidr == "" {
		tunIPv4Cidr = "10.90.0.1/24"
		slog.Debug("TUN_IPV4_CIDR is not asigned, auto value used\n")
	}

	tunIpv4Ip, tunIpv4Net, err = net.ParseCIDR(tunIPv4Cidr)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	s := strings.Split(tunIPv4Cidr, "/")
	cidrSize, _ := strconv.Atoi(s[1])
	if cidrSize > 24 {
		log.Fatalf("TUN_IPV4_CIDR '%d' size is low. cidr range must bigger or equal to 24\n", cidrSize)
	}

	slog.Debug(fmt.Sprintf("TUN_IPV4_CIDR: %s/%s", &tunIpv4Net.IP, s[1]))
}
