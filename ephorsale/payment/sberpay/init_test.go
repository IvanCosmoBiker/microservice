package sberpay

import (
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"fmt"
	"testing"
)

var TransactionTest *transaction.Transaction

func init() {
	TransactionTest = transaction.InitTransaction()
	TransactionTest.Payment.Login = "P690209812752-api"
	TransactionTest.Payment.Password = "xG-zG-z-YZmE8,f"
	TransactionTest.Payment.Sum = 300
	TransactionTest.Payment.Description = "TEST"
	TransactionTest.Payment.TokenType = transaction.TypeTokenSberPayAndroid
}

func TestHold(t *testing.T) {
	sber := NewSberPayStruct{}
	sberpay := sber.New(true)
	result := sberpay.Hold(TransactionTest)
	fmt.Sprintf("%+v", result)
}
