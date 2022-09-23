package v1

import (
	config "ephorservices/config"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ServerHttp struct {
	RouterHttp *mux.Router
	Server     *http.Server
}

func (server *ServerHttp) Init(url, port string) {
	server.InitRouter()
	server.InitServer(url, port)
	return
}

func (server *ServerHttp) StartListener() {
	if err := server.Server.ListenAndServe(); err != nil {
		log.Println(err)
		runtime.Goexit()
	}
	log.Println("Start listeners Server")
	return
}

func (server *ServerHttp) InitRouter() {
	server.RouterHttp = mux.NewRouter()
	return
}

func (server *ServerHttp) InitServer(url, port string) {
	Address := fmt.Sprintf("%s:%s", url, port)
	log.Println(Address)
	s := &http.Server{
		Addr:         Address,
		Handler:      server.RouterHttp,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	server.Server = s
	return
}

func (server *ServerHttp) SetHandlerListener(address string, handler func(w http.ResponseWriter, req *http.Request)) *mux.Route {
	router := server.RouterHttp.HandleFunc(address, handler)
	return router
}
func (server *ServerHttp) SetMiddleWare() {
	server.RouterHttp.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(log.Writer(), h)
	})
	server.RouterHttp.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	return
}

/*
   Завершает работу http сервера аккуратно. Закрывает соединения для новых и ждёт завершения текущих соединений используя пакет context.
*/
func (server *ServerHttp) CloseListener() {
	if err := server.Server.Close(); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	} else {
		log.Print("HTTP server Shutdown completed successfully")
	}
	return
}

func Init(conf *config.Config) *ServerHttp {
	httpServer := &ServerHttp{}
	httpServer.Init(conf.Services.Http.Address, conf.Services.Http.Port)
	return httpServer
}
