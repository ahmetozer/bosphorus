package cmd

import (
	"flag"
	"sync"

	"github.com/ahmetozer/wstransit/client/flags"
	"github.com/ahmetozer/wstransit/client/tcp"
	"github.com/ahmetozer/wstransit/client/tun"
)

func Client(args []string) {

	f := flag.NewFlagSet("client", flag.ExitOnError)
	f.Var(&flags.TcpFlag, "tcp", "tunnel example localhost:8022;ahmet.engineer;127.0.0.1:22")
	//f.Var(&flags.TcpRawFlag, "tcpraw", "tunnel example localhost:8022;ahmet.engineer;127.0.0.1:22")
	f.Var(&flags.TunFlag, "tun", "tunnel example utun10;https://ahmet.engineer")

	f.Parse(args)

	wg := new(sync.WaitGroup)
	if len(flags.TcpFlag) > 0 {
		tcp.TcpStart(flags.TcpFlag, wg, tcp.TCPTypeListen)
	}

	// if len(flags.TcpRawFlag) > 0 {
	// 	tcp.TcpStart(flags.TcpRawFlag, wg, tcp.TCPTypeRawListen)
	// }

	if len(flags.TunFlag) > 0 {
		tun.TunStart(flags.TunFlag, wg)
	}

	wg.Wait()

}
