package modul

import (
	"ephorservices/ephorsale/payment/interface/payment"
	"ephorservices/ephorsale/payment/sbp/modul/core"
	response "ephorservices/ephorsale/payment/sbp/modul/response"
	"ephorservices/ephorsale/payment/sbp/modul/transport"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	logger "ephorservices/pkg/logger"
	"fmt"
	"time"
)

type SbpModul struct {
	State int
	Name  string
	Core  *core.Core
	Http  *transport.Http
	Debug bool
}

type NewSbpModul struct {
	SbpModul
}

func (sm NewSbpModul) New(debug bool) payment.Payment /* тип interfaceBank.Bank*/ {
	Core := core.InitCore()
	Http := transport.InitHttp(debug)
	return &NewSbpModul{
		SbpModul: SbpModul{
			Name:  "SbpModul",
			State: 0,
			Core:  Core,
			Http:  Http,
			Debug: debug,
		},
	}
}

func (sm *SbpModul) GetName() string {
	return sm.Name
}

func (sm *SbpModul) Hold(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	url, data, err := sm.Core.MakeRequestGetQrCode(tran.Payment.SbpPoint, tran)
	if err != nil {
		return sm.setStatusResponseErr(fmt.Sprintf("%v", err), fmt.Sprintf("%v", err))
	}
	header := make(map[string]interface{})
	header["Host"] = core.HOST
	header["Authorization"] = "Bearer " + tran.Payment.Login
	Response := &response.ResponseGetQR{}
	funSend := sm.Http.Send(sm.Http.Call, "POST", url, header, data, Response, 3, 2)
	code, errResp := funSend("POST", url, header, data, Response)
	if errResp != nil {
		return sm.setStatusResponseErr(fmt.Sprintf("%v", errResp), fmt.Sprintf("%v", errResp))
	}
	if code != 200 {
		return sm.setStatusResponseErr(fmt.Sprintf("%v", Response.Massage), fmt.Sprintf("%v", Response.Massage))
	}
	tran.Payment.InvoiceId = Response.Payload
	tran.Payment.OrderId = Response.QrcId
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["orderId"] = Response.QrcId
	ResponsePaymentSystem["invoiceId"] = Response.Payload
	ResponsePaymentSystem["message"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["description"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitWait
	logger.Log.Info("Заказ принят")
	return ResponsePaymentSystem
}

func (sm *SbpModul) checkStatusQr(header map[string]interface{}, url string, tran *transaction.Transaction) (*response.ResponseGetStatus, error) {
	Response := &response.ResponseGetStatus{}
	funSend := sm.Http.Send(sm.Http.Call, "GET", url, header, nil, Response, 3, 2)
	code, errResp := funSend("GET", url, header, nil, Response)
	if errResp != nil {
		return Response, errResp
	}
	if code != 200 {
		return Response, nil
	}
	return Response, nil
}

func (sm *SbpModul) setStatusResponseErr(message, description string) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["message"] = message
	ResponsePaymentSystem["description"] = description
	ResponsePaymentSystem["status"] = transaction.TransactionState_Error
	logger.Log.Errorf("%+v", ResponsePaymentSystem)
	return ResponsePaymentSystem
}

func (sm *SbpModul) Status(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Оплата подтверждена"
	ResponsePaymentSystem["description"] = "Оплата подтверждена"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (sm *SbpModul) Debit(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	timeout := time.NewTimer(2 * time.Minute)
	header := make(map[string]interface{})
	header["Host"] = core.HOST
	header["Authorization"] = "Bearer " + tran.Payment.Login
	params := make(map[string]interface{})
	params["qrId"] = tran.Payment.OrderId
	url, _ := sm.Core.MakeRequestGetStatusQrCode(params)
	for {
		select {
		case <-timeout.C:
			{
				ResponsePaymentSystem["status"] = false
				ResponsePaymentSystem["message"] = "Нет ответа от платёжной системы"
				ResponsePaymentSystem["description"] = "Нет ответа от платёжной системы"
				ResponsePaymentSystem["code"] = transaction.TransactionState_Error
				return ResponsePaymentSystem
			}
		case <-time.After(5 * time.Second):
			{
				Response, err := sm.checkStatusQr(header, url, tran)
				if err != nil {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = fmt.Sprintf("%v", err)
					ResponsePaymentSystem["description"] = fmt.Sprintf("%v", err)
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				}
				switch Response.Status {
				case "NotStarted":
					return sm.setStatusResponseErr("Операции по QR коду не существует", "Операции по QR коду не существует")
				case "Rejected":
					return sm.setStatusResponseErr("Операция отклонена", "Операция отклонена")
				case "TimedOut":
					return sm.setStatusResponseErr("Время ожидания операции превышено", "Время ожидания операции превышено")
				case "Accepted":
					ResponsePaymentSystem["status"] = true
					ResponsePaymentSystem["message"] = "Платёж подтверждён"
					ResponsePaymentSystem["description"] = "Платёж подтверждён"
					tran.Payment.OperationId = Response.OperationId
					ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
					return ResponsePaymentSystem
				}
			}
		}
	}
}

func (sm *SbpModul) Return(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	operationId := tran.Payment.OperationId
	if len(operationId) < 1 {
		ResponsePaymentSystem["message"] = "Нет operationId для возврата"
		ResponsePaymentSystem["description"] = "Нет operationId для возврата"
		ResponsePaymentSystem["status"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	params := make(map[string]interface{})
	params["sbpOperId"] = operationId
	params["type"] = "full"
	url, data, err := sm.Core.MakeRequestReturnMoney(params)
	if err != nil {
		return sm.setStatusResponseErr(fmt.Sprintf("%v", err), fmt.Sprintf("%v", err))
	}
	header := make(map[string]interface{})
	header["Host"] = core.HOST
	header["Authorization"] = "Bearer " + tran.Payment.Login
	Response := &response.ResponseReturnQr{}
	funSend := sm.Http.Send(sm.Http.Call, "POST", url, header, data, Response, 3, 2)
	code, errResp := funSend("POST", url, header, data, Response)
	if errResp != nil {
		return sm.setStatusResponseErr(fmt.Sprintf("%v", errResp), fmt.Sprintf("%v", errResp))
	}
	if code != 200 {
		return sm.setStatusResponseErr(fmt.Sprintf("%v", Response.Massage), fmt.Sprintf("%v", Response.Massage))
	}
	ResponsePaymentSystem["message"] = "Деньги возвращены"
	ResponsePaymentSystem["description"] = "Деньги возвращены"
	ResponsePaymentSystem["status"] = transaction.TransactionState_ReturnMoney
	return ResponsePaymentSystem
}

func (sm *SbpModul) Timeout() {

}
