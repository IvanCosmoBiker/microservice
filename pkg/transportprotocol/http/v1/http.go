package v1

import (
	"bytes"
	transactionDispetcher "ephorservices/ephorsale/transaction"
	logger "ephorservices/pkg/logger"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ServerHttp struct {
	RouterHttp      *mux.Router
	Server          *http.Server
	Dispetcher      *transactionDispetcher.TransactionDispetcher
	MiddleWareFuncs []mux.MiddlewareFunc
}

func (server *ServerHttp) Init(url, port string) {
	server.InitRouter()
	server.InitServer(url, port)
	server.InitMiddleWareFunc()
}

func (server *ServerHttp) InitRouter() {
	server.RouterHttp = mux.NewRouter()
}

func (server *ServerHttp) InitServer(url, port string) {
	Address := fmt.Sprintf("%s:%s", url, port)
	s := &http.Server{
		Addr:         Address,
		Handler:      server.RouterHttp,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	server.Server = s
}

func (server *ServerHttp) InitMiddleWareFunc() {
	server.MiddleWareFuncs = make([]mux.MiddlewareFunc, 0, 4)
	server.MiddleWareFuncs = append(server.MiddleWareFuncs, handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
}

func (server *ServerHttp) AddMiddleWareFunc(middleware ...mux.MiddlewareFunc) {
	server.MiddleWareFuncs = append(server.MiddleWareFuncs, middleware...)
}

func (server *ServerHttp) SetHandlerListener(address string, handler func(w http.ResponseWriter, req *http.Request), method ...string) *mux.Route {
	router := server.RouterHttp.HandleFunc(address, handler)
	router.Methods(method...)
	return router
}

func (server *ServerHttp) SetMiddleWare() {
	server.RouterHttp.Use(server.MiddleWareFuncs...)
}

func (server *ServerHttp) SetBody(r *http.Request, data []byte) {
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
}

func (server *ServerHttp) StartListener() {
	if err := server.Server.ListenAndServe(); err != nil {
		logger.Log.Error(err.Error())
		runtime.Goexit()
	}
	logger.Log.Info("Start listeners Server")
}

/*
   Завершает работу http сервера аккуратно. Закрывает соединения для новых и ждёт завершения текущих соединений используя пакет context.
*/
func (server *ServerHttp) CloseListener() {
	if err := server.Server.Close(); err != nil {
		logger.Log.Info(fmt.Sprintf("HTTP server Shutdown: %v", err))
	} else {
		logger.Log.Info("HTTP server Shutdown completed successfully")
	}
}

func Init(addres, port string) *ServerHttp {
	httpServer := &ServerHttp{}
	httpServer.Init(addres, port)
	return httpServer
}
