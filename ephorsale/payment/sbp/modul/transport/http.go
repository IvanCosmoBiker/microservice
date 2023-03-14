package transport

import (
	"bytes"
	"encoding/json"
	response_modul "ephorservices/ephorsale/payment/sbp/modul/response"
	logger "ephorservices/pkg/logger"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Circuit func(method, url string, header map[string]interface{}, json_request []byte, response response_modul.Response) (int, error)

type Http struct {
	Status int
	Debug  bool
}

func InitHttp(debug bool) *Http {
	return &Http{
		Debug: debug,
	}
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, json_request []byte, response response_modul.Response, retries int, delay int) Circuit {
	return func(method, url string, header map[string]interface{}, json_request []byte, response response_modul.Response) (int, error) {
		for r := 0; r < retries; r++ {
			code, err := circuit(method, url, header, json_request, response)
			if err == nil || r >= retries {
				return code, err
			}
			logger.Log.Infof("Request %d failed Modul Sbp; retrying in %v", r+1, delay)
			select {
			case <-time.After(time.Duration(delay) * time.Second):
				return code, errors.New("TimeOut for send resuest Modul Sbp")
			}
		}
		return 0, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, json_request []byte, response response_modul.Response) (int, error) {
	code := 0
	log.Printf("%s", method)
	req, errReq := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	if errReq != nil {
		logger.Log.Errorf("%+v", errReq)
		return code, errReq
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for key, data := range header {
		req.Header.Set(key, parserTypes.ParseTypeInString(data))
	}
	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return code, err
	}
	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		logger.Log.Errorf("%+v", errBody)
		return code, errBody
	}
	code = resp.StatusCode
	defer resp.Body.Close()
	json.Unmarshal([]byte(body), response)
	if h.Debug {
		logger.Log.Infof("%+v", response)
		logger.Log.Infof("%s", body)
		logger.Log.Infof("%s", url)
	}
	return code, nil
}
