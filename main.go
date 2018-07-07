package main

import (
	"github.com/strider2038/serial-uuid-generator/config"
	"github.com/strider2038/serial-uuid-generator/server"
)

func main() {
	configuration := config.Config{
		DatabaseUrl: "postgres://user:password@localhost/generator?sslmode=disable",
		TableName:   "public.uuid_sequence",
		RangeStep:   100,
	}
	generatorServer := server.NewUUIDGeneratorServer(configuration)
	generatorServer.Run()
}
