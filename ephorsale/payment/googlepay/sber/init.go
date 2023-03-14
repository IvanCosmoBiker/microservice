package sber

import (
	"encoding/json"
	config "ephorservices/config"
	"ephorservices/ephorsale/payment/interface/payment"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"log"
	"time"
)

type SberGooglePay struct {
	Name              string
	Counter           int
	Status            int
	cfg               *config.Config
	Core              *Core
	Http              *Http
	UrlCreateOrder    string
	UrlGetStatusOrder string
	UrlDepositeOrder  string
	UrlReverse        string
}

type NewSberGooglePayStruct struct {
	SberGooglePay
}

func (sberGPay NewSberGooglePayStruct) New(conf *config.Config) payment.Payment /* тип Payment*/ {
	Core := InitCore(conf)
	Http := InitHttp(conf)
	return &NewSberGooglePayStruct{
		SberGooglePay: SberGooglePay{
			Name:              "SberGooglePay",
			Counter:           0,
			Status:            0,
			Core:              Core,
			Http:              Http,
			UrlCreateOrder:    "https://3dsec.sberbank.ru/payment/google/payment.do",
			UrlGetStatusOrder: "https://3dsec.sberbank.ru/payment/google/getOrderStatusExtended.do",
			UrlDepositeOrder:  "https://3dsec.sberbank.ru/payment/google/deposit.do",
			UrlReverse:        "https://3dsec.sberbank.ru/payment/google/reverse.do",
			cfg:               conf,
		},
	}
}

func (sgp *SberGooglePay) HoldMoney(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	dataPush, err := sgp.Core.makeOrderRequestCreateOrder(tran)
	if sgp.cfg.Debug {
		log.Printf("%+v", dataPush)
		log.Printf("%+v", tran)
	}
	if err != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%sgp", err)
		ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Создание транзакции)"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	header := make(map[string]interface{})
	funcSend := sgp.Http.Send(sgp.Http.Call, "POST", sgp.UrlCreateOrder, header, dataPush, 3, 2)
	Response, errResp := funcSend("POST", sgp.UrlCreateOrder, header, dataPush)
	if errResp != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	if Response.StatusCode != 200 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = "Сервис банка Сбер не отвечает"
		ResponsePaymentSystem["description"] = "Сервис банка Сбер не отвечает"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	success := Response.GetData("success").(bool)
	if success != true {
		errorResp := Response.Data["error"].(map[string]interface{})
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = errorResp["message"]
		ResponsePaymentSystem["description"] = errorResp["description"]
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["description"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyHoldWait
	tran.Payment.OrderId = Response.Data["orderId"].(string)
	return ResponsePaymentSystem
}

func (sgp *SberGooglePay) DebitHoldMoney(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})

	dataPush, err := sgp.Core.makeRequestDepositOrder(tran)
	if err != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%sgp", err)
		ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Списание денег транзакции)"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	header := make(map[string]interface{})
	funcSend := sgp.Http.Send(sgp.Http.Call, "POST", sgp.UrlDepositeOrder, header, dataPush, 3, 2)
	Response, errResp := funcSend("POST", sgp.UrlDepositeOrder, header, dataPush)
	if errResp != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	if Response.StatusCode != 200 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = "Сервис банка Сбер не отвечает"
		ResponsePaymentSystem["description"] = "Сервис банка Сбер не отвечает"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	success := Response.GetData("success").(bool)
	if success != true {
		errorResp := Response.Data["error"].(map[string]interface{})
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = errorResp["message"]
		ResponsePaymentSystem["description"] = errorResp["description"]
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги списаны"
	ResponsePaymentSystem["description"] = "Деньги списаны"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	return ResponsePaymentSystem
}

func (sgp *SberGooglePay) GetStatusHoldMoney(tran *transaction.Transaction) map[string]interface{} {
	return sgp.getStatus(tran)
}

func (sgp *SberGooglePay) getStatus(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	dataPush, err := sgp.Core.makeRequestStatusOrder(tran)
	if err != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%s", err)
		ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Опрос статуса транзакции)"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	header := make(map[string]interface{})
	timer := time.NewTimer(2 * time.Minute)
	for {
		select {
		case <-timer.C:
			{
				ResponsePaymentSystem["status"] = false
				ResponsePaymentSystem["message"] = "Нет ответа от банка"
				ResponsePaymentSystem["description"] = "Нет ответа от банка"
				ResponsePaymentSystem["code"] = transaction.TransactionState_Error
				return ResponsePaymentSystem
			}
		case <-time.After(5 * time.Second):
			{
				funcSend := sgp.Http.Send(sgp.Http.Call, "POST", sgp.UrlGetStatusOrder, header, dataPush, 3, 2)
				Response, errResp := funcSend("POST", sgp.UrlGetStatusOrder, header, dataPush)
				if errResp != nil {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
					ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
					ResponsePaymentSystem["code"] = transaction.TransactionState_Error
					return ResponsePaymentSystem
				}
				errCode := parserTypes.ParseTypeInterfaceToInt(Response.GetData("errorCode"))
				OrderStatus := parserTypes.ParseTypeInterfaceToInt(Response.GetData("orderStatus"))
				if errCode != 0 {
					ResponsePaymentSystem["status"] = false
					ResponsePaymentSystem["message"] = Response.GetData("errorMessage").(string)
					ResponsePaymentSystem["description"] = Response.GetData("errorMessage").(string)
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
				ResponsePaymentSystem["status"] = true
				ResponsePaymentSystem["message"] = "Платёж подтверждён"
				ResponsePaymentSystem["description"] = "Платёж подтверждён"
				ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitWait
				return ResponsePaymentSystem

			}
		}
	}

}

func (sgp *SberGooglePay) ReturnMoney(tran *transaction.Transaction) map[string]interface{} {
	ResponsePaymentSystem := make(map[string]interface{})
	requestOrder := make(map[string]interface{})
	requestOrder["orderId"] = tran.Payment.OrderId
	data, err := json.Marshal(requestOrder)
	if err != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%sgp", err)
		ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Возврат денег)"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	header := make(map[string]interface{})
	funcSend := sgp.Http.Send(sgp.Http.Call, "POST", sgp.UrlReverse, header, data, 3, 2)
	Response, errResp := funcSend("POST", sgp.UrlReverse, header, data)
	if errResp != nil {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["description"] = fmt.Sprintf("%v", errResp)
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	if Response.StatusCode != 200 {
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = "Сервис банка Сбер не отвечает"
		ResponsePaymentSystem["description"] = "Сервис банка Сбер не отвечает"
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	success := Response.GetData("success").(bool)
	if success != true {
		errorResp := Response.Data["error"].(map[string]interface{})
		ResponsePaymentSystem["status"] = false
		ResponsePaymentSystem["message"] = errorResp["message"]
		ResponsePaymentSystem["description"] = errorResp["description"]
		ResponsePaymentSystem["code"] = transaction.TransactionState_Error
		return ResponsePaymentSystem
	}
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Деньги возвращены"
	ResponsePaymentSystem["description"] = "Деньги возвращены"
	ResponsePaymentSystem["code"] = transaction.TransactionState_ReturnMoney
	return ResponsePaymentSystem
}

func (sgp *SberGooglePay) Timeout() {

}
