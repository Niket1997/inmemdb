package main

import (
	"flag"
	"github.com/Niket1997/inmemdb/server"
	"log"

	"github.com/Niket1997/inmemdb/config"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the inmemdb server")
	flag.IntVar(&config.Port, "port", 7379, "port for the inmemdb server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("starting the inmemdb server :)")
	server.RunSyncTCPServer()
}
