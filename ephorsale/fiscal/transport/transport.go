package transport

import (
	"bytes"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Circuit func(method, url string, header map[string]interface{}, json_request []byte) (int, error)

type Http struct {
	Status int
	Debug  bool
}

func New(debug bool) *Http {
	return &Http{
		Debug: debug,
	}
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, json_request []byte, retries int, delay int) Circuit {
	return func(method, url string, header map[string]interface{}, json_request []byte) (int, error) {
		for r := 0; r < retries; r++ {
			code, err := circuit(method, url, header, json_request)
			if err == nil || r >= retries {
				return code, err
			}
			log.Printf("Request %d failed  send; retrying in %v", r+1, delay)
			select {
			case <-time.After(time.Duration(delay) * time.Second):
				return code, errors.New("TimeOut for send resuest")
			}
		}
		return 0, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, json_request []byte) (int, error) {
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
	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return code, err
	}
	defer resp.Body.Close()
	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		log.Printf("%v", errBody)
		return code, errBody
	}
	if h.Debug {
		log.Printf("%+v", body)
		log.Printf("%s", body)
		log.Printf("%s", url)
	}
	return code, nil
}

// func SendToService(method, url string, json_request []byte) {
// 	log.Printf("%v", string(json_request))
// 	log.Printf("%s", url)
// 	req, err := http.NewRequest(method, url, bytes.NewBuffer(json_request))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Close = true
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("%v", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	log.Println(resp)
// 	log.Println(string(body))
// }
