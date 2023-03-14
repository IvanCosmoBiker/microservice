package nanokass

import (
	config "ephorservices/config"
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
	Name string
	Core *core.Core
	Http *transport.Http
	cfg  *config.Config
}

type NewNanokassaStruct struct {
	Nanokassa
}

func (nos *NewNanokassaStruct) New(conf *config.Config) fr.Fiscal {
	Core := core.InitCore(conf)
	Http := transport.InitHttp(conf)
	return &NewNanokassaStruct{
		Nanokassa: Nanokassa{
			Name: "Nanokassa",
			cfg:  conf,
			Core: Core,
			Http: Http,
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
		return n.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
	}
	log.Printf("%v", Response)
	if code > 299 {
		return n.MakeResultError(fmt.Sprintf("%s - %v", "HTTP code is", code), errors.New(fmt.Sprintf("%s - %v", "HTTP code is", code)), result)
	}
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
	timeout := time.NewTimer(n.cfg.Services.EphorFiscal.ExecuteMinutes * time.Minute)
	headers := make(map[string]interface{})
	for {
		select {
		case <-timeout.C:
			{
				result["status"] = fr.Status_Error
				result["f_desc"] = fmt.Sprintf("Cancelled by a Timeout of %s", n.Name)
				result["fr_status"] = fr.Status_Error
				return result
			}
		case <-time.After(n.cfg.Services.EphorFiscal.SleepSecond * time.Second):
			{
				Response := &responseNanokass.ResponseStatusCheck{}
				funcSend := n.Http.Send(n.Http.Call, "GET", url, headers, []byte(""), Response, 3, 2)
				code, errResp := funcSend("GET", url, headers, []byte(""), Response)
				if errResp != nil {
					return n.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
				}
				if code > 299 {
					result["fr_id"] = tran.Fiscal.ResiptId
					result["state"] = fr.Status_Error
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = Response.Error
				}
				if len(Response.Error) > 0 {
					result["fr_id"] = tran.Fiscal.ResiptId
					result["state"] = fr.Status_Error
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = Response.Error
				}
				if Response.Check_status == 1 || Response.Check_status == 3 {
					tran.Fiscal.Fields.Fp = fmt.Sprintf("%v", Response.Check_num_fp)
					tran.Fiscal.Fields.Fn = parserTypes.ParseTypeInString(Response.Check_fn_num)
					tran.Fiscal.Fields.Fd = parserTypes.ParseTypeInFloat64(Response.Check_num_fd)
					tran.Fiscal.Fields.DateFisal = Response.Check_dt_ofdtime
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
