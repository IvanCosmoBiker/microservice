package fermaOfd

import (
	responseOfd "ephorservices/ephorsale/fiscal/fermaOfd/response"
	"ephorservices/ephorsale/fiscal/interface/fr"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"regexp"
	"time"
)

type FermaOfd struct {
	retries        int
	delay          int
	Name           string
	ExecuteMinutes int
	SleepSecond    int
	Core           *Core
	Http           *Http
}

type NewFermaOfdStruct struct {
	FermaOfd
}

func (nos *NewFermaOfdStruct) New(executeMinutes, sleepSecond int, debug bool) fr.Fiscal {
	Core := InitCore()
	Http := InitHttp(debug)
	return &NewFermaOfdStruct{
		FermaOfd: FermaOfd{
			retries:        10,
			delay:          5,
			Name:           "FermaOfd",
			ExecuteMinutes: executeMinutes,
			SleepSecond:    sleepSecond,
			Core:           Core,
			Http:           Http,
		},
	}
}

func (fo *FermaOfd) MakeAuth(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	jsonRequest, err := fo.Core.MakeRequestAuth(tran)
	if err != nil {
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", err)
		return result
	}
	url := fmt.Sprintf("https://%s/api/Authorization/CreateAuthToken", tran.Fiscal.Config.Dev_addr)
	Response := &responseOfd.ResponseAuth{}
	headers := make(map[string]interface{})
	method := "POST"
	funcSend := fo.Http.Send(fo.Http.Call, method, url, headers, jsonRequest, Response, fo.retries, fo.delay)
	code, errResp := funcSend(method, url, headers, jsonRequest, Response)
	if errResp != nil {
		jsonRequest = nil
		funcSend = nil
		Response = nil
		headers = nil
		url = ""
		return fo.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
	}
	if code != 200 {
		tran.Fiscal.Code = code
		tran.Fiscal.StatusCode = code
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", Response.Error.Message)
		jsonRequest = nil
		funcSend = nil
		Response = nil
		headers = nil
		url = ""
		return result
	}
	tran.Fiscal.AuthToken = Response.Data.AuthToken
	if len(tran.Fiscal.AuthToken) == 0 {
		jsonRequest = nil
		funcSend = nil
		Response = nil
		headers = nil
		url = ""
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("Error: %s No auth token", fo.Name)
		return result
	}
	tran.Fiscal.Code = code
	tran.Fiscal.StatusCode = code
	result["status"] = fr.Status_InQueue
	jsonRequest = nil
	funcSend = nil
	Response = nil
	headers = nil
	url = ""
	return result
}

func (fo *FermaOfd) SendCheck(tran *transaction.Transaction) map[string]interface{} {
	resultAuth := fo.MakeAuth(tran)
	if resultAuth["status"] != fr.Status_InQueue {
		return resultAuth
	}
	result := make(map[string]interface{})
	jsonRequest, err := fo.Core.MakeRequestSendCheck(tran)
	if err != nil {
		jsonRequest = nil
		tran.Fiscal.Code = 400
		tran.Fiscal.StatusCode = 400
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", err)
		return result
	}
	method := "POST"
	url := fmt.Sprintf("https://%s/api/kkt/cloud/receipt?AuthToken=%s", tran.Fiscal.Config.Dev_addr, tran.Fiscal.AuthToken)
	Response := &responseOfd.ResponseSendCheck{}
	headers := make(map[string]interface{})
	funcSend := fo.Http.Send(fo.Http.Call, method, url, headers, jsonRequest, Response, fo.retries, fo.delay)
	code, errResp := funcSend(method, url, headers, jsonRequest, Response)
	if errResp != nil {
		jsonRequest = nil
		funcSend = nil
		Response = nil
		headers = nil
		return fo.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
	}
	if code != 200 {
		if Response.Error.Code == 1019 {
			re := regexp.MustCompile("--([a-f0-9-]+)--")
			match := re.FindStringSubmatch(Response.Error.Message)
			if len(match) > 1 {
				tran.Fiscal.ResiptId = match[1]
				result["status"] = fr.Status_Retry_Status
				return result
			}
			tran.Fiscal.Code = 200
			tran.Fiscal.StatusCode = 200
		}
		result["status"] = fr.Status_Error
		jsonRequest = nil
		funcSend = nil
		Response = nil
		headers = nil
		return result
	}
	tran.Fiscal.ResiptId = Response.Data.ReceiptId
	tran.Fiscal.Code = code
	tran.Fiscal.StatusCode = code
	if len(tran.Fiscal.ResiptId) == 0 {
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("Error: %s No ReceiptId", fo.Name)
		funcSend = nil
		Response = nil
		headers = nil
		jsonRequest = nil
		return result
	}
	result["status"] = fr.Status_InQueue
	jsonRequest = nil
	funcSend = nil
	Response = nil
	headers = nil
	return result
}

