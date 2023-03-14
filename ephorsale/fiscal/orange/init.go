package orange

import (
	"ephorservices/ephorsale/fiscal/interface/fr"
	core "ephorservices/ephorsale/fiscal/orange/core"
	responseOrange "ephorservices/ephorsale/fiscal/orange/response"
	transport "ephorservices/ephorsale/fiscal/orange/transport"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"fmt"
	"log"
	"strings"
	"time"
)

type Orange struct {
	retries        int
	delay          int
	Name           string
	ExecuteMinutes int
	SleepSecond    int
	Core           *core.Core
	Http           *transport.Http
}

type NewOrangeStruct struct {
	Orange
}

func (nos *NewOrangeStruct) New(executeMinutes, sleepSecond int, pathCert string, debug bool) fr.Fiscal {
	Core := core.InitCore(pathCert)
	Http := transport.InitHttp(debug)
	return &NewOrangeStruct{
		Orange: Orange{
			retries:        10,
			delay:          5,
			Name:           "Orange",
			ExecuteMinutes: executeMinutes,
			SleepSecond:    sleepSecond,
			Core:           Core,
			Http:           Http,
		},
	}
}

func (o *Orange) SendCheck(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	signature, resiptId, jsonRequest, err := o.Core.MakeRequestSendCheck(tran)
	if err != nil {
		signature = ""
		resiptId = ""
		jsonRequest = nil
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", err)
		return result
	}
	tran.Fiscal.Signature = signature
	url := fmt.Sprintf("https://%v:%v/api/v2/documents/", tran.Fiscal.Config.Dev_addr, tran.Fiscal.Config.Dev_port)
	Response := &responseOrange.ResponseCreateCheck{}
	headers := make(map[string]interface{})
	method := "POST"
	headers["x-Signature"] = tran.Fiscal.Signature
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Content-Length"] = len(jsonRequest)
	funcSend := o.Http.Send(o.Http.Call, method, url, headers, tran.Fiscal.Config.Auth_public_key, tran.Fiscal.Config.Auth_private_key, jsonRequest, Response, o.retries, o.delay)
	code, errResp := funcSend(method, url, headers, tran.Fiscal.Config.Auth_public_key, tran.Fiscal.Config.Auth_private_key, jsonRequest, Response)
	log.Printf("%v", code)
	if errResp != nil {
		Response = nil
		funcSend = nil
		headers = nil
		return o.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
	}
	tran.Fiscal.ResiptId = resiptId
	if code == 409 {
		Response = nil
		funcSend = nil
		headers = nil
		tran.Fiscal.Code = 200
		tran.Fiscal.StatusCode = 200
		result["status"] = fr.Status_InQueue
		return result
	}
	if code != 201 {
		funcSend = nil
		headers = nil
		tran.Fiscal.Code = 200
		tran.Fiscal.StatusCode = 200
		result["fr_id"] = resiptId
		result["status"] = fr.Status_Error
		result["fr_status"] = fr.Status_Error
		result["f_desc"] = strings.Join(Response.Errors[:], "\n")
		Response = nil
		return result
	}
	Response = nil
	funcSend = nil
	headers = nil
	result["status"] = fr.Status_InQueue
	return result
}

func (o *Orange) GetStatus(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	url := o.Core.MakeRequestStatusQr(tran)
	timeout := time.NewTimer(time.Duration(o.ExecuteMinutes) * time.Minute)
	headers := make(map[string]interface{})
	headers["x-Signature"] = tran.Fiscal.Signature
	headers["Content-Type"] = "application/json; charset=utf-8"
	Response := &responseOrange.ResponseStatusCheck{}
	method := "GET"
	funcSend := o.Http.Send(o.Http.Call, method, url, headers, tran.Fiscal.Config.Auth_public_key, tran.Fiscal.Config.Auth_private_key, []byte(""), Response, o.retries, o.delay)
	for {
		select {
		case <-timeout.C:
			{
				tran.Fiscal.Code = 400
				tran.Fiscal.StatusCode = 400
				result["status"] = fr.Status_Error
				result["f_desc"] = fmt.Sprintf("Cancelled by a Timeout of %s", o.Name)
				result["fr_status"] = fr.Status_Error
				return result
			}
		case <-time.After(time.Duration(o.SleepSecond) * time.Second):
			{

				code, errResp := funcSend(method, url, headers, tran.Fiscal.Config.Auth_public_key, tran.Fiscal.Config.Auth_private_key, []byte(""), Response)
				if errResp != nil {
					funcSend = nil
					return o.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
				}
				log.Printf("%v", code)
				if code > 299 {
					funcSend = nil
					tran.Fiscal.Code = code
					tran.Fiscal.StatusCode = code
					result["fr_id"] = tran.Fiscal.ResiptId
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = strings.Join(Response.Errors[:], "\n")
					fmt.Printf("%+v\n", Response)
					Response = nil
					return result
				}
				if code == 200 {
					tran.Fiscal.Code = code
					tran.Fiscal.StatusCode = code
					TimeOrange, _ := tran.DateTime.ParseDateOfLayout("2006-01-02T15:04:05", Response.ProcessedAt)
					tran.Fiscal.Fields.Fp = Response.Fp
					tran.Fiscal.Fields.Fn = Response.FsNumber
					tran.Fiscal.Fields.Fd = float64(Response.DocumentNumber)
					tran.Fiscal.Fields.DateFisal = TimeOrange
					result["fr_id"] = tran.Fiscal.ResiptId
					result["status"] = fr.Status_Complete
					result["fr_status"] = fr.Status_Complete
					result["f_desc"] = strings.Join(Response.Errors[:], "\n")
					Response = nil
					return result
				}
			}
		}
	}
}

func (o *Orange) GetQrPicture(tran *transaction.Transaction) string {
	return o.Core.EncodeUrlToBase64(tran)
}

func (o *Orange) GetQrUrl(tran *transaction.Transaction) string {
	return o.Core.MakeUrlQr(tran)
}

func (o *Orange) MakeResultError(massage string, err error, result map[string]interface{}) map[string]interface{} {
	o.ClearResultMap(result)
	result["status"] = transaction.TransactionState_Error
	result["f_desc"] = fmt.Sprintf("%v", err)
	result["message"] = massage
	return result
}

func (o *Orange) ClearResultMap(result map[string]interface{}) {
	for k := range result {
		delete(result, k)
	}
}
