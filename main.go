package main

import (
	"os"

	"github.com/ahmetozer/wstransit/client"
	"github.com/ahmetozer/wstransit/server"
)

func main() {

	if len(os.Args) == 1 {
		println(help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "client":
		client.Start(os.Args[2:])
	case "server":
		server.Start(os.Args[2:])
	default:
		println(help)
	}
}

const help = `Todo`
