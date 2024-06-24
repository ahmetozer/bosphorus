package tun

import (
	"log"
	"sync"

	"github.com/ahmetozer/bosphorus/client/flags"
)

func TunStart(a flags.ArrFlag, wg *sync.WaitGroup) {
	for _, flag := range a {
		connectionString, err := parseTUNConnectionString(flag)
		if err != nil {
			log.Fatalf("tcp flag parse: %s, %s ", err, flag)
		}
		go tunListener(connectionString, wg)
		wg.Add(1)
	}
}
