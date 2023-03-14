package manager

import (
	"context"
	config "ephorservices/config"
	payment "ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/payment/interface/manager"
	interfacePayment "ephorservices/ephorsale/payment/interface/payment"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	MockPaymentSystem "ephorservices/test/mock/payment/paymentSystem"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var Config *config.Config
var TransactionTest *transaction.Transaction
var Manager manager.ManagerPayment

func init() {
	newConfig := &config.Config{}
	newConfig.Transport.Mqtt.Login = "device"
	newConfig.Transport.Mqtt.Address = "188.225.18.140"
	newConfig.Transport.Mqtt.Password = "ephor2021"
	newConfig.Transport.Mqtt.Port = "1883"
	newConfig.Transport.Mqtt.ExecuteTimeSeconds = 30
	newConfig.Transport.Mqtt.BackOffPolicySendMassage = make([]time.Duration, 3)
	newConfig.Transport.Mqtt.BackOffPolicyConnection = make([]time.Duration, 3)
	Config = newConfig
	tran := transaction.Transaction{}
	TransactionTest = &tran
	var ctx context.Context
	Manager = payment.Init(Config, ctx)

}

func makeMockPaymentsSystem(t *testing.T) {
	var Payments = make(map[uint8]interfacePayment.Payment, len(payment.TypePayment))
	for _, item := range payment.TypePayment {
		controller := gomock.NewController(t)
		defer controller.Finish()
		MockInterface := MockPaymentSystem.NewMockPayment(controller)
		Payments[item] = MockInterface
	}
	Manager.SetPayment(Payments)
}

func TestMANAGER_StartPaymentNoAvalableBank(t *testing.T) {
	var mapNoAvalableBank = make(map[string]interface{})
	mapNoAvalableBank["error"] = "No Avalable Type Payment System"
	mapNoAvalableBank["ps_desc"] = "No Avalable Type Payment System"
	mapNoAvalableBank["ps_order"] = "none"
	mapNoAvalableBank["status"] = transaction.TransactionState_Error
	result := Manager.Payment(TransactionTest)
	assert.Equal(t, mapNoAvalableBank, result)
}

func TestMANAGER_StartPaymentAvalableBank(t *testing.T) {
	makeMockPaymentsSystem(t)
	TransactionTest.Payment.Type = interfacePayment.TypeVendPay
	var mapAvalableBank = make(map[string]interface{})
	mapAvalableBank["status"] = transaction.TransactionState_MoneyHoldStart
	result := Manager.Payment(TransactionTest)
	assert.Equal(t, result, mapAvalableBank)
}

func TestMANAGER_HoldMoneyErr(t *testing.T) {
	makeMockPaymentsSystem(t)
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = false
	ResponsePaymentSystem["message"] = "err"
	ResponsePaymentSystem["description"] = "err"
	ResponsePaymentSystem["code"] = transaction.TransactionState_Error
	paymentSystem, _ := Manager.GetPaymentOfType(interfacePayment.TypeSber)
	paymentSystem.(*MockPaymentSystem.MockPayment).EXPECT().HoldMoney(TransactionTest).Return(ResponsePaymentSystem)
	TransactionTest.Payment.Type = interfacePayment.TypeSber
	mapErrHoldMoney := make(map[string]interface{})
	mapErrHoldMoney["status"] = transaction.TransactionState_Error
	mapErrHoldMoney["ps_order"] = "none"
	mapErrHoldMoney["ps_desc"] = "err"
	mapErrHoldMoney["error"] = "err"
	resultMock := Manager.Hold(TransactionTest)
	assert.Equal(t, mapErrHoldMoney, resultMock)
}

func TestMANAGER_HoldMoneyOk(t *testing.T) {
	makeMockPaymentsSystem(t)
	ResponsePaymentSystem := make(map[string]interface{})
	ResponsePaymentSystem["status"] = true
	ResponsePaymentSystem["message"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["description"] = "Заказ принят, ожидание оплаты"
	ResponsePaymentSystem["code"] = transaction.TransactionState_MoneyDebitOk
	ResponsePaymentSystem["invoiceId"] = "test"
	ResponsePaymentSystem["orderId"] = "test"
	paymentSystem, _ := Manager.GetPaymentOfType(interfacePayment.TypeSber)
	paymentSystem.(*MockPaymentSystem.MockPayment).EXPECT().HoldMoney(TransactionTest).Return(ResponsePaymentSystem)
	TransactionTest.Payment.Type = interfacePayment.TypeSber
	mapOkHoldMoney := make(map[string]interface{})
	mapOkHoldMoney["status"] = transaction.TransactionState_MoneyDebitOk
	mapOkHoldMoney["ps_order"] = "test"
	mapOkHoldMoney["ps_invoice_id"] = "test"
	mapOkHoldMoney["ps_tid"] = nil
	mapOkHoldMoney["ps_desc"] = "Заказ принят, ожидание оплаты"
	mapOkHoldMoney["error"] = ""
	resultMock := Manager.Hold(TransactionTest)
	assert.Equal(t, mapOkHoldMoney, resultMock)
}

func TestMANAGER_HoldMoneyEmptyMap(t *testing.T) {
	var Payments = make(map[uint8]interfacePayment.Payment, 0)
	Manager.SetPayment(Payments)
	mapNoAvalableBank := make(map[string]interface{})
	mapNoAvalableBank["error"] = "No Avalable Type Payment System"
	mapNoAvalableBank["ps_desc"] = "No Avalable Type Payment System"
	mapNoAvalableBank["ps_order"] = "none"
	mapNoAvalableBank["status"] = transaction.TransactionState_Error
	resultMock := Manager.Hold(TransactionTest)
	assert.Equal(t, mapNoAvalableBank, resultMock)
}
