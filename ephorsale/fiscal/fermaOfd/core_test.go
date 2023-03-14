package fermaOfd

import (
	"ephorservices/ephorsale/fiscal/interface/fr"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TransactionTest *transaction.Transaction
var CoreTest *Core

func init() {
	TransactionTest = transaction.InitTransaction()
	TransactionTest.Fiscal.Config.Phone = "+7-930-152-97-65"
	TransactionTest.Fiscal.Config.Email = "ysamberi@bk.ru"
	TransactionTest.Fiscal.Config.Login = "test"
	TransactionTest.Fiscal.Config.Password = "12345"
	TransactionTest.Fiscal.Config.Name = "OFD_TEST"
	TransactionTest.Fiscal.Config.Type = fr.Fr_OFD
	TransactionTest.Fiscal.Config.AutomatNumber = 45
	TransactionTest.Fiscal.Config.Dev_addr = "ferma-test.ofd.ru"
	TransactionTest.Fiscal.Config.Dev_port = 0
	TransactionTest.Fiscal.Config.Ofd_addr = "ferma-test.ofd.ru"
	TransactionTest.Fiscal.Config.Ofd_port = 0
	TransactionTest.Fiscal.Config.Inn = "1234567890"
	TransactionTest.Fiscal.Config.Auth_public_key = ""
	TransactionTest.Fiscal.Config.Sign_private_key = ""
	TransactionTest.Fiscal.Config.Auth_private_key = ""
	TransactionTest.Fiscal.Config.Param1 = "4010004"
	TransactionTest.Fiscal.Config.Use_sn = 1
	TransactionTest.Fiscal.Config.Add_fiscal = 1
	TransactionTest.Fiscal.Config.Ffd_version = 1
	TransactionTest.Fiscal.Config.MaxSum = 50000
	TransactionTest.Fiscal.Config.CancelCheck = 0
	TransactionTest.TaxSystem.Type = 1
	TransactionTest.Config.AccountId = 1
	TransactionTest.Date = "2023-03-10 10:00:00"
	for i := 0; i < 2; i++ {
		payment_device := "DA"
		if i == 1 {
			payment_device = "CA"
		}
		product := transaction.Product{
			Name:           fmt.Sprintf("Product_%v", i),
			Payment_device: payment_device,
			Price_list:     int32(1),
			Type:           int32(1),
			Ware_id:        int32(0),
			Select_id:      fmt.Sprintf("%v", i),
			Value:          float64(500),
			Price:          float64(500),
			Tax_rate:       int32(0),
			Quantity:       int64(1000),
			Fiscalization:  true,
		}
		TransactionTest.Products = append(TransactionTest.Products, &product)
	}
	CoreTest = InitCore()
}

func TestMakeRequestAuth(t *testing.T) {
	testString := `{"Login":"test","Password":"12345"}`
	json, _ := CoreTest.MakeRequestAuth(TransactionTest)
	assert.Equal(t, testString, string(json))
}

func TestMakeRequestSendCheckTwoPaymets(t *testing.T) {
	testString := `{"Request":{"Inn":"1234567890","Type":"Income","InvoiceId":"23101678442400","CallbackUrl":"","CustomerReceipt":{"KktFA":true,"TaxationSystem":"Common","Email":"ysamberi@bk.ru","Phone":"","PaymentType":1,"AutomaticDeviceNumber":0,"BillAddress":"","Items":[{"Label":"Product_0","Price":5,"Quantity":1,"Amount":5,"Vat":"VatNo","PaymentMethod":4,"PaymentType":1},{"Label":"Product_1","Price":5,"Quantity":1,"Amount":5,"Vat":"VatNo","PaymentMethod":4,"PaymentType":1}],"PaymentItems":[{"PaymentType":1,"Sum":5},{"PaymentType":0,"Sum":5}]}}}`
	json, _ := CoreTest.MakeRequestSendCheck(TransactionTest)
	assert.Equal(t, testString, string(json))
}

func TestMakeRequestSendCheck(t *testing.T) {
	for _, product := range TransactionTest.Products {
		product.Payment_device = "CA"
	}
	testString := `{"Request":{"Inn":"1234567890","Type":"Income","InvoiceId":"23101678442400","CallbackUrl":"","CustomerReceipt":{"KktFA":true,"TaxationSystem":"Common","Email":"ysamberi@bk.ru","Phone":"","PaymentType":1,"AutomaticDeviceNumber":0,"BillAddress":"","Items":[{"Label":"Product_0","Price":5,"Quantity":1,"Amount":5,"Vat":"VatNo","PaymentMethod":4,"PaymentType":1},{"Label":"Product_1","Price":5,"Quantity":1,"Amount":5,"Vat":"VatNo","PaymentMethod":4,"PaymentType":1}],"PaymentItems":[{"PaymentType":0,"Sum":10}]}}}`
	json, _ := CoreTest.MakeRequestSendCheck(TransactionTest)
	assert.Equal(t, testString, string(json))
}
