package main

import (
	"flag"
	"log"

	"github.com/Niket1997/inmemdb/config"
	"github.com/Niket1997/inmemdb/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the inmemdb server")
	flag.IntVar(&config.Port, "port", 7379, "port for the inmemdb server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("starting the inmemdb server :)")
	err := server.RunAsyncTCPServer()
	if err != nil {
		log.Println(err)
		return
	}
}