func (fo *FermaOfd) GetStatus(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	jsonData, err := fo.Core.MakeRequestStatusCheck(tran)
	if err != nil {
		jsonData = nil
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", err)
		return result
	}
	method := "POST"
	timeout := time.NewTimer(time.Duration(fo.ExecuteMinutes) * time.Minute)
	headers := make(map[string]interface{})
	url := fmt.Sprintf("https://%s/api/kkt/cloud/status?AuthToken=%s", tran.Fiscal.Config.Dev_addr, tran.Fiscal.AuthToken)
	Response := &responseOfd.ResponseStatusCheck{}
	funcSend := fo.Http.Send(fo.Http.Call, method, url, headers, jsonData, Response, fo.retries, fo.delay)
	for {
		select {
		case <-timeout.C:
			{
				tran.Fiscal.Code = 400
				tran.Fiscal.StatusCode = 400
				result["status"] = fr.Status_Error
				result["f_desc"] = fmt.Sprintf("Cancelled by a Timeout of %s", fo.Name)
				result["fr_status"] = fr.Status_Error
				Response = nil
				funcSend = nil
				headers = nil
				method = ""
				url = ""
				jsonData = nil
				return result
			}
		case <-time.After(time.Duration(fo.SleepSecond) * time.Second):
			{

				code, errResp := funcSend(method, url, headers, jsonData, Response)
				if errResp != nil {
					Response = nil
					funcSend = nil
					headers = nil
					method = ""
					url = ""
					jsonData = nil
					return fo.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
				}
				if code > 299 {
					Response = nil
					funcSend = nil
					headers = nil
					method = ""
					url = ""
					jsonData = nil
					tran.Fiscal.Code = code
					tran.Fiscal.StatusCode = code
					result["fr_id"] = tran.Fiscal.ResiptId
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = Response.Error.Message
				}
				if code == 200 {
					if Response.Data.StatusName == "CONFIRMED" {
						tran.Fiscal.Code = code
						tran.Fiscal.StatusCode = code
						tran.Fiscal.Fields.Fp = Response.Data.Device.FPD
						tran.Fiscal.Fields.Fn = Response.Data.Device.FN
						tran.Fiscal.Fields.Fd = parserTypes.ParseTypeInFloat64(Response.Data.Device.FDN)
						TimeKass, _ := tran.DateTime.ParseDateOfLayout("2006-01-02T15:04:05", Response.Data.ReceiptDateUtc)
						fmt.Printf("TIME OF FISCAL: %s\n", TimeKass)
						tran.Fiscal.Fields.DateFisal = TimeKass
						result["fr_id"] = tran.Fiscal.ResiptId
						result["status"] = fr.Status_Complete
						result["fr_status"] = fr.Status_Complete
						result["f_desc"] = "нет ошибок"
						Response = nil
						funcSend = nil
						headers = nil
						method = ""
						TimeKass = ""
						url = ""
						jsonData = nil
						return result
					}
					if Response.Data.StatusName == "KKT_ERROR" {
						result["fr_id"] = tran.Fiscal.ResiptId
						result["status"] = fr.Status_Error
						result["fr_status"] = fr.Status_Error
						result["f_desc"] = Response.Data.StatusMessage
						Response = nil
						funcSend = nil
						url = ""
						headers = nil
						method = ""
						jsonData = nil
						return result
					}

				}
			}
		}
	}
	return result
}

func (fo *FermaOfd) GetQrPicture(tran *transaction.Transaction) string {
	return fo.Core.EncodeUrlToBase64(tran)
}

func (fo *FermaOfd) GetQrUrl(tran *transaction.Transaction) string {
	return fo.Core.MakeUrlQr(tran)
}

func (fo *FermaOfd) MakeResultError(massage string, err error, result map[string]interface{}) map[string]interface{} {
	fo.ClearResultMap(result)
	result["status"] = transaction.TransactionState_Error
	result["f_desc"] = fmt.Sprintf("%v", err)
	result["message"] = massage
	return result
}

func (fo *FermaOfd) ClearResultMap(result map[string]interface{}) {
	for k := range result {
		delete(result, k)
	}
}
