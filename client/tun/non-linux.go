//go:build !linux

package tun

import (
	"log"
	"sync"

	"github.com/ahmetozer/bosphorus/client/flags"
)

func TunStart(a flags.ArrFlag, wg *sync.WaitGroup) {
	log.Printf("tun is only supported for linux")
}
