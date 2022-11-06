package tun

import (
	"encoding/binary"
	"math/bits"
	"net"
)

type uint128 [16]byte

func IpUint64Add(ip net.IP, a uint64) net.IP {

	lo, carry := bits.Add64(binary.BigEndian.Uint64(ip[8:16]), a, 0)

	blo := make([]byte, 8)
	bhi := make([]byte, 8)

	binary.BigEndian.PutUint64(blo, lo)
	binary.BigEndian.PutUint64(bhi, binary.BigEndian.Uint64(ip[0:8])+carry)

	return append(bhi, blo...)
}
