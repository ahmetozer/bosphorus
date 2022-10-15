package client

import (
	"flag"
	"sync"
)

func Main(args []string) {

	f := flag.NewFlagSet("client", flag.ExitOnError)
	f.Var(&tcpFlag, "tcp", "tunnel example localhost:8022;ahmet.engineer;127.0.0.1:8118")
	f.Parse(args)

	wg := new(sync.WaitGroup)
	if len(tcpFlag) > 0 {
		tcpStart(tcpFlag, wg)
	}

	wg.Wait()

}
