package main

import (
	"log"
)

type Config struct {
	Port          string
	ServerAddress string
	ServerMode    bool
}

func main() {
	log.Println(getProcessName())
}
