package fermaOfd

import (
	"data/transaction"
	responseOfd "ephorfiscal/fiscal/fermaOfd/response"
	"ephorfiscal/fiscal/interface/fr"
	pb "ephorfiscal/service"
	config "ephorservices/config"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"regexp"
	"time"
)

type FermaOfd struct {
	Name string
	Core *Core
	Http *Http
	cfg  *config.Config
}

type NewFermaOfdStruct struct {
	FermaOfd
}

func (nos *NewFermaOfdStruct) New(conf *config.Config) fr.Fiscal {
	Core := InitCore(conf)
	Http := InitHttp(conf)
	return &NewFermaOfdStruct{
		FermaOfd: FermaOfd{
			Name: "FermaOfd",
			cfg:  conf,
			Core: Core,
			Http: Http,
		},
	}
}

func (fo *FermaOfd) MakeAuth(req *pb.Request) map[string]interface{} {
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
	funcSend := fo.Http.Send(fo.Http.Call, "POST", url, headers, jsonRequest, Response, 3, 2)
	code, errResp := funcSend("POST", url, headers, jsonRequest, Response)
	if errResp != nil {
		return fo.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
	}
	if code != 200 {
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", Response.Error.Message)
		return result
	}
	tran.Fiscal.AuthToken = Response.Data.AuthToken
	if len(tran.Fiscal.AuthToken) == 0 {
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("Error: %s No auth token", fo.Name)
		return result
	}
	result["status"] = fr.Status_InQueue
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
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", err)
		return result
	}
	url := fmt.Sprintf("https://%s/api/kkt/cloud/receipt?AuthToken=%s", tran.Fiscal.Config.Dev_addr, tran.Fiscal.AuthToken)
	Response := &responseOfd.ResponseSendCheck{}
	headers := make(map[string]interface{})
	funcSend := fo.Http.Send(fo.Http.Call, "POST", url, headers, jsonRequest, Response, 3, 2)
	code, errResp := funcSend("POST", url, headers, jsonRequest, Response)
	if errResp != nil {
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
		}
		result["status"] = fr.Status_Error
		return result
	}
	tran.Fiscal.ResiptId = Response.Data.ReceiptId
	if len(tran.Fiscal.ResiptId) == 0 {
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("Error: %s No ReceiptId", fo.Name)
		return result
	}
	result["status"] = fr.Status_InQueue
	return result
}

func (fo *FermaOfd) GetStatus(tran *transaction.Transaction) map[string]interface{} {
	result := make(map[string]interface{})
	jsonData, err := fo.Core.MakeRequestStatusCheck(tran)
	if err != nil {
		result["status"] = fr.Status_Error
		result["f_desc"] = fmt.Sprintf("%v", err)
		return result
	}
	timeout := time.NewTimer(fo.cfg.Services.EphorFiscal.ExecuteMinutes * time.Minute)
	headers := make(map[string]interface{})
	url := fmt.Sprintf("https://%s/api/kkt/cloud/status?AuthToken=%s", tran.Fiscal.Config.Dev_addr, tran.Fiscal.AuthToken)
	for {
		select {
		case <-timeout.C:
			{
				result["status"] = fr.Status_Error
				result["f_desc"] = fmt.Sprintf("Cancelled by a Timeout of %s", fo.Name)
				result["fr_status"] = fr.Status_Error
				return result
			}
		case <-time.After(fo.cfg.Services.EphorFiscal.SleepSecond * time.Second):
			{
				Response := &responseOfd.ResponseStatusCheck{}
				funcSend := fo.Http.Send(fo.Http.Call, "POST", url, headers, jsonData, Response, 3, 2)
				code, errResp := funcSend("POST", url, headers, jsonData, Response)
				if errResp != nil {
					return fo.MakeResultError(fmt.Sprintf("%v", errResp), errResp, result)
				}
				if code > 299 {
					result["fr_id"] = tran.Fiscal.ResiptId
					result["status"] = fr.Status_Error
					result["fr_status"] = fr.Status_Error
					result["f_desc"] = Response.Error.Message
				}
				if code == 200 {
					tran.Fiscal.Fields.Fp = Response.Data.Device.FPD
					tran.Fiscal.Fields.Fn = Response.Data.Device.FN
					tran.Fiscal.Fields.Fd = parserTypes.ParseTypeInFloat64(Response.Data.Device.FDN)
					tran.Fiscal.Fields.DateFisal = Response.Data.ReceiptDateUtc
					result["fr_id"] = tran.Fiscal.ResiptId
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
