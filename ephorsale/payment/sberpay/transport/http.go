package transport

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	responseInterface "ephorservices/ephorsale/payment/sberpay/response"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var sert = []byte(`-----BEGIN CERTIFICATE-----
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
MW4bZis=
-----END CERTIFICATE-----`)

type Circuit func(method, url string, header map[string]interface{}, json_request []byte, response responseInterface.Response) (int, error)

type Http struct {
	Status int
	Debug  bool
}

func InitHttp(debug bool) *Http {
	return &Http{
		Debug: debug,
	}
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, json_request []byte, response responseInterface.Response, retries int, delay int) Circuit {
	return func(method, url string, header map[string]interface{}, json_request []byte, response responseInterface.Response) (int, error) {
		for r := 0; r < retries; r++ {
			code, err := circuit(method, url, header, json_request, response)
			if err == nil || r >= retries {
				return code, err
			}
			log.Printf("Request %d failed SberPay; retrying in %v", r+1, delay)
			select {
			case <-time.After(time.Duration(delay) * time.Second):
				return code, errors.New("TimeOut for send resuest SberPay")
			}
		}
		return 0, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, json_request []byte, response responseInterface.Response) (int, error) {
	code := 0
	log.Printf("%s", method)
	req, errReq := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	if errReq != nil {
		log.Printf("%v", errReq)
		return code, errReq
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for key, data := range header {
		req.Header.Set(key, parserTypes.ParseTypeInString(data))
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	req.Close = true
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return code, err
	}
	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
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
	return code, nil
}

func (h *Http) setHeader(header map[string]interface{}, req *http.Request) {
	for key, data := range header {
		req.Header.Set(key, parserTypes.ParseTypeInString(data))
	}
}
