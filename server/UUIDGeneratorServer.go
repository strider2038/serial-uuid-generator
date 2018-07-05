package server

import (
	"net/http"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/strider2038/serial-uuid-generator/service"
)

type uuidGeneratorServer struct {
	handler http.Handler
}

func NewUUIDGeneratorServer() Server {
	generator := new(service.Generator)

	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	rpcServer.RegisterService(generator, "")

	router := mux.NewRouter()
	router.Handle("/rpc", rpcServer)

	generatorServer := new(uuidGeneratorServer)
	generatorServer.handler = router

	return generatorServer
}

func (server *uuidGeneratorServer) Run() error {
	fmt.Println("Starting Serial UUID Generator uuidGeneratorServer on port 3000...")

	return http.ListenAndServe(":3000", server.handler)
}
