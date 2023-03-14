package sberpay

import (
	"ephorservices/ephorsale/payment/interface/payment"
	"ephorservices/ephorsale/payment/sberpay/core"
	response "ephorservices/ephorsale/payment/sberpay/response"
	"ephorservices/ephorsale/payment/sberpay/transport"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"log"
	"time"
)

type SberPay struct {
	Name                     string
	Counter                  int
	State                    int
	Core                     *core.Core
	Http                     *transport.Http
	UrlCreateOrderSberPay    string
	UrlGetStatusOrderSberPay string
	UrlDepositeOrderSberPay  string
	UrlReverseSberPay        string
	UrlRefundSberPay         string
	Debug                    bool
}

type NewSberPayStruct struct {
	SberPay
}

func (sberPay NewSberPayStruct) New(debug bool) payment.Payment /* тип Payment*/ {
	Core := core.InitCore()
	Http := transport.InitHttp(debug)
	return &NewSberPayStruct{
		SberPay: SberPay{
			Name:                     "SberPay",
			Counter:                  0,
			State:                    0,
			Core:                     Core,
			Http:                     Http,
			UrlCreateOrderSberPay:    "https://securepayments.sberbank.ru/payment/rest/registerPreAuth.do",
			UrlGetStatusOrderSberPay: "https://securepayments.sberbank.ru/payment/rest/getOrderStatusExtended.do",
			UrlDepositeOrderSberPay:  "https://securepayments.sberbank.ru/payment/rest/deposit.do",
			UrlReverseSberPay:        "https://securepayments.sberbank.ru//payment/rest/reverse.do",
			UrlRefundSberPay:         "https://securepayments.sberbank.ru//payment/rest/refund.do",
			Debug:                    debug,
		},
	}
}

func (sb *SberPay) GetName() string {
	return sb.Name
}

func (sb *SberPay) Hold(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	stringCreateOrder, err := sb.Core.MakeRequestCreateOrder(tran)
	if err != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%v", err)
		ResponsePaymentSystem["description"] = fmt.Sprintf("%v", err)
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	Response := &response.ResponseCreateOrder{}
	url := fmt.Sprintf("%s?%s", sb.UrlCreateOrderSberPay, stringCreateOrder)
	header := make(map[string]interface{})
	funcSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 30)
	StatusCode, errResp := funcSend("POST", url, header, nil, Response)
	if errResp != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	if StatusCode != 200 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = Response.ErrorMessage
		ResponsePaymentSystem["description"] = Response.ErrorMessage
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	if len(Response.ErrorMessage) > 0 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = Response.ErrorMessage
		ResponsePaymentSystem["description"] = Response.ErrorMessage
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	log.Printf("%+v", Response)
	if Response.ExternalParams.SbolInactive != false {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = "Банк недоступен"
		ResponsePaymentSystem["description"] = "Банк недоступен"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["description"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyHoldWait
	ResponsePaymentSystem["orderId"] = Response.OrderId
	ResponsePaymentSystem["invoiceId"] = Response.ExternalParams.SbolBankInvoiceId
	tran.Payment.SbolBankInvoiceId = Response.ExternalParams.SbolBankInvoiceId
	tran.Payment.OrderId = Response.OrderId
	if sb.Debug {
		log.Printf(tran.Payment.SbolBankInvoiceId)
		log.Printf(tran.Payment.OrderId)
	}
	return ResponsePaymentSystem
}

func (sb *SberPay) Status(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = false
	ResponsePaymentSystem["message"] = "Нет ответа от банка"
	ResponsePaymentSystem["description"] = "Нет ответа от банка"
	requestStatus, err := sb.Core.MakeRequestStatusOrder(tran)
	if err != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%s", err)
		ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Опрос статуса sberpay транзакции)"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	url := fmt.Sprintf("%s?%s", sb.UrlGetStatusOrderSberPay, requestStatus)
	header := make(map[string]interface{})
	timeout := time.NewTimer(1 * time.Minute)
	for {
		select {
		case <-timeout.C:
			{
				ResponsePaymentSystem["status"] = false
				ResponsePaymentSystem["message"] = "Нет ответа от банка"
				ResponsePaymentSystem["description"] = "Нет ответа от банка"
				ResponsePaymentSystem["orderId"] = tran.Payment.OrderId
				ResponsePaymentSystem["invoiceId"] = tran.Payment.SbolBankInvoiceId
				ResponsePaymentSystem["code"] = transaction.TransactionState_Error
				return ResponsePaymentSystem
			}
		case <-time.After(5 * time.Second):
			{
				Response := &response.ResponseStatusOrder{}
				funcSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 30)
				StatusCode, errResp := funcSend("POST", url, header, nil, Response)
				log.Printf("%+v", Response)
				if errResp != nil {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
					ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				}
				if StatusCode != 200 {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = Response.ErrorMessage
					ResponsePaymentSystem["description"] = Response.ErrorMessage
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				}
				errCode := parserTypes.ParseTypeInterfaceToInt(Response.ErrorCode)
				OrderStatus := parserTypes.ParseTypeInterfaceToInt(Response.OrderStatus)
				if errCode != 0 {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = Response.ErrorMessage
					ResponsePaymentSystem["description"] = Response.ErrorMessage
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				}
				if OrderStatus == 3 {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = "Авторизация отменена"
					ResponsePaymentSystem["description"] = "Авторизация отменена"
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				} else if OrderStatus == 4 {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = "По транзакции была проведена операция возврата"
					ResponsePaymentSystem["description"] = "По транзакции была проведена операция возврата"
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				} else if OrderStatus == 6 {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = "Авторизация отклонена"
					ResponsePaymentSystem["description"] = "Авторизация отклонена"
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				}
				if OrderStatus == 1 || OrderStatus == 2 {
					tran.Status = transaction.TransactionState_MoneyDebitWait
					ResponsePaymentSystem["status"] = true
					ResponsePaymentSystem["message"] = "Платёж подтверждён"
					ResponsePaymentSystem["description"] = "Платёж подтверждён"
					ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitWait
					return ResponsePaymentSystem
				}

			}
		}
	}
	return ResponsePaymentSystem
}

