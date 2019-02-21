package main

import "fmt"

type Config struct {
	Port          string
	ServerAddress string
	ServerMode    bool
}

func main() {
	fmt.Printf(getProcessName())
}
