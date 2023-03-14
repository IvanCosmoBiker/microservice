package nanokass

import (
	"ephorservices/ephorsale/fiscal/interface/fr"
	core "ephorservices/ephorsale/fiscal/nanokass/core"
	responseNanokass "ephorservices/ephorsale/fiscal/nanokass/response"
	transport "ephorservices/ephorsale/fiscal/nanokass/transport"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Nanokassa struct {
	Name           string
	ExecuteMinutes int
	SleepSecond    int
	Core           *core.Core
	Http           *transport.Http
}

type NewNanokassaStruct struct {
	Nanokassa
}

func (nos *NewNanokassaStruct) New(executeMinutes, sleepSecond int, debug bool) fr.Fiscal {
	Core := core.InitCore()
	Http := transport.InitHttp(debug)
	return &NewNanokassaStruct{
		Nanokassa: Nanokassa{
			Name:           "Nanokassa",
			ExecuteMinutes: executeMinutes,
			SleepSecond:    sleepSecond,
			Core:           Core,
			Http:           Http,
		},
	}
}

func (n *Nanokassa) SendCheck(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	jsonRequest, url, err := n.Core.MakeRequestSendCheck(tran)
	if err != nil {
		return n.MakeResultError(fmt.Sprintf("%v", err), err, result)
	}
	Response := &responseNanokass.ResponseSendCheck{}
	headers := make(map[string]interface{})
	funcSend := n.Http.Send(n.Http.Call, "POST", url, headers, jsonRequest, Response, 3, 2)
	code, errResp := funcSend("POST", url, headers, jsonRequest, Response)
	log.Printf("%v", Response)
	if errResp != nil {
		tran.Fiscal.Code = 400
		tran.Fiscal.StatusCode = 400
		return n.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
	}
	log.Printf("%v", Response)
	if code > 299 {
		tran.Fiscal.Code = code
		tran.Fiscal.StatusCode = code
		return n.MakeResultError(fmt.Sprintf("%s - %v", "HTTP code is", code), errors.New(fmt.Sprintf("%s - %v", "HTTP code is", code)), result)
	}
	tran.Fiscal.Code = code
	tran.Fiscal.StatusCode = code
	receiptId := make([]string, 0)
	log.Printf("%+v", Response)
	receiptId = append(receiptId, Response.Nuid, Response.Qnuid)
	log.Printf("%+v", receiptId)
	tran.Fiscal.ResiptId = strings.Join(receiptId, ":")
	log.Printf("%+v", tran.Fiscal.ResiptId)
	result["status"] = fr.Status_InQueue
	result["fr_status"] = true
	return result
}

func (n *Nanokassa) GetStatus(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	url := n.Core.MakeRequestStatusCheck(tran)
	timeout := time.NewTimer(time.Duration(n.ExecuteMinutes) * time.Minute)
	headers := make(map[string]interface{})
	Response := &responseNanokass.ResponseStatusCheck{}
	for {
		select {
		case <-timeout.C:
			{
				tran.Fiscal.Code = 400
				tran.Fiscal.StatusCode = 400
				result["status"] = fr.Status_Error
				result["f_desc"] = fmt.Sprintf("Cancelled by a Timeout of %s", n.Name)
				result["fr_status"] = fr.Status_Error
				return result
			}
		case <-time.After(time.Duration(n.SleepSecond) * time.Second):
			{

				funcSend := n.Http.Send(n.Http.Call, "GET", url, headers, []byte(""), Response, 3, 2)
				code, errResp := funcSend("GET", url, headers, []byte(""), Response)
				if errResp != nil {
					tran.Fiscal.Code = code
					tran.Fiscal.StatusCode = code
					return n.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
				}
				if code > 299 {
					tran.Fiscal.Code = code
					tran.Fiscal.StatusCode = code
					result["fr_id"] = tran.Fiscal.ResiptId
					result["state"] = fr.Status_Error
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = Response.Error
				}
				if len(Response.Error) > 0 {
					tran.Fiscal.Code = 400
					tran.Fiscal.StatusCode = 400
					result["fr_id"] = tran.Fiscal.ResiptId
					result["state"] = fr.Status_Error
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = Response.Error
				}
				if Response.Check_status == 1 || Response.Check_status == 3 {
					tran.Fiscal.Code = code
					tran.Fiscal.StatusCode = code
					tran.Fiscal.Fields.Fp = fmt.Sprintf("%v", Response.Check_num_fp)
					tran.Fiscal.Fields.Fn = parserTypes.ParseTypeInString(Response.Check_fn_num)
					tran.Fiscal.Fields.Fd = parserTypes.ParseTypeInFloat64(Response.Check_num_fd)
					TimeKass, _ := tran.DateTime.ParseDateOfLayout("2006-01-02T15:04:05", Response.Check_dt_ofdtime)
					tran.Fiscal.Fields.DateFisal = TimeKass
					result["fr_id"] = tran.Fiscal.ResiptId
					result["state"] = fr.Status_Complete
					result["status"] = fr.Status_Complete
					result["fr_status"] = fr.Status_Complete
					result["f_desc"] = "нет ошибок"
					return result
				}
			}
		}
	}
	return result
}

func (n *Nanokassa) GetQrPicture(tran *transaction.Transaction) string {
	return n.Core.EncodeUrlToBase64(tran)
}

func (n *Nanokassa) GetQrUrl(tran *transaction.Transaction) string {
	return n.Core.MakeUrlQr(tran)
}

func (n *Nanokassa) MakeResultError(massage string, err error, result map[string]interface{}) map[string]interface{} {
	n.ClearResultMap(result)
	result["status"] = transaction.TransactionState_Error
	result["fr_status"] = false
	result["state"] = fr.Status_Error
	result["f_desc"] = fmt.Sprintf("%v", err)
	result["message"] = massage
	return result
}

func (n *Nanokassa) ClearResultMap(result map[string]interface{}) {
	for k := range result {
		delete(result, k)
	}
}
