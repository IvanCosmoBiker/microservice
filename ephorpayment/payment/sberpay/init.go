package sberpay

import (
	config "ephorservices/config"
	"ephorservices/ephorpayment/payment/interface/payment"
	core "ephorservices/ephorpayment/payment/sberpay/core"
	response "ephorservices/ephorpayment/payment/sberpay/response"
	transport "ephorservices/ephorpayment/payment/sberpay/transport"
	pb "ephorservices/ephorpayment/service"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"log"
	"time"
)

type SberPay struct {
	Name                     string
	Counter                  int
	Status                   int
	Core                     *core.Core
	Http                     *transport.Http
	UrlCreateOrderSberPay    string
	UrlGetStatusOrderSberPay string
	UrlDepositeOrderSberPay  string
	UrlReverseSberPay        string
	UrlRefundSberPay         string
	cfg                      *config.Config
}

type NewSberPayStruct struct {
	SberPay
}

func (sberPay NewSberPayStruct) New(conf *config.Config) payment.Payment /* тип Payment*/ {
	Core := core.InitCore(conf)
	Http := transport.InitHttp(conf)
	return &NewSberPayStruct{
		SberPay: SberPay{
			Name:                     "SberPay",
			Counter:                  0,
			Status:                   0,
			Core:                     Core,
			Http:                     Http,
			UrlCreateOrderSberPay:    "https://securepayments.sberbank.ru/payment/rest/registerPreAuth.do",
			UrlGetStatusOrderSberPay: "https://securepayments.sberbank.ru/payment/rest/getOrderStatusExtended.do",
			UrlDepositeOrderSberPay:  "https://securepayments.sberbank.ru/payment/rest/deposit.do",
			UrlReverseSberPay:        "https://securepayments.sberbank.ru//payment/rest/reverse.do",
			UrlRefundSberPay:         "https://securepayments.sberbank.ru//payment/rest/refund.do",
			cfg:                      conf,
		},
	}
}

func (sb *SberPay) HoldMoney(tran *pb.Request) map[string]interface{} {
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
	funcSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 2)
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
	tran.SbolBankInvoiceId = Response.ExternalParams.SbolBankInvoiceId
	tran.OrderId = Response.OrderId
	if sb.cfg.Debug {
		log.Printf(tran.SbolBankInvoiceId)
		log.Printf(tran.OrderId)
	}
	return ResponsePaymentSystem
}

func (sb *SberPay) GetStatusHoldMoney(tran *pb.Request) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
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
	timeout := time.NewTimer(sb.cfg.Services.EphorPayment.Config.ExecuteMinutes * time.Minute)
	for {
		select {
		case <-timeout.C:
			{
				ResponsePaymentSystem["status"] = false
				ResponsePaymentSystem["message"] = "Нет ответа от банка"
				ResponsePaymentSystem["description"] = "Нет ответа от банка"
				ResponsePaymentSystem["orderId"] = tran.OrderId
				ResponsePaymentSystem["invoiceId"] = tran.SbolBankInvoiceId
				ResponsePaymentSystem["code"] = transaction.TransactionState_Error
				return ResponsePaymentSystem
			}
		case <-time.After(sb.cfg.Services.EphorPayment.Config.IntervalTime * time.Second):
			{
				Response := &response.ResponseStatusOrder{}
				funcSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 2)
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
					tran.Status = int32(transaction.TransactionState_MoneyDebitWait)
					ResponsePaymentSystem["status"] = true
					ResponsePaymentSystem["message"] = "Платёж подтверждён"
					ResponsePaymentSystem["description"] = "Платёж подтверждён"
					ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitWait
					return ResponsePaymentSystem
				}

			}
		}
	}
}

func (sb *SberPay) DebitHoldMoney(tran *pb.Request) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	if tran.DebitSum == 0 {
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
	funcSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 2)
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
		return sb.DebitHoldMoney(tran)
	}
	tran.Status = int32(transaction.TransactionState_MoneyDebitOk)
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги списаны"
	ResponsePaymentSystem["description"] = "Деньги списаны"
	return ResponsePaymentSystem
}

func (sb *SberPay) ReturnMoney(tran *pb.Request) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	stringRequest, _ := sb.Core.MakeRequestReturnMoney(tran)
	url := ""
	if tran.Status == int32(transaction.TransactionState_MoneyDebitWait) {
		url = fmt.Sprintf("%s?%s", sb.UrlReverseSberPay, stringRequest)
	}
	if tran.Status == int32(transaction.TransactionState_MoneyDebitOk) {
		url = fmt.Sprintf("%s?%s", sb.UrlRefundSberPay, stringRequest)
	}
	Response := &response.ResponseReturnOrder{}
	header := make(map[string]interface{})
	funSend := sb.Http.Send(sb.Http.Call, "POST", url, header, nil, Response, 3, 2)
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
