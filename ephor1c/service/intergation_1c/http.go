package intergation_1c

import (
	"bytes"
	"encoding/json"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Code int
	Data map[string]interface{}
}

type Circuit func(method, url string, header map[string]interface{}, json_request []byte, response *Response) (int, error)

type Http struct {
	Status int
	Retrie int
	Delay  int
}

func New() *Http {
	return &Http{
		Retrie: 3,
		Delay:  5,
	}
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, json_request []byte, response *Response) Circuit {
	return func(method, url string, header map[string]interface{}, json_request []byte, response *Response) (int, error) {
		for r := 0; r < h.Retrie; r++ {
			code, err := circuit(method, url, header, json_request, response)
			if err == nil || r >= h.Retrie {
				return code, err
			}
			log.Printf("Request %d failed 1C; retrying in %v", r+1, h.Delay)
			select {
			case <-time.After(time.Duration(h.Delay) * time.Second):
				return code, errors.New("TimeOut for send resuest 1C")
			}
		}
		return 0, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, json_request []byte, response *Response) (int, error) {
	code := 0
	log.Printf("%s", method)
	req, errReq := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	if errReq != nil {
		log.Printf("%v", errReq)
		return code, errReq
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	h.setHeader(header, req)
	req.Close = true
	client := &http.Client{}
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
	return code, nil
}

func (h *Http) setHeader(header map[string]interface{}, req *http.Request) {
	for key, data := range header {
		req.Header.Set(key, parserTypes.ParseTypeInString(data))
	}
}
