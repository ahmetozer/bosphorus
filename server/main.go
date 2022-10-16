package server

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func Start(args []string) {
	http.Handle("/", websocket.Handler(handleNewClient))
	if err := http.ListenAndServe(os.Getenv("ADDR"), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
