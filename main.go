package main

import "github.com/strider2038/serial-uuid-generator/server"

func main() {
	generatorServer := server.NewUUIDGeneratorServer()
	generatorServer.Run()
}
