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

var c = new(config) // new to give us a pointer to a zero'd config so we don't accidentally null pointer like an idiot haha

func main() {
	err := awconf.LoadConfig("getwindowprocname", &c)
	if err != nil {
		log.Println(err)
		log.Println("Running in client mode") // Can't run the server if we don't know what port to listen on at all
		c.ServerMode = false
	}

	if c.ServerMode {
		log.Println("Running in server mode")
		startServer()
	} else {
		fmt.Printf(getProcessName()) // Printing only this to stdout and everything else to stderr so the program can more easily be used in scripts
	}
}
