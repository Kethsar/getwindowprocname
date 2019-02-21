package main

import (
	"fmt"
	"log"

	"github.com/awused/awconf"
)

type config struct {
	Port          string
	ServerAddress string
	ServerMode    bool
}

var c config

func main() {
	err := awconf.LoadConfig("getwindowprocname", c)
	if err != nil {
		log.Println(err)
		log.Println("Running in client mode")
		c.ServerMode = false
	}

	if c.ServerMode {
		log.Println("Running in server mode (not yet implemented)")
	} else {
		fmt.Printf(getProcessName())
	}
}
