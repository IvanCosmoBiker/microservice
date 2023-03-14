package sber

import (
	"bytes"
	"encoding/json"
	config "ephorservices/config"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Circuit func(method, url string, header map[string]interface{}, json_request []byte) (Response, error)

type Http struct {
	cfg *config.Config
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, json_request []byte, retries int, delay int) Circuit {
	return func(method, url string, header map[string]interface{}, json_request []byte) (Response, error) {
		for r := 0; r < retries; r++ {
			response, err := circuit(method, url, header, json_request)
			if err == nil || r >= retries {
				return response, err
			}
			log.Printf("Request %d failed SberPay; retrying in %v", r+1, delay)
			select {
			case <-time.After(time.Duration(delay) * time.Second):
				return Response{}, errors.New("TimeOut for send resuest SberPay")
			}
		}
		return Response{}, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, json_request []byte) (Response, error) {
	ResponseHttp := Response{}
	req, errReq := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	if errReq != nil {
		log.Printf("%v", errReq)
		return ResponseHttp, errReq
	}
	h.setHeader(header, req)
	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		ResponseHttp.StatusCode = 0
		return ResponseHttp, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	ResponseHttp.StatusCode = resp.StatusCode
	json.Unmarshal([]byte(body), &ResponseHttp.Data)
	if h.cfg.Debug {
		fmt.Printf("%s", body)
		log.Printf("%s", url)
		log.Printf("%+v", ResponseHttp)
	}
	return ResponseHttp, nil
}

func (h *Http) setHeader(header map[string]interface{}, req *http.Request) {
	for key, data := range header {
		req.Header.Set(key, parserTypes.ParseTypeInString(data))
	}
}

func InitHttp(conf *config.Config) *Http {
	return &Http{
		cfg: conf,
	}
}
