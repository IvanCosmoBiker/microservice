package orange

import (
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal/interface/fr"
	core "ephorservices/ephorsale/fiscal/orange/core"
	responseOrange "ephorservices/ephorsale/fiscal/orange/response"
	transport "ephorservices/ephorsale/fiscal/orange/transport"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	datetime "ephorservices/pkg/datetime"
	"fmt"
	"log"
	"strings"
	"time"
)

type Orange struct {
	Name string
	Core *core.Core
	Http *transport.Http
	cfg  *config.Config
}

type NewOrangeStruct struct {
	Orange
}

func (nos *NewOrangeStruct) New(conf *config.Config) fr.Fiscal {
	Core := core.InitCore(conf)
	Http := transport.InitHttp(conf)
	return &NewOrangeStruct{
		Orange: Orange{
			Name: "Orange",
			cfg:  conf,
			Core: Core,
			Http: Http,
		},
	}
}

func (o *Orange) SendCheck(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	signature, resiptId, jsonRequest, err := o.Core.MakeRequestSendCheck(tran)
	if err != nil {
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", err)
		return result
	}
	tran.Fiscal.Signature = signature
	url := fmt.Sprintf("https://%v:%v/api/v2/documents/", tran.Fiscal.Config.Dev_addr, tran.Fiscal.Config.Dev_port)
	Response := &responseOrange.ResponseCreateCheck{}
	headers := make(map[string]interface{})
	headers["x-Signature"] = tran.Fiscal.Signature
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Content-Length"] = len(jsonRequest)
	funcSend := o.Http.Send(o.Http.Call, "POST", url, headers, []byte(tran.Fiscal.Config.Auth_public_key), []byte(tran.Fiscal.Config.Auth_private_key), jsonRequest, Response, 3, 2)
	code, errResp := funcSend("POST", url, headers, []byte(tran.Fiscal.Config.Auth_public_key), []byte(tran.Fiscal.Config.Auth_private_key), jsonRequest, Response)
	log.Printf("%v", code)
	if errResp != nil {
		return o.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
	}
	tran.Fiscal.ResiptId = resiptId
	if code == 409 {
		result["status"] = fr.Status_InQueue
		return result
	}
	if code != 201 {
		result["fr_id"] = resiptId
		result["status"] = fr.Status_Error
		result["fr_status"] = fr.Status_Error
		result["f_desc"] = strings.Join(Response.Errors[:], "\n")
		return result
	}
	result["status"] = fr.Status_InQueue
	return result
}

func (o *Orange) GetStatus(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	url := o.Core.MakeRequestStatusQr(tran)
	timeout := time.NewTimer(o.cfg.Services.EphorFiscal.ExecuteMinutes * time.Minute)
	headers := make(map[string]interface{})
	headers["x-Signature"] = tran.Fiscal.Signature
	headers["Content-Type"] = "application/json; charset=utf-8"
	for {
		select {
		case <-timeout.C:
			{
				result["status"] = fr.Status_Error
				result["f_desc"] = fmt.Sprintf("Cancelled by a Timeout of %s", o.Name)
				result["fr_status"] = fr.Status_Error
				return result
			}
		case <-time.After(o.cfg.Services.EphorFiscal.SleepSecond * time.Second):
			{
				Response := &responseOrange.ResponseStatusCheck{}
				funcSend := o.Http.Send(o.Http.Call, "GET", url, headers, []byte(tran.Fiscal.Config.Auth_public_key), []byte(tran.Fiscal.Config.Auth_private_key), []byte(""), Response, 3, 2)
				code, errResp := funcSend("GET", url, headers, []byte(tran.Fiscal.Config.Auth_public_key), []byte(tran.Fiscal.Config.Auth_private_key), []byte(""), Response)
				if errResp != nil {
					return o.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
				}
				log.Printf("%v", code)
				if code > 299 {
					result["fr_id"] = tran.Fiscal.ResiptId
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = strings.Join(Response.Errors[:], "\n")
					return result
				}
				if code == 200 {
					date, _ := datetime.Init()
					TimeOrange, _ := date.ParseDateOfLayout("2006-01-02T15:04:05", Response.ProcessedAt)
					TimeUtcOrange, _ := date.SubtractFromTime(TimeOrange, 10800)
					tran.Fiscal.Fields.Fp = Response.Fp
					tran.Fiscal.Fields.Fn = Response.FsNumber
					tran.Fiscal.Fields.Fd = float64(Response.DocumentNumber)
					tran.Fiscal.Fields.DateFisal = TimeUtcOrange
					result["fr_id"] = tran.Fiscal.ResiptId
					result["status"] = fr.Status_Complete
					result["fr_status"] = fr.Status_Complete
					result["f_desc"] = strings.Join(Response.Errors[:], "\n")
					return result
				}
			}
		}
	}
	return result
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
