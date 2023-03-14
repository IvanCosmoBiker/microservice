package orange

import (
	"encoding/json"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal/interface/fr"
	core "ephorservices/ephorsale/fiscal/orange/core"
	orangeRequestCheck "ephorservices/ephorsale/fiscal/orange/request"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var Trasaction *transaction.Transaction
var ConfigTest *config.Config

var PrivateKeySign = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDXNEfOOxwuSauj
uZMsjvR0BOsr4ehjK3rbxX3wLdwh0nxAJ91vJ6nP2C54wULJCF8J+FGLKVcP07dV
h7m7o1rWXRuiLd31Flj1/gLQeZo0eJgPeACtLR5+xvL46I73vJmyHPlTKHpKvrph
9Nd7KnqkzDT6yDdGKaPuYqRn0wtW4mlyalbma6i4MEaOBb8mAmpA/PCGXUtwIj9X
CB4T5yFNWBvyexewgsoaNvUGAwrPbpvCuAflmmsdYHYIKfjv/ZO5xjudINm2adXZ
yVXMHHM0A/ACaPrgnxDbPOTG7bRutqge9L/yh1mHeHXgz9XObXlV6d1n0kmBDXEy
gdeItuDVAgMBAAECggEBALtaEXVScps9mcbkxWMSZXEn4xEGEEldzgzMp3JUioOL
eo5j5lxh3G1NGFAaeCkKN6s3Ws5bRCdMOxykF6dqdKeQ0YDki4pWVUZ7SDn007IA
lulIoNYjJJxcWaUm2WiF8gxlOw4RfD3cQ+kJvhrFBZa5DRqS+cQEdmoPyG93BTUy
N7Sp97m8D374e4mAavCjj0G316x4g3okADVi1QsPvbu4tSpx9x9iRXSzYJegLdW7
DIEStoiJGYk/mS1GjzkquH19c5hvtnRVFXXlZsdy0VL3N1UH2iEKIsDGKr7dQsbN
vP7VEVERSRT7vcFITFh4ePL/RkI/2Qrt0u9QKKQ4iNUCgYEA2kFujDBmt5GK2ZW3
w4JzKRmQVIp1+YMhqbLsiqWiKA5AwABOOuykMD5k2lhCTYN7NsFcRZsTzk0bWiYV
zj6mnjde64w3v/L2/LvXvg1oCSYAdu2oGy8lWrEHWE9JJukxZ/iSYwoVJb46SizZ
qBdhC+jTeOsI/67as7asCvqECisCgYEA/GvCSBE8A14dAz/suRR6p9cEwh8JAN8b
09GZpr4qA0SilEPOB3qUKgFBcyR0buSdD6MUKKZNcK7vqbpe88B7tCpt4G/f1xVe
NZVDsgsXbXz6BIqvrL/faZ1DOPGPjkVFMk9N8T+tfKywY+/qU2tjBKn10QXurk07
6p8OkCDyQP8CgYEAiSTGb0bWtJC63CCM+UhWTsQmgkkC+sdgdr7cjf6oV10laMCI
Z9RdE4eRXfZJq2VsHisAbSiWGHMxNcNqvk916UNH3OEeAvqMIqFyXpUUA3OipRiP
Io3MfiFxSReBEvdDOV7jtWIXicDv5b4rAsm2DIK/p2KhI/DeskCd+MQUBkMCgYEA
7velakzGn/mNRfJSzbURmav6GT0AbQ7LbXDVIgKOC6ICuJKojnQBqPKfX753bDSK
bK9a+lDWp4M16V1DX0gu1JYGh5/iLeFQ2zGAcSIG/+R9XadeQRE1FOuJJHOsEGiL
5eEmTOqX95wVMceD842KpHOzADu5htIfkzMZumE2d0kCgYB+1aC77TFg0aDQq9hM
3+lrTJ5q1DNVsq/1IZGjyAXv3Guxm/urrG75ztOWWikg+q5mXINUnBDYPu0DnJCF
chf4nYFOAGJ5VQZrVblcQxoN+X3hM1Y4/6Pzg4cPiFGU+QTQjpaajr2voS7GWt7m
cAuNQXAHU4LCkppfW+xm9g3vBw==
-----END RSA PRIVATE KEY-----`)

var CertTest = `-----BEGIN CERTIFICATE-----
MIIDqzCCApMCBgCgs7XxEDANBgkqhkiG9w0BAQsFADBxMQswCQYDVQQGEwJSVTEP
MA0GA1UECAwGTW9zY293MQ8wDQYDVQQHDAZNb3Njb3cxEzARBgNVBAoMCk9yYW5n
ZWRhdGExDzANBgNVBAsMBk5lYnVsYTEaMBgGA1UEAwwRd3d3Lm9yYW5nZWRhdGEu
cnUwHhcNMTgxMDExMTE1NjExWhcNMjMxMDExMTE1NjExWjCBwDELMAkGA1UEBhMC
UlUxFTATBgNVBAgMDDY5MDIwOTgxMjc1MjEPMA0GA1UEBwwGTW9zY293MUAwPgYD
VQQKDDfQmNCfX9CS0L7QudGC0LrQtdCy0LjRh1/QkNC70LXQutGB0LXQuV/QntC7
0LXQs9C+0LLQuNGHMRMwEQYDVQQLDApFLWNvbW1lcmNlMQ8wDQYDVQQDDAZXV1cu
cnUxITAfBgkqhkiG9w0BCQEWEmV4YW1wbGVAZXhhbXBsZS5ydTCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBALPSmaet3vuR1vbyuuraZjdowYLxPhmRQE9J
Q0rP/aButZH2XIqv8IDwomt2OALHNXk5G0ziIKfAhjnCkT0xWQALIdwSGvUx/VY8
1++1a0hplpPMNwNpthOXTcZpuxeOGbWqGlUKmvlBr7X7qFHFB8M4K+xs6MzpXMEf
tXOSVaAWaSqP69Zei0OC7z2Mn/CTKmdtsbJN/UjQKhjFL+PDNIacbVHtZVVLmN9z
e9iwtd0b/JhMlRceBE4Pw+hcqKAUYB/tmygVwb680Mb0VMvvvtkRm0woB4IVK3hd
Ai9mbEdm/D7+pccywgJCqcrNFp1663rmfrTquhynyk8VRgwi8jcCAwEAATANBgkq
hkiG9w0BAQsFAAOCAQEA20Ai9afvtZba09fGwgRW9Q5Zs2lCx52vcDCO3rFDkSmC
w0V6pAHMcUfOh1W/HH98FWuYkf2EY6AYCihtX21sl98EDoBc9f2r+tUuWpRpUS7I
GqD2XI2Zc7EqrOW3rncGKZmvwweKK9FqaRaNZcRxrRB796FR33p5Te7Qjai8ZLd+
v3XL8IFKRbUwfNoQJUK0kghqP8sqFTAvkim3gJ85ZAjibexIZahfGvS+WnSki2+g
fjhw4rC48woAhBSKox/4LWVjDps+Tf0uia8xuXakmA8/JWJ/oxBvfkj0ldx0sONj
pAC0WAjf4zgw3+frDwWxHIMzTxAoRAVhxUalr9o0CQ==
-----END CERTIFICATE-----`
var KeyTest = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCz0pmnrd77kdb2
8rrq2mY3aMGC8T4ZkUBPSUNKz/2gbrWR9lyKr/CA8KJrdjgCxzV5ORtM4iCnwIY5
wpE9MVkACyHcEhr1Mf1WPNfvtWtIaZaTzDcDabYTl03GabsXjhm1qhpVCpr5Qa+1
+6hRxQfDOCvsbOjM6VzBH7VzklWgFmkqj+vWXotDgu89jJ/wkypnbbGyTf1I0CoY
xS/jwzSGnG1R7WVVS5jfc3vYsLXdG/yYTJUXHgROD8PoXKigFGAf7ZsoFcG+vNDG
9FTL777ZEZtMKAeCFSt4XQIvZmxHZvw+/qXHMsICQqnKzRadeut65n606rocp8pP
FUYMIvI3AgMBAAECggEAD3VxZCrcWoAlHMGtM/dmhijpSdp3XjdQcgB4Wnwa76nU
ziGBvyJ06IDHVbmqAwMhI7S3Fhryd7ljUJ/bYIlXf1t1o7eivaV4g+tjHOZZvLXn
DfmmWRLDZlfBhecdAF9k8msXLGxm+jqdYmWqCK2Jh0zS6dZLBSKiqK+TJ8ZSuhpO
B9wjjL53tZO8nuy4q9LPVkU0Ah4V9R7u/Sfri5hLm5YNpzX5TrZb3R0ejFUW6MLN
eEVBNQwhhxjBl3u1w/EoaYUSz7WEqpzqUZUm/GIdeEW4tR+wylXNQ92HOY7ixRqs
+L8p4ecF7xlUuXLmDTBWEiniN3druGvvZhJ3nlUAIQKBgQDtbC1RfQnFk6KyrEWk
LRa+2F989t7xIX3gZuH1jVplxTLFyiQ8K7VlishLzIbrzdsnRfXsB88//AFvk6Zo
TDAXwFcSxHyVrtpQnPCnkYINiAJK7RMqWLKrrYO39hYNKDef95snchvwco7HXZda
unEQSBYhqniDqJM342bxak6r7wKBgQDB5KL4dQw7XtbGwtXE8KzlSN/BWXbcBmIv
szWfDRQ80A8s+JSziWLUNOh0OuPwgdytMjngTJa3Fm2Je6OQlnZlg5k+eHwqVFWC
JApgL5BqDTYqKixQZG2DcUyQJlDPoLWsRtZVTGoo+rK2CxR+55CRv0/BZfm/Ou86
0Ka2BAn2OQKBgHnRWWdIOq1PVNlMHudf4x3EsynRGQ6r2oQ7BZESF+HDzotBbloZ
KxeQn7iUll2C4AFEmiuizinMSYhQP7+f58UoAQU2H55FeuqFu8yekhYTROngvkap
//KqMr0+3I2fpvrC9q7Ek6VJggy07qW0p7Js6j4X04HqCq9QVE9l9jutAoGAbzqW
MyoSZky1sTg8IcpfpPj1Q5nrEbWnxe1sqV17apeA3S+NPqFlzI69e0/9Sw90ZPcX
NJE7NLTtCZ2f62YlbX7c/nVn5XCTzSCXwy4GDpCdrfqbiVLTcEAix97zJOjwz2+j
rTM1A2Ut+DjK/TIiQToaqruxVf6dFoRz3p7aiCECgYEAtTTP6ypvNiQQPloMTlvz
8C7X1JnlzFY+Z83P5zfJ5PMTrASFm388mTUFmEOdM3TFBrOnTB4MyAs7X7A1Vf8B
bQq0BVh6lkClVlGKins1ojCE12DV+O2SwQ0h2rimSOuH108cqRwXpOK+xWqHphik
5Qh/w6trz2eBXlwfnECBGCA=
-----END PRIVATE KEY-----`

func init() {
	Trasaction = GenerateDataTransactionForTest()
}

func GenerateDataTransactionForTest() *transaction.Transaction {
	cfg := config.Config{}
	cfg.Services.EphorFiscal.PathCert = "../../../../cert"
	ConfigTest = &cfg
	tran := transaction.Transaction{}
	tran.Fiscal.Config.Name = "Orange"
	tran.Fiscal.Config.Type = fr.Fr_EphorServerOrangeData
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
		product := &transaction.Product{}
		product.Price = float64(100)
		product.Quantity = int64(i + 1)
		product.Tax_rate = int32(transaction.TaxRate_NDS18)
		product.Name = fmt.Sprintf("Product%v", i)
		tran.Products = append(tran.Products, product)
	}
	return &tran
}