func (sb *SberPay) Debit(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	if tran.Payment.DebitSum == 0 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = "не задана сумма для списания"
		ResponsePaymentSystem["description"] = "не задана сумма для списания"
		return ResponsePaymentSystem
	}
	dataPush, err := sb.Core.MakeRequestDepositOrder(tran)
	if err != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%s", err)
		ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Списание денег транзакции)"
		return ResponsePaymentSystem
	}
	Response := &response.ResponseDebitOrder{}
	header := make(map[string]interface{})
	url := fmt.Sprintf("%s?%s", sb.UrlDepositeOrderSberPay, dataPush)
	funcSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 30)
	StatusCode, errResp := funcSend("POST", url, header, nil, Response)
	if errResp != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
		return ResponsePaymentSystem
	}
	if StatusCode != 200 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = Response.ErrorMessage
		ResponsePaymentSystem["description"] = Response.ErrorMessage
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	errCode := parserTypes.ParseTypeInterfaceToInt(Response.ErrorCode)
	if errCode != 0 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = Response.ErrorMessage
		ResponsePaymentSystem["description"] = Response.ErrorMessage
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	if errCode == 7 {
		return sb.Debit(tran)
	}
	tran.Status = transaction.TransactionState_MoneyDebitOk
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги списаны"
	ResponsePaymentSystem["description"] = "Деньги списаны"
	return ResponsePaymentSystem
}

func (sb *SberPay) Return(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	stringRequest, _ := sb.Core.MakeRequestReturnMoney(tran)
	url := ""
	if tran.Status == transaction.TransactionState_MoneyDebitWait {
		url = fmt.Sprintf("%s?%s", sb.UrlReverseSberPay, stringRequest)
	}
	if tran.Status == transaction.TransactionState_MoneyDebitOk {
		url = fmt.Sprintf("%s?%s", sb.UrlRefundSberPay, stringRequest)
	}
	Response := &response.ResponseReturnOrder{}
	header := make(map[string]interface{})
	funSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 30)
	StatusCode, errResp := funSend("POST", url, header, nil, Response)
	log.Printf("%+v", Response)
	if errResp != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	if StatusCode != 200 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = Response.ErrorMessage
		ResponsePaymentSystem["description"] = Response.ErrorMessage
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	errCode := parserTypes.ParseTypeInterfaceToInt(Response.ErrorCode)
	if errCode != 0 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = Response.ErrorMessage
		ResponsePaymentSystem["description"] = Response.ErrorMessage
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги возвращены"
	ResponsePaymentSystem["description"] = "Деньги возвращены"
	ResponsePaymentSystem["code"] = transaction.TransactionState_ReturnMoney
	return ResponsePaymentSystem
}

func (sb *SberPay) Timeout() {

}
