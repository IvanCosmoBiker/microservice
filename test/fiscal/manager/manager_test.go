package fiscal

import (
	config "ephorservices/config"
	fiscal "ephorservices/ephorsale/fiscal"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var Config *config.Config
var TransactionTest *transaction.Transaction
var FiscalManager *fiscal.FiscalManager

func init() {
	newConfig := &config.Config{}
	newConfig.Db.Password = "go2021"
	newConfig.Db.Address = "188.225.18.140"
	newConfig.Db.DatabaseName = "cardtest"
	newConfig.Db.Login = "goadmin"
	newConfig.Db.Port = 6432
	newConfig.Db.PreferSimpleProtocol = false
	newConfig.Db.PgConnectionMax = 5
	newConfig.Db.PgConnectionPool = 10
	newConfig.Db.PgConnectionMin = 5
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
	FiscalManager = &fiscal.FiscalManager{
		Cfg: Config,
	}
}

func TestCalcSumEmptyProducts(t *testing.T) {
	result := FiscalManager.CalcSumProducts(TransactionTest)
	assert.Equal(t, 0, result)
}

func TestCalcSumProducts(t *testing.T) {
	for i := 0; i < 2; i++ {
		product := &transaction.Product{
			Name:           fmt.Sprintf("test%v", i),
			Payment_device: "CA",
			Type:           3,
			Select_id:      "1",
			Value:          3200,
			Quantity:       1000,
		}
		TransactionTest.Products = append(TransactionTest.Products, product)
	}
	result := FiscalManager.CalcSumProducts(TransactionTest)
	assert.Equal(t, 6400000, result)
}

func TestGetFr(t *testing.T) {

}