func TestGenerateDataForCheck(t *testing.T) {
	orangeCore := core.InitCore(ConfigTest)
	requestTest := orangeRequestCheck.RequestSendCheck{}
	request := &orangeRequestCheck.RequestSendCheck{}
	for _, product := range Trasaction.Products {
		entryPayments := orangeRequestCheck.Payment{}
		entryPositions := orangeRequestCheck.Position{}
		quantity := parserTypes.ParseTypeInFloat64(product.Quantity)
		price := parserTypes.ParseTypeInFloat64(product.Price)
		entryPayments.Type = 2
		entryPayments.Amount = math.Round(quantity * price)

		entryPositions.PaymentMethodType = 4
		entryPositions.PaymentSubjectType = 1
		entryPositions.Quantity = int64(quantity)
		entryPositions.Price = math.Round(price)
		entryPositions.Tax = uint8(orangeCore.ConvertTax(parserTypes.ParseTypeInterfaceToInt(product.Tax_rate)))
		entryPositions.Text = parserTypes.ParseTypeInString(product.Name)
		requestTest.Content.CheckClose.Payments = append(requestTest.Content.CheckClose.Payments, entryPayments)
		requestTest.Content.Positions = append(requestTest.Content.Positions, entryPositions)
	}
	err := orangeCore.GenerateDataForCheck(request, Trasaction)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, request.Content.CheckClose.Payments, requestTest.Content.CheckClose.Payments, "they should be equal")
	assert.Equal(t, request.Content.Positions, requestTest.Content.Positions, "they should be equal")
}

