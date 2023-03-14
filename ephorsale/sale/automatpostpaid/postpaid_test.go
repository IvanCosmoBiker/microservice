package automatpostpaid

import (
	"context"
	config "ephorservices/config"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	transport_manager "ephorservices/ephorsale/transport"
	connectionDb "ephorservices/pkg/orm/db"
	"testing"
)

var ContextApp context.Context
var Transaction *transaction.Transaction

func init() {
	Config := &config.Config{}
	ContextApp = context.Background()
	//"goadmin", "go2021", "188.225.18.140", "cardtest", uint16(6432), uint16(10), uint16(2), true, true
	connectionDb.Init("goadmin", "go2021", "188.225.18.140", "cardtest", uint16(6432), uint16(10), uint16(2), uint16(2), 10, false, true, ContextApp)
	transaction_dispetcher.New(Config, ContextApp)
	Transaction = transaction_dispetcher.Dispetcher.NewTransaction()
	//Transaction
}
