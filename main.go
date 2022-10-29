package main

import (
	"os"

	"github.com/ahmetozer/wstransit/cmd"
)

func main() {

	if len(os.Args) == 1 {
		println(help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "client":
		cmd.Client(os.Args[2:])
	case "server":
		cmd.Server(os.Args[2:])
	default:
		println(help)
	}
}

const help = `Todo`
