package transport

import (
	"bytes"
	"encoding/json"
	responseNanokass "ephorservices/ephorsale/fiscal/nanokass/response"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Circuit func(method, url string, header map[string]interface{}, json_request []byte, response responseNanokass.Response) (int, error)

type Http struct {
	Status int
	Debug  bool
}

func InitHttp(debug bool) *Http {
	return &Http{
		Debug: debug,
	}
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, json_request []byte, response responseNanokass.Response, retries int, delay int) Circuit {
	return func(method, url string, header map[string]interface{}, json_request []byte, response responseNanokass.Response) (int, error) {
		for r := 0; r < retries; r++ {
			code, err := circuit(method, url, header, json_request, response)
			if err == nil || r >= retries {
				return code, err
			}
			log.Printf("Request %d failed Nanokassa; retrying in %v", r+1, delay)
			select {
			case <-time.After(time.Duration(delay) * time.Second):
				return code, errors.New("TimeOut send Nanokassa")
			}
		}
		return 0, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, json_request []byte, response responseNanokass.Response) (int, error) {
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
