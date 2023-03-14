package coolerprepaid

import (
	config "ephorservices/config"
	interfacePayment "ephorservices/ephorsale/payment/interfacePayment"
	"ephorservices/ephorsale/sale/coolerprepaid"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	MockPayment "ephorservices/test/mock/payment/manager"
	MockPaymentSystem "ephorservices/test/mock/payment/paymentSystem"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var Config *config.Config
var TransactionTest *transaction.Transaction
var Cooler *coolerprepaid.SaleCoolerPrePaid

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
	TransactionTest.Payment.Type = interfacePayment.TypeSber
	Cooler = &coolerprepaid.SaleCoolerPrePaid{
		Cfg: Config,
	}
}

func makePaymentMock(t *testing.T) {

}

func makeFiscalMock(t *testing.T) {

}

func makeQueueBrokerMock(t *testing.T) {

}

// var (
// 	IceboxStatus_Drink  = 1 // выдача напитков
// 	IceboxStatus_Icebox = 2 // дверь открыта
// 	IceboxStatus_End    = 3 // выдача завершена
// )

// var (
// 	VendState_Session   = 1 //[11]PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
// 	VendState_Approving = 2 //[14] Продукт выбран. Ожидание оплаты.
// 	VendState_Vending   = 3 //[12]PAY_OK_AUTOMAT_PREPARE Оплата успешна, ТА готовит продукт
// 	VendState_VendOk    = 4 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
// 	VendState_VendError = 5 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
// )

// var (
// 	VendError_VendFailed       = 769 //769 Ошибка выдачи продукта
// 	VendError_SessionCancelled = 770 //770
// 	VendError_SessionTimeout   = 771 //771
// 	VendError_WrongProduct     = 772 //772
// 	VendError_VendCancelled    = 773 //773
// 	VendError_ApprovingTimeout = 774 //774
// )

// func (scpp *SaleCoolerPrePaid) Payment(tran *transaction.Transaction) (map[string]interface{}, interfacePayment.Payment) {
// 	resultDb := make(map[string]interface{})
// 	resultPayment, PaymentSystem := scpp.PaymentM.StartPayment(tran)
// 	if resultPayment["status"] != transaction.TransactionState_MoneyHoldStart {
// 		return resultPayment, nil
// 	}
// 	resultDb["id"] = tran.Config.Tid
// 	resultDb["status"] = transaction.TransactionState_MoneyHoldWait
// 	scpp.Dispetcher.StoreTransaction.SetByParams(resultDb)
// 	resultPay := scpp.PaymentM.HoldMoney(tran, PaymentSystem)
// 	if resultPay["status"] == transaction.TransactionState_Error {
// 		return resultPay, nil
// 	}
// 	return resultPay, PaymentSystem
// }

func generatePaymentErrResponse(t *testing.T) (*gomock.Controller, *MockPayment.MockManagerPayment) {
	response := make(map[string]interface{})
	response["status"] = transaction.TransactionState_Error
	controller := gomock.NewController(t)
	MockInterface := MockPayment.NewMockManagerPayment(controller)
	MockInterface.EXPECT().InitPayments(Config).Return()
	MockInterface.InitPayments(Config)
	return controller, MockInterface
}

// func generatePaymentOkResponse(t *testing.T) (*gomock.Controller, *MockPaymentSystem.MockPayment) {

// }

// func generateFiscalErrResponse(t *testing.T) (*gomock.Controller, *MockPaymentSystem.MockPayment) {

// }

// func generateFiscalOkResponse(t *testing.T) (*gomock.Controller, *MockPaymentSystem.MockPayment) {

// }

// func generateBrokerResponse(t *testing.T,err ) (*gomock.Controller, *MockPaymentSystem.MockPayment) {

// }

func TestCOOLERPREPAID_PaymentErr(t *testing.T) {
	controller, paymentManager := generatePaymentErrResponse(t)
	defer controller.Finish()
	controllerPaymentSystem := gomock.NewController(t)
	defer controllerPaymentSystem.Finish()
	mockPaymentSystem := MockPaymentSystem.NewMockPayment(controller)
	paymentSystem := mockPaymentSystem.EXPECT()
	mapErrPaymentStartPayment := make(map[string]interface{})
	mapErrPaymentStartPayment["status"] = transaction.TransactionState_Error
	paymentManager.EXPECT().StartPayment(TransactionTest).Return(mapErrPaymentStartPayment, paymentSystem)
	Cooler.PaymentM = paymentManager
	TransactionTest.Payment.Type = interfacePayment.TypeSber
	result, _ := Cooler.Payment(TransactionTest)
	fmt.Printf("%v\n", result)
	assert.Equal(t, true, result)
}
