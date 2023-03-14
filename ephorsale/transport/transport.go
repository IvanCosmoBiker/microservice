package transport

import (
	"bytes"
	"context"
	"io"

	dispetcher "ephorservices/ephorsale/transaction"
	logger "ephorservices/pkg/logger"
	MqttBroker "ephorservices/pkg/mqttmanager"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (lrw *LoggingResponseWriter) Write(p []byte) (int, error) {
	return lrw.buf.Write(p)
}

type Transport struct {
	Ctx            context.Context
	RequestManager *transportHttp.ServerHttp
	QueueManager   *MqttBroker.BrokerManager
}

var TransportManager *Transport

func New(ctx context.Context) *Transport {
	transport := &Transport{
		Ctx: ctx,
	}
	TransportManager = transport
	return transport
}

func (t *Transport) InitHttp(url, port string) {
	t.RequestManager = transportHttp.Init(url, port)
	t.RequestManager.AddMiddleWareFunc(t.MiddleWareLogin)
	t.RequestManager.SetMiddleWare()
}

func (t *Transport) InitMqtt(address, port, login, password, clientId string,
	protocolVersion, disconnect uint,
	backOffPolicySendMassage, backOffPolicyConnection []time.Duration,
	executeTimeSeconds int) error {
	MqttBroker, err := MqttBroker.New(t.Ctx)
	if err != nil {
		return err
	}
	t.QueueManager = MqttBroker
	MqttBroker.SetConfig(address,
		port, login, password, clientId,
		protocolVersion, disconnect,
		backOffPolicySendMassage, backOffPolicyConnection,
		executeTimeSeconds)
	err = t.QueueManager.Start()
	return err
}

func (t *Transport) MiddleWareLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		body, err := ioutil.ReadAll(r.Body)
		//defer r.Body.Close()
		if err != nil {
			logger.Log.Errorf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		t.addLog(r, body)
		t.RequestManager.SetBody(r, body)
		mrw := &LoggingResponseWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
		}
		next.ServeHTTP(mrw, r)
		t.setLog(r, mrw.buf, start)
		logger.Log.Info(fmt.Sprintf("%v", time.Since(start)))
		if _, err := io.Copy(w, mrw.buf); err != nil {
			logger.Log.Errorf("Failed to send out response: %v", err)
		}
	})
}

func (t *Transport) addLog(r *http.Request, body []byte) {
	logAdd := make(map[string]interface{})
	logAdd["date"] = dispetcher.Dispetcher.Date.Now()
	logAdd["address"] = r.RemoteAddr
	logAdd["request_id"] = dispetcher.Dispetcher.Date.UnixNano()
	logAdd["request_uri"] = r.RequestURI
	logAdd["request_data"] = string(body)
	structEntry, err := dispetcher.Dispetcher.StoreLog.AddByParams(logAdd)
	logger.Log.Info(fmt.Sprintf("%+v", structEntry))
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	idStrings := []string{strconv.Itoa(structEntry.Id)}
	logger.Log.Info(fmt.Sprintf("LOG::: %+v", structEntry))
	r.Header["log_id"] = idStrings
	structEntry = nil
}

func (t *Transport) setLog(r *http.Request, res *bytes.Buffer, start time.Time) {
	Now := time.Now()
	logIdSlice := r.Header["log_id"]
	logger.Log.Info(fmt.Sprintf("%v", logIdSlice))
	if len(logIdSlice) < 1 {
		logger.Log.Error("Not id log in header")
		return
	}
	for _, id := range logIdSlice {
		logId, errConvert := strconv.Atoi(id)
		logger.Log.Info(fmt.Sprintf("%T", logId))
		if errConvert != nil {
			logger.Log.Error(errConvert.Error())
			continue
		}
		logSet := make(map[string]interface{})
		logSet["id"] = id
		logSet["response"] = string(res.Bytes()[:])
		logSet["runtime"] = Now.Second() - start.Second()
		dispetcher.Dispetcher.StoreLog.SetByParams(logSet)
	}
}

func (t *Transport) Listen() error {
	go t.RequestManager.StartListener()
	return nil
}

func (t *Transport) Close() error {
	t.RequestManager.CloseListener()
	t.QueueManager.Shutdown(t.Ctx)
	return nil
}

func (t *Transport) CloseHttp() {
	t.RequestManager.CloseListener()
}

func (t *Transport) CloseMqtt() {
	t.QueueManager.Shutdown(t.Ctx)
}
