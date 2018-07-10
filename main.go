package main

import (
	"github.com/strider2038/serial-uuid-generator/config"
	"github.com/strider2038/serial-uuid-generator/server"
)

func main() {
	configuration := config.LoadConfigFromEnvironment()
	generatorServer := server.NewUUIDGeneratorServer(configuration)
	generatorServer.Run()
}
