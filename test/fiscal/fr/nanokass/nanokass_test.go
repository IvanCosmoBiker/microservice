package nanokass

import (
	"encoding/base64"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal/factory"
	interfaceFiscal "ephorservices/ephorsale/fiscal/interface/fr"
	core "ephorservices/ephorsale/fiscal/nanokass/core"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var Trasaction *transaction.Transaction
var ConfigTest *config.Config

func init() {
	Trasaction = GenerateDataTransactionForTest()
}

func GenerateDataTransactionForTest() *transaction.Transaction {
	cfg := config.Config{}
	cfg.Services.EphorFiscal.PathCert = "../../../../cert"
	f := float64(0.05)
	cfg.Services.EphorFiscal.ExecuteMinutes = time.Duration(int(f))
	ConfigTest = &cfg
	tran := transaction.Transaction{}
	tran.Fiscal.Config.Name = "Nanokass"
	tran.Fiscal.Config.Type = interfaceFiscal.Fr_NanoKassa
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
	tran.Fiscal.Config.Param1 = "nanokass"
	tran.Fiscal.Config.Use_sn = 1
	tran.Fiscal.Config.Add_fiscal = 1
	tran.Fiscal.Config.Id_shift = ""
	tran.Fiscal.Config.Fr_disable_cash = 0
	tran.Fiscal.Config.Fr_disable_cashless = 0
	tran.Fiscal.Config.Ffd_version = 1
	tran.Fiscal.Config.Auth_public_key = ""
	tran.Fiscal.Config.Auth_private_key = ""
	tran.Fiscal.Config.Sign_private_key = "111111:222222222222"
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

var testData = []byte("absd")
var HMACTEST = "7da5f94b19e9b17947488e67cc1dd37aa6a269941f06a0f5e47a1629f242fca787f785121ae18b54659f36a90b084c8e0b22018ebb2d17796b943623e9118561"
var DE = "faX5SxnpsXlHSI5nzB3TeqaiaZQfBqD15HoWKfJC/KeH94USGuGLVGWfNqkLCEyOCyIBjrstF3lrlDYj6RGFYa/5NRF4OKCtHsINx6x70Xypf2pg"
var CryptoTestAes256 = []byte{169, 127, 106, 96}
var IVdata = []byte{175, 249, 53, 17, 120, 56, 160, 173, 30, 194, 13, 199, 172, 123, 209, 124}
var pw = []byte{235, 106, 108, 33, 122, 94, 87, 200, 171, 59, 131, 7, 111, 154, 253, 40, 40, 146, 44, 131, 153, 224, 28, 184, 19, 204, 65, 187, 11, 179, 250, 128}

func TestGenerateRandomBytesInString(t *testing.T) {
	Core := core.InitCore(ConfigTest)
	Random16 := Core.GenerateRandomBytesInString(16)
	Random32 := Core.GenerateRandomBytesInString(32)
	if len(Random32) != 32 {
		t.Error("Empty Random")
	}
	if len(Random16) != 16 {
		t.Error("Empty Random")
	}
}

func TestEncryptAES(t *testing.T) {
	Core := core.InitCore(ConfigTest)
	byteData, err := Core.EncryptAES(testData, pw, IVdata)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, byteData, CryptoTestAes256, "they should be equal")
}

func TestHMAC512(t *testing.T) {
	Core := core.InitCore(ConfigTest)
	mk := core.HMAC_FIRST
	dataAES, err := Core.EncryptAES(testData, pw, IVdata)
	if err != nil {
		fmt.Println(err)
	}
	decodeMk, _ := base64.StdEncoding.DecodeString(mk)
	key512 := make([]byte, 0)
	key512 = append(key512, IVdata...)
	key512 = append(key512, dataAES...)
	hmac := Core.HMAC512(key512, decodeMk)
	hmac16 := fmt.Sprintf("%x", hmac)
	assert.Equal(t, hmac16, HMACTEST, "they should be equal")

}

func TestDE(t *testing.T) {
	Core := core.InitCore(ConfigTest)
	mk := core.HMAC_FIRST
	dataAES, err := Core.EncryptAES(testData, pw, IVdata)
	if err != nil {
		fmt.Println(err)
	}
	decodeMk, _ := base64.StdEncoding.DecodeString(mk)
	key512 := make([]byte, 0)
	key512 = append(key512, IVdata...)
	key512 = append(key512, dataAES...)
	hmac := Core.HMAC512(key512, decodeMk)
	DataDe := make([]byte, 0)
	DataDe = append(DataDe, hmac...)
	DataDe = append(DataDe, IVdata...)
	DataDe = append(DataDe, dataAES...)
	returnDataDE := base64.StdEncoding.EncodeToString(DataDe)
	assert.Equal(t, returnDataDE, DE, "they should be equal")
}

/*
Testing when cloud kass not answer and time is out
*/
func TestGetStatusErrTimeOut(t *testing.T) {
	kass := factory.GetFiscal(interfaceFiscal.Fr_ServerNanoKassa, ConfigTest)
	result := kass.GetStatus(Trasaction)
	resultTest := make(map[string]interface{})
	resultTest["status"] = interfaceFiscal.Status_Error
	resultTest["f_desc"] = "Cancelled by a Timeout of Nanokassa"
	resultTest["fr_status"] = interfaceFiscal.Status_Error
	assert.Equal(t, result, resultTest, "they should be equal")
}
