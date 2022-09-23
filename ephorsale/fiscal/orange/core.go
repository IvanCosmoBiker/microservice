package orange

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	config "ephorservices/config"
	transaction "ephorservices/ephorsale/transaction"
	transactionStruct "ephorservices/ephorsale/transaction"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"
)

var pathFileCert = "/var/www/html/test/cert"

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

type Core struct {
	cfg *config.Config
}

func Init(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}

func (c *Core) MakeUrlQr(date string, summ int, frResponse map[string]interface{}) string {
	t, _ := time.Parse(layoutISO, date)
	valueSumm := summ / 100
	stringResult := fmt.Sprintf("t=%s&s=%v&fn=%v&i=%v&fp=%v&n=1", fmt.Sprintf("%s", t.Format(layoutQr)), fmt.Sprintf("%v.00", valueSumm), frResponse["fn"], frResponse["fd"], frResponse["fp"])
	log.Println(stringResult)
	return stringResult
}

func (c *Core) RsaEncrypt(origData, privatesignKey []byte) ([]byte, error) {
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
	signature, err := rsa.SignPKCS1v15(rng, pub, crypto.SHA256, origData[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (c *Core) ComputeSignature(data string, privatesignKey []byte) (string, error) {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, data)
	buf.Write([]byte(data))
	h := sha256.New()
	h.Write(buf.Bytes())
	hash := h.Sum(nil)
	result, err := ofd.RsaEncrypt([]byte(hash), privatesignKey)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(result)
	return str, nil
}

func (c *Core) GenerateDataForCheck(transaction transactionStruct.Transaction) ([]map[string]interface{}, []map[string]interface{}) {
	var payments []map[string]interface{}
	var positions []map[string]interface{}
	entryPayments := make(map[string]interface{})
	entryPositions := make(map[string]interface{})
	for _, product := range transaction.Products {
		quantity := float64(product["quantity"].(float64))
		price := float64(product["value"].(float64))
		entryPayments["type"] = 2
		entryPayments["amount"] = math.Round(quantity * price)
		entryPayments["paymentMethodType"] = 4
		entryPayments["paymentSubjectType"] = 1

		entryPositions["quantity"] = product["quantity"]
		entryPositions["price"] = math.Round(price)
		entryPositions["tax"] = ofd.ConvertTax(product["tax_rate"].(int))
		entryPositions["text"] = product["name"]
		payments = append(payments, entryPayments)
		positions = append(positions, entryPositions)
	}
	return payments, positions
}

func (c *Core) ReadFileCertificate() (string, string, error) {
	crtFilte := fmt.Sprintf("%s%s", pathFileCert, "/ephorOrangeData.crt")
	keyFile := fmt.Sprintf("%s%s", pathFileCert, "/ephorOrangeData.key")
	fcrt, errcrt := os.Open(crtFilte)
	if errcrt != nil {
		return "", "", errcrt
	}
	var chunkCrt []byte
	bufcrt := make([]byte, 2048)
	for {
		n, err := fcrt.Read(bufcrt)
		if err != nil && err != io.EOF {
			fmt.Println("read buf fail", err)
			return "", "", err
		}
		if n == 0 {
			break
		}
		chunkCrt = append(chunkCrt, bufcrt[:n]...)
	}
	fcrt.Close()
	fkey, errkey := os.Open(keyFile)
	if errkey != nil {
		return "", "", errkey
	}
	var chunkkey []byte
	bufkey := make([]byte, 2048)
	for {
		n, err := fkey.Read(bufkey)
		if err != nil && err != io.EOF {
			fmt.Println("read buf fail", err)
			return "", "", err
		}
		if n == 0 {
			break
		}
		chunkkey = append(chunkkey, bufkey[:n]...)
	}
	fkey.Close()
	return string(chunkCrt), string(chunkkey), nil
}

func (c *Core) ConvertTax(tax int) int {
	switch tax {
	case transaction.TaxRate_NDSNone:
		return 6
		fallthrough
	case transaction.TaxRate_NDS0:
		return 5
		fallthrough
	case transaction.TaxRate_NDS10:
		return 2
		fallthrough
	case transaction.TaxRate_NDS18:
		return 1
	}
	return 6
}

func (c *Core) ConvertTaxationSystem(taxsystem int) int {
	switch taxsystem {
	case transaction.TaxSystem_OSN:
		return 0
		fallthrough
	case transaction.TaxSystem_USND:
		return 1
		fallthrough
	case transaction.TaxSystem_USNDMR:
		return 2
		fallthrough
	case transaction.TaxSystem_ENVD:
		return 3
		fallthrough
	case transaction.TaxSystem_ESN:
		return 4
		fallthrough
	case transaction.TaxSystem_Patent:
		return 5
	}
	return 0
}
