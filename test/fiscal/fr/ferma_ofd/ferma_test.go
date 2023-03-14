package ferma_ofd

import (
	"encoding/json"
	config "ephorservices/config"
	//fermaRequestCheck "ephorservices/ephorsale/fiscal/ferma/request"
	//"ephorservices/ephorsale/fiscal/interfaceFiscal"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var Trasaction *transaction.Transaction
var ConfigTest *config.Config

func init() {
	Trasaction = GenerateDataTransactionForTest()
}

func GenerateDataTransactionForTest() *transaction.Transaction {
	cfg := config.Config{}
	cfg.Services.EphorFiscal.PathCert = "../../../../cert"
	ConfigTest = &cfg
	tran := transaction.Transaction{}
	tran.Fiscal.Config.Name = "Orange"
	tran.Fiscal.Config.Type = interfaceFiscal.Fr_EphorServerOrangeData
	tran.Fiscal.Config.Dev_interface = 2
	tran.Fiscal.Config.Login = "test"
	tran.Fiscal.Config.Password = "test"
	tran.Fiscal.Config.Phone = "test"
	tran.Fiscal.Config.Email = "test"
	tran.Fiscal.Config.Dev_addr = "127.0.0.1"
	tran.Fiscal.Config.Dev_port = 6060
	tran.Fiscal.Config.Ofd_addr = "orange"
	tran.Fiscal.Config.Ofd_port = 1887
	tran.Fiscal.Config.Inn = "1345655467"
	tran.Fiscal.Config.Param1 = "orange"
	tran.Fiscal.Config.Use_sn = 1
	tran.Fiscal.Config.Add_fiscal = 1
	tran.Fiscal.Config.Id_shift = ""
	tran.Fiscal.Config.Fr_disable_cash = 0
	tran.Fiscal.Config.Fr_disable_cashless = 0
	tran.Fiscal.Config.Ffd_version = 1
	tran.Fiscal.Config.Auth_public_key = ""
	tran.Fiscal.Config.Auth_private_key = ""
	tran.Fiscal.Config.Sign_private_key = ""
	tran.TaxSystem.Type = transaction.TaxSystem_ENVD
	for i := 0; i < 3; i++ {
		product := transaction.Product{}
		product.Price = float64(100)
		product.Quantity = int64(i + 1)
		product.Tax_rate = int32(transaction.TaxRate_NDS18)
		product.Name = fmt.Sprintf("Product%v", i)
		tran.Products = append(tran.Products, product)
	}
	return &tran
}

// func TestGenerateDataForCheck(t *testing.T) {
// 	orangeCore := InitCore(ConfigTest)
// 	requestTest := fermaRequestCheck.RequestSendCheck{}
// 	request := &fermaRequestCheck.RequestSendCheck{}
// 	for _, product := range Trasaction.Products {
// 		entryPayments := fermaRequestCheck.Payment{}
// 		entryPositions := fermaRequestCheck.Position{}
// 		quantity := parserTypes.ParseTypeInFloat64(product.Quantity)
// 		price := parserTypes.ParseTypeInFloat64(product.Price)
// 		entryPayments.Type = 2
// 		entryPayments.Amount = math.Round(quantity * price)

// 		entryPositions.PaymentMethodType = 4
// 		entryPositions.PaymentSubjectType = 1
// 		entryPositions.Quantity = int64(quantity)
// 		entryPositions.Price = math.Round(price)
// 		entryPositions.Tax = uint8(orangeCore.ConvertTax(parserTypes.ParseTypeInterfaceToInt(product.Tax_rate)))
// 		entryPositions.Text = parserTypes.ParseTypeInString(product.Name)
// 		requestTest.Content.CheckClose.Payments = append(requestTest.Content.CheckClose.Payments, entryPayments)
// 		requestTest.Content.Positions = append(requestTest.Content.Positions, entryPositions)
// 	}
// 	err := orangeCore.GenerateDataForCheck(request, Trasaction)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	assert.Equal(t, request.Content.CheckClose.Payments, requestTest.Content.CheckClose.Payments, "they should be equal")
// 	assert.Equal(t, request.Content.Positions, requestTest.Content.Positions, "they should be equal")
// }
