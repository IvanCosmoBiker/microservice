package fermaOfd

import (
	"bytes"
	"encoding/json"
	responseOfd "ephorservices/ephorsale/fiscal/fermaOfd/response"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Circuit func(method, url string, header map[string]interface{}, json_request []byte, response responseOfd.Response) (int, error)

type Http struct {
	Debug  bool
	Status int
}

func InitHttp(debug bool) *Http {
	return &Http{
		Debug: debug,
	}
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, json_request []byte, response responseOfd.Response, retries int, delay int) Circuit {
	return func(method, url string, header map[string]interface{}, json_request []byte, response responseOfd.Response) (int, error) {
		for r := 0; r < retries; r++ {
			code, err := circuit(method, url, header, json_request, response)
			if err == nil || r >= retries {
				return code, err
			}
			log.Printf("Request %d failed Ferma Ofd; retrying in %v", r+1, delay)
			select {
			case <-time.After(time.Duration(delay) * time.Second):
				return code, errors.New("TimeOut send request Ferma Ofd")
			}
		}
		return 0, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, json_request []byte, response responseOfd.Response) (int, error) {
	code := 0
	req, errReq := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	if errReq != nil {
		req = nil
		log.Printf("%v", errReq)
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
		client = nil
		resp = nil
		req = nil
		fmt.Println(err)
		return code, err
	}
	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		body = nil
		req = nil
		client = nil
		resp = nil
		log.Printf("%v", errBody)
		return code, errBody
	}
	code = resp.StatusCode
	defer resp.Body.Close()
	json.Unmarshal([]byte(body), response)
	if h.Debug {
		log.Printf("%+v", response)
		log.Printf("%s", body)
		log.Printf("%s", url)
	}
	body = nil
	req = nil
	client = nil
	resp = nil
	return code, nil
}