func TestReadFileCertificate(t *testing.T) {
	orangeCore := core.InitCore(ConfigTest)
	cert, key, err := orangeCore.ReadFileCertificate(Trasaction)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, cert, CertTest, "they should be equal")
	assert.Equal(t, key, KeyTest, "they should be equal")
}

func TestComputeSignature(t *testing.T) {
	orangeCore := core.InitCore(ConfigTest)
	request := &orangeRequestCheck.RequestSendCheck{}
	err := orangeCore.GenerateDataForCheck(request, Trasaction)
	if err != nil {
		t.Error(err)
	}
	data, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	signature, err := orangeCore.ComputeSignature(data, PrivateKeySign)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, signature, "cPm4KGSiJ5b/dywJqA3qA18UGzYDE+o4rXAqiJ8TVVldMlqcH0VAtN8AGjO1SOrp/jRIW6XpfL8T5TxmXeY6G9bjTJjlF0rqq8+EIqcIzScFtzgfebim4q1z2kHS0PMGXjullVFJTxKaovywxgtFg4n8uba7uf3nEw0pukEPgm1illQ3iL9SimhuVb7UR62TxP7EVe+QOvGs/vHpoGQsG7I0Qr3QZ8UTrOmgNHI6Cf8MAvKnTKHO+MDoGvwc/egY4v4fZdi1OoiODfMVmF+siN/K6Juq22g3MXEGvKSA0iULtTyGOXodWFdHOc1fDeXmViKXtf67rDFrbQRu1rZ4Cg==", "they should be equal")
}
