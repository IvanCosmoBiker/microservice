package transport

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	responseOrange "ephorservices/ephorsale/fiscal/orange/response"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Circuit func(method, url string, header map[string]interface{}, cert, key string, json_request []byte, response responseOrange.Response) (int, error)

type Http struct {
	Status int
	Debug  bool
}

func InitHttp(debug bool) *Http {
	return &Http{
		Debug: debug,
	}
}

func (h *Http) Send(circuit Circuit, method, url string, header map[string]interface{}, cert, key string, json_request []byte, response responseOrange.Response, retries int, delay int) Circuit {
	return func(method, url string, header map[string]interface{}, cert, key string, json_request []byte, response responseOrange.Response) (int, error) {
		for r := 0; r < retries; r++ {
			code, err := circuit(method, url, header, cert, key, json_request, response)
			if err == nil || r >= retries {
				return code, err
			}
			log.Printf("Request %d failed Orange Data send; retrying in %v", r+1, delay)
			select {
			case <-time.After(time.Duration(60) * time.Second):
				return code, errors.New("TimeOut send request Orange Data")
			}
		}
		return 0, nil
	}
}

func (h *Http) Call(method, url string, header map[string]interface{}, cert, key string, json_request []byte, response responseOrange.Response) (int, error) {
	code := 0
	log.Printf("%s", method)
	req, errReq := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	if errReq != nil {
		log.Printf("%v", errReq)
		return code, errReq
	}
	for key, data := range header {
		req.Header.Set(key, parserTypes.ParseTypeInString(data))
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	fmt.Printf("\nHEAD: %+v\n", req.Header)
	req.Close = true
	fmt.Printf("\nCERTIFICATE %v\n", cert)
	fmt.Printf("\nKEY %v\n", key)
	Cert, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{Cert},
		InsecureSkipVerify: true,
	}

	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
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
	log.Printf("%+v", body)
	code = resp.StatusCode
	defer resp.Body.Close()
	json.Unmarshal(body, response)
	fmt.Printf("\nRESPONSE : %v\n", response)
	if h.Debug {
		log.Printf("%+v", response)
		log.Printf("%s", body)
		log.Printf("%s", url)
	}
	return code, nil
}
