package cmd

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ahmetozer/bosphorus/server"
	"golang.org/x/net/websocket"
)

func Server(args []string) {

	if os.Getenv("ADDR") == "" {
		os.Setenv("ADDR", ":8080")
	}
	httpSrv := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           websocket.Handler(server.HandleNewClient),
		Addr:              os.Getenv("ADDR"),
	}

	log.Printf("server is listening on %s", os.Getenv("ADDR"))
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
