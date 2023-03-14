package core

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal/interface/fr"
	orangeRequestCheck "ephorservices/ephorsale/fiscal/orange/request"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	randString "ephorservices/pkg/randgeneratestring"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

var privateKeySign = []byte(`-----BEGIN RSA PRIVATE KEY-----
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

var (
	layoutQr = "20060102T150405"
)

type Core struct {
	cfg *config.Config
}

func InitCore(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}

func (c *Core) MakeRequestSendCheck(tran *transaction.Transaction) (signature string, resipeId string, jsonRequest []byte, err error) {
	request := &orangeRequestCheck.RequestSendCheck{}
	orderString := randString.Init()
	orderString.RandGiud()
	request.Id = orderString.String
	orderString = nil
	request.Inn = tran.Fiscal.Config.Inn
	request.Group = tran.Fiscal.Config.Param1
	request.Key = tran.Fiscal.Config.Inn
	request.Content.FfdVersion = c.FfdVersion(tran.Fiscal.Config.Ffd_version)
	request.Content.SettlementAddress = tran.Point.Address
	request.Content.SettlementPlace = tran.Point.PointName
	request.Content.AutomatNumber = parserTypes.ParseTypeInString(tran.Config.AutomatId)
	if tran.Fiscal.Config.Use_sn == 1 {
		request.Content.AutomatNumber = parserTypes.ParseTypeInString(tran.Fiscal.Config.AutomatNumber)
	}
	if tran.Fiscal.Config.Type == fr.Fr_EphorServerOrangeData || tran.Fiscal.Config.Type == fr.Fr_EphorOrangeData {
		request.Key = "4010004"
	}
	if tran.Fiscal.Config.CancelCheck == 0 {
		request.Content.Type = 1
	} else {
		request.Content.Type = 2
	}
	request.IgnoreItemCodeCheck = true
	request.Content.CheckClose.TaxationSystem = c.ConvertTaxationSystem(tran.TaxSystem.Type)
	c.GenerateDataForCheck(request, tran)
	jsonRequest, err = json.Marshal(request)
	if err != nil {
		return signature, request.Id, []byte(""), err
	}
	if tran.Fiscal.Config.Type != fr.Fr_EphorServerOrangeData && tran.Fiscal.Config.Type != fr.Fr_EphorOrangeData {
		signature, err = c.ComputeSignature(jsonRequest, []byte(tran.Fiscal.Config.Sign_private_key))
		return signature, request.Id, jsonRequest, err
	}
	cert, key, err := c.ReadFileCertificate(tran)
	if err != nil {
		return signature, request.Id, jsonRequest, err
	}
	tran.Fiscal.Config.Auth_public_key = cert
	tran.Fiscal.Config.Auth_private_key = key
	log.Printf("%v", tran.Fiscal.Config)
	signature, err = c.ComputeSignature(jsonRequest, privateKeySign)
	return signature, request.Id, jsonRequest, err
}

func (c *Core) MakeRequestStatusQr(tran *transaction.Transaction) string {
	return fmt.Sprintf("https://%s:%v/api/v2/documents/%s/status/%s", tran.Fiscal.Config.Dev_addr, tran.Fiscal.Config.Dev_port, tran.Fiscal.Config.Inn, tran.Fiscal.ResiptId)
}

func (c *Core) EncodeUrlToBase64(tran *transaction.Transaction) string {
	stringUrl := c.MakeUrlQr(tran)
	return base64.StdEncoding.EncodeToString([]byte(stringUrl))
}

func (c *Core) MakeUrlQr(tran *transaction.Transaction) string {
	stringResult := fmt.Sprintf("t=%s&s=%v&fn=%v&i=%v&fp=%v&n=1", tran.Fiscal.Fields.DateFisal, fmt.Sprintf("%v.00", tran.Fiscal.SumCheck), tran.Fiscal.Fields.Fn, tran.Fiscal.Fields.Fd, tran.Fiscal.Fields.Fp)
	log.Println(stringResult)
	return stringResult
}

func (c *Core) FfdVersion(version int) uint8 {
	switch version {
	case 1:
		return 2
	case 2:
		return 4
	}
	return 2
}

func (c *Core) RsaEncrypt(origData, privatesignKey []byte) (encrypt []byte, err error) {
	rng := rand.Reader
	block, _ := pem.Decode(privateKeySign)
	if block == nil {
		return nil, errors.New("is not pem")
	}
	pubInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PrivateKey)
	encrypt, err = rsa.SignPKCS1v15(rng, pub, crypto.SHA256, origData[:])
	if err != nil {
		return nil, err
	}
	return encrypt, nil
}

func (c *Core) ComputeSignature(data interface{}, privatesignKey []byte) (signature string, err error) {
	buf := bytes.Buffer{}
	err = binary.Write(&buf, binary.BigEndian, data)
	if err != nil {
		return signature, err
	}
	h := sha256.New()
	h.Write(buf.Bytes())
	hash := h.Sum(nil)
	result, err := c.RsaEncrypt([]byte(hash), privatesignKey)
	if err != nil {
		return signature, err
	}
	signature = base64.StdEncoding.EncodeToString(result)
	return signature, nil
}

func (c *Core) GenerateDataForCheck(request *orangeRequestCheck.RequestSendCheck, tran *transaction.Transaction) error {
	for _, product := range tran.Products {
		entryPayments := orangeRequestCheck.Payment{}
		entryPositions := orangeRequestCheck.Position{}
		var value float64 = 0
		if product.Value == value {
			value = product.Price
		} else {
			value = product.Value
		}
		quantity := parserTypes.ParseTypeInFloat64(product.Quantity) / float64(1000)
		price := value / float64(100)
		entryPayments.Type = 2
		entryPayments.Amount = math.Round(quantity * price)
		tran.Fiscal.SumCheck += entryPayments.Amount
		entryPositions.PaymentMethodType = 4
		entryPositions.PaymentSubjectType = 1
		entryPositions.Quantity = int64(quantity)
		entryPositions.Price = math.Round(price)
		entryPositions.Tax = uint8(c.ConvertTax(parserTypes.ParseTypeInterfaceToInt(product.Tax_rate)))
		entryPositions.Text = parserTypes.ParseTypeInString(product.Name)
		request.Content.CheckClose.Payments = append(request.Content.CheckClose.Payments, entryPayments)
		request.Content.Positions = append(request.Content.Positions, entryPositions)
	}
	return nil
}

func (c *Core) ReadFileCert() (string, error) {
	crtFilte := fmt.Sprintf("%s%s", c.cfg.Services.EphorFiscal.PathCert, "/ephorOrangeData.crt")
	fcrt, errcrt := os.Open(crtFilte)
	if errcrt != nil {
		return "", errcrt
	}
	var chunkCrt []byte
	bufcrt := make([]byte, 2048)
	for {
		n, err := fcrt.Read(bufcrt)
		if err != nil && err != io.EOF {
			fmt.Println("read buf fail", err)
			return "", err
		}
		if n == 0 {
			break
		}
		chunkCrt = append(chunkCrt, bufcrt[:n]...)
	}
	fcrt.Close()
	return string(chunkCrt), nil
}

func (c *Core) ReadFileKey() (string, error) {
	keyFile := fmt.Sprintf("%s%s", c.cfg.Services.EphorFiscal.PathCert, "/ephorOrangeData.key")
	fkey, errkey := os.Open(keyFile)
	if errkey != nil {
		return "", errkey
	}
	var chunkkey []byte
	bufkey := make([]byte, 2048)
	for {
		n, err := fkey.Read(bufkey)
		if err != nil && err != io.EOF {
			fmt.Printf("read buf fail: %s", err.Error())
			return "", err
		}
		if n == 0 {
			break
		}
		chunkkey = append(chunkkey, bufkey[:n]...)
	}
	fkey.Close()
	return string(chunkkey), nil
}

func (c *Core) ReadFileCertificate(tran *transaction.Transaction) (cert, key string, err error) {
	defer func(cert, key string, err error) (string, string, error) {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			cert = ""
			key = ""
			return cert, key, err
		}
		return cert, key, err
	}(cert, key, err)
	cert, err = c.ReadFileCert()
	if err != nil {
		return "", "", err
	}
	key, err = c.ReadFileKey()
	if err != nil {
		return "", "", err
	}
	cert = strings.ReplaceAll(cert, "\r", "")
	key = strings.ReplaceAll(key, "\r", "")
	return cert, key, nil
}

func (c *Core) ConvertTax(tax int) int {
	switch tax {
	case transaction.TaxRate_NDSNone:
		return 6
	case transaction.TaxRate_NDS0:
		return 5
	case transaction.TaxRate_NDS10:
		return 2
	case transaction.TaxRate_NDS18:
		return 1
	}
	return 6
}

func (c *Core) ConvertTaxationSystem(taxsystem int) uint8 {
	switch taxsystem {
	case transaction.TaxSystem_OSN:
		return 0
	case transaction.TaxSystem_USND:
		return 1
	case transaction.TaxSystem_USNDMR:
		return 2
	case transaction.TaxSystem_ENVD:
		return 3
	case transaction.TaxSystem_ESN:
		return 4
	case transaction.TaxSystem_Patent:
		return 5
	}
	return 0
}
