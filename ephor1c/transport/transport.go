package transport

import (
	"ephorservices/ephor1c/config"
	warehouse_service "ephorservices/ephor1c/transport/server/warehouse"
	loggerEvent "ephorservices/pkg/logger"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"fmt"
	"net/http"
	"time"
)

type Transport struct {
	Config         *config.Config
	EventLogger    *loggerEvent.LoggerEvent
	RequestManager *transportHttp.ServerHttp
}

func New(cfg *config.Config) *Transport {
	return &Transport{
		Config: cfg,
	}
}

func (t *Transport) InitServer(url, port string, logger *loggerEvent.LoggerEvent) {
	t.RequestManager = transportHttp.Init(url, port)
	t.RequestManager.EventLogger = logger
	t.RequestManager.AddMiddleWareFunc(t.MiddleWareLogin)
	t.RequestManager.SetMiddleWare()
	t.InitUrl()
	t.Listen()
}

func (t *Transport) MiddleWareLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//t.addLog(r)
		next.ServeHTTP(w, r)
		t.EventLogger.Info(fmt.Sprintf("%v", time.Since(start)))
		t.EventLogger.Info(fmt.Sprintf("%v", w))
	})
}

// func (t *Transport) addLog(r *http.Request) {
// 	json_data, _ := ioutil.ReadAll(r.Body)
// 	defer t.RequestManager.setBody(r, json_data)
// 	defer r.Body.Close()
// 	logAdd := make(map[string]interface{})
// 	logAdd["date"] = time.Now()
// 	logAdd["address"] = r.RemoteAddr
// 	logAdd["request_id"] = time.StampMilli()
// 	logAdd["request_uri"] = r.RequestURI
// 	logAdd["request_data"] = string(json_data[:])
// 	//structEntry, err := .Dispetcher.StoreLog.AddByParams(logAdd)
// 	// if err != nil {
// 	// 	server.EventLogger.Info(err.Error())
// 	// 	return
// 	// }
// 	//idStrings := []string{strconv.Itoa(structEntry.Id)}
// 	//r.Header["log_id"] = idStrings
// }

func (t *Transport) InitUrl() {
	WareHouseService := warehouse_service.New()
	WareHouseService.InitUrl(t.RequestManager)
}

func (t *Transport) Listen() {
	go t.RequestManager.StartListener()
}
