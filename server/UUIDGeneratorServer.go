package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	log "github.com/sirupsen/logrus"
	"github.com/strider2038/serial-uuid-generator/config"
	"github.com/strider2038/serial-uuid-generator/generator"
	"github.com/strider2038/serial-uuid-generator/service"
)

type uuidGeneratorServer struct {
	handler http.Handler
	config  config.Config
}

func NewUUIDGeneratorServer(config config.Config) Server {
	initLogger(config)

	generatorServer := new(uuidGeneratorServer)
	generatorServer.handler = createRouter(config)
	generatorServer.config = config

	return generatorServer
}

func initLogger(config config.Config) {
	log.SetLevel(config.LogLevel)

	if config.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
}

func createRouter(config config.Config) *mux.Router {
	valueStorage := generator.NewPostgresValueStorage(config.DatabaseUrl, config.TableName)
	valueGenerator := generator.NewSequenceValueGenerator(valueStorage, config.RangeStep)
	generatorService := service.NewGenerator(valueGenerator)

	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	rpcServer.RegisterService(generatorService, "")

	router := mux.NewRouter()
	router.Handle("/rpc", rpcServer)

	return router
}

func (server *uuidGeneratorServer) Run() error {
	log.
		WithFields(log.Fields{
			"port":     server.config.Port,
			"logLevel": server.config.LogLevel.String(),
		}).
		Info("Starting Serial UUID Generator server")

	host := fmt.Sprintf(":%d", server.config.Port)

	return http.ListenAndServe(host, server.handler)
}
