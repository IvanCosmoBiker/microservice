package core

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"ephorservices/ephorsale/fiscal/interface/fr"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var publicKey = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1zRHzjscLkmro7mTLI70\ndATrK+HoYyt628V98C3cIdJ8QCfdbyepz9gueMFCyQhfCfhRiylXD9O3VYe5u6Na\n1l0boi3d9RZY9f4C0HmaNHiYD3gArS0efsby+OiO97yZshz5Uyh6Sr66YfTXeyp6\npMw0+sg3Rimj7mKkZ9MLVuJpcmpW5muouDBGjgW/JgJqQPzwhl1LcCI/VwgeE+ch\nTVgb8nsXsILKGjb1BgMKz26bwrgH5ZprHWB2CCn47/2TucY7nSDZtmnV2clVzBxz\nNAPwAmj64J8Q2zzkxu20braoHvS/8odZh3h14M/Vzm15VendZ9JJgQ1xMoHXiLbg\n1QIDAQAB\n-----END PUBLIC KEY-----"
var CoreOrange *Core
var Tran *transaction.Transaction

func init() {
	CoreOrange = InitCore("D:/bibikov/commit/cert", true)
	Tran = transaction.InitTransaction()
	Tran.Fiscal.Config.Name = "Orange_Test"
	Tran.Fiscal.Config.Type = fr.Fr_EphorServerOrangeData
	Tran.Fiscal.Config.AutomatNumber = 45
	Tran.Fiscal.Config.Dev_addr = "apip.orangadata.ru"
	Tran.Fiscal.Config.Dev_port = 12003
	Tran.Fiscal.Config.Ofd_addr = "apip.orangadata.ru"
	Tran.Fiscal.Config.Ofd_port = 12003
	Tran.Fiscal.Config.Inn = "test_inn"
	Tran.Fiscal.Config.Auth_public_key = ""
	Tran.Fiscal.Config.Sign_private_key = ""
	Tran.Fiscal.Config.Auth_private_key = ""
	Tran.Fiscal.Config.Param1 = "4010004"
	Tran.Fiscal.Config.Use_sn = 1
	Tran.Fiscal.Config.Add_fiscal = 1
	Tran.Fiscal.Config.Ffd_version = 1
	Tran.Fiscal.Config.MaxSum = 50000
	Tran.Fiscal.Config.CancelCheck = 0
	Tran.TaxSystem.Type = 1
	Tran.Config.AccountId = 1
	Tran.Date = "2023-03-10 10:00:00"
	for i := 0; i < 2; i++ {
		product := transaction.Product{
			Name:           fmt.Sprintf("Product_%v", i),
			Payment_device: "DA",
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
		Tran.Products = append(Tran.Products, &product)
	}
}

func DecodeSignature(message []byte, rawSignature, rawPubKey string) error {
	block, _ := pem.Decode([]byte(rawPubKey))
	if block == nil {
		return errors.New("pem block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pubKey := key.(*rsa.PublicKey)
	signature, err := base64.StdEncoding.DecodeString(rawSignature)
	if err != nil {
		return err
	}
	pack := CoreOrange.Pack([]byte("3031300d060960864801650304020105000420"))
	pack = append(pack, message...)
	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, pack)
	if err != nil {
		buf = nil
		return err
	}
	h := sha256.New()
	h.Write(buf.Bytes())
	hash := h.Sum(nil)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return err
	}
	fmt.Println("Successfully verified message with signature and public key")
	return nil
}

func TestMakingRequestCheck(t *testing.T) {
	signature, resiptId, jsonRequest, _ := CoreOrange.MakeRequestSendCheck(Tran)
	signatureTest := "zqxib0wKjX2gIpDY71WMMoLoOLL5saF8eJTJbnCyBwAfGZZVqMzaUpP+6IFU84w5re4t/ngH8T2vfqAqg70+7fo09czXCCVR1DRQZU5xAvyb/MTIdyzo7qYlWRlJEGdattKzaSv5ls3N+Uf7mpO+kUR4E21HfC/Jg/xo1ddnNNSRbsKir/Z80DM0686pebw3I5/DGRGMbrLb8FlQGYk/Huz0lOgTkskD9TL56MIK5B9l9L+mof5GyDyaC/zO7sS9Z9eJ6FyKMzCMnA5tAlvj+mn/2Yf8wRwJpmGRO1/CNdwHMKuykn8+jlUCLI8HyCwIx7IsUSXKx9mtDo9wuOW0rA=="
	resipeTest := "23101678442400"
	jsonTest := `{"Id":"23101678442400","Inn":"test_inn","Group":"4010004","Key":"4010004","Content":{"FfdVersion":2,"Type":1,"AutomatNumber":"45","Positions":[{"Quantity":1,"Price":5,"Tax":6,"Text":"Product_0","PaymentMethodType":4,"PaymentSubjectType":1},{"Quantity":1,"Price":5,"Tax":6,"Text":"Product_1","PaymentMethodType":4,"PaymentSubjectType":1}],"CheckClose":{"Payments":[{"Type":2,"Amount":5},{"Type":2,"Amount":5}],"TaxationSystem":0},"SettlementAddress":"","SettlementPlace":""}}`
	assert.Equal(t, resipeTest, resiptId, "they should be equal")
	assert.Equal(t, signatureTest, signature, "they should be equal")
	assert.Equal(t, jsonTest, string(jsonRequest), "they should be equal")
}

func TestMakingRequestCheckMarking(t *testing.T) {
	Tran.Fiscal.Config.Ffd_version = 2
	for _, product := range Tran.Products {
		product.Marking = 1
		product.Barcode = "95634726937"
	}
	signature, resiptId, jsonRequest, _ := CoreOrange.MakeRequestSendCheck(Tran)
	signatureTest := "JB8Qd76pv0xTlmr2dR5+Vt5VRPL0ZBkrmPLGoUjeYy+QVzxh9lsKclVdzaaPG8S2ldXzR24FyKmJmFyJsIzPimGrQhgG4sDMFeke9BuqvQjx1+QaBeA4AIMofm4NAovSEmuKOydxNlhFySOqUz+XC6WKUOF4rx6QsVruV3rLRATw4drjYQIFrGCPo6Im4IYHoqyQEs9sCIWb3iJ/ay08J2+0MYTVYiSAJLbbqdl9OvTGGLBwpXLDskzqvVJFvU7sbbQ1tM05IIOj0sCY5ShDJ8zH9R2a3DMvyMjuM3lnB36WTT+HlrAcO88GhBkRk9cQcqG54sTynL19NRD/O6W9CA=="
	resipeTest := "23101678442400"
	jsonTest := `{"Id":"23101678442400","Inn":"test_inn","Group":"4010004","Key":"4010004","ignoreItemCodeCheck":true,"Content":{"FfdVersion":4,"Type":1,"AutomatNumber":"45","Positions":[{"itemCode":"095634726937","Quantity":1,"Price":5,"Tax":6,"Text":"Product_0","PaymentMethodType":4,"PaymentSubjectType":33,"quantityMeasurementUnit":"0","plannedStatus":2,"IndustryAttribute":{"FoivId":"030","CauseDocumentDate":"16.06.2022","CauseDocumentNumber":"174","Value":"crpt=mrk\u0026mode=vend"}},{"itemCode":"095634726937","Quantity":1,"Price":5,"Tax":6,"Text":"Product_1","PaymentMethodType":4,"PaymentSubjectType":33,"quantityMeasurementUnit":"0","plannedStatus":2,"IndustryAttribute":{"FoivId":"030","CauseDocumentDate":"16.06.2022","CauseDocumentNumber":"174","Value":"crpt=mrk\u0026mode=vend"}}],"CheckClose":{"Payments":[{"Type":2,"Amount":5},{"Type":2,"Amount":5}],"TaxationSystem":0},"SettlementAddress":"","SettlementPlace":"","fsItemCodeType":"0"}}`
	assert.Equal(t, resipeTest, resiptId, "they should be equal")
	assert.Equal(t, signatureTest, signature, "they should be equal")
	assert.Equal(t, jsonTest, string(jsonRequest), "they should be equal")

}

func TestGenerateSignature(t *testing.T) {
	Tran.Fiscal.Config.Ffd_version = 2
	for _, product := range Tran.Products {
		product.Marking = 1
		product.Barcode = "95634726937"
	}
	signature, _, data, _ := CoreOrange.MakeRequestSendCheck(Tran)
	err := DecodeSignature(data, signature, publicKey)
	assert.Equal(t, nil, err)
}
