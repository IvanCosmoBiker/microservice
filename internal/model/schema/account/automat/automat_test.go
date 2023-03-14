package automat

import (
	"context"
	config "ephorservices/config"
	connectionPostgresql "ephorservices/pkg/db"
	storeAutomat "ephorservices/pkg/model/schema/account/automat/store"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func InitConfig() *config.Config {
	cfg := config.Config{}
	cfg.Debug = true
	cfg.Db.Login = "postgres"
	cfg.Db.Password = "123"
	cfg.Db.Address = "127.0.0.1"
	cfg.Db.DatabaseName = "local"
	cfg.Db.PreferSimpleProtocol = false
	cfg.Db.Port = uint16(5432)
	cfg.Db.PgConnectionPool = uint16(10)
	cfg.Db.PgConnectionMin = uint16(10)
	cfg.Db.PgConnectionMax = uint16(10)
	return &cfg
}

func InitConnection(ctx context.Context, cfg *config.Config) *connectionPostgresql.Manager {
	conn, err := connectionPostgresql.Init(cfg, ctx)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return conn
}

func TestAddTransaction(t *testing.T) {
	cfg := InitConfig()
	ctx := context.Background()
	CoreConn := InitConnection(ctx, cfg)
	storeTran := storeAutomat.NewStore(CoreConn)
	args := make(map[string]interface{})
	args["account_id"] = 1
	args["id"] = 20
	transaction, err := storeTran.GetOneWithOptions(args)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(20), transaction.Id, "they should be equal")
}

// func TestSetTransaction(t *testing.T) {
// 	date, _ := dateTime.Init()
// 	cfg := InitConfig()
// 	ctx := context.Background()
// 	CoreConn := InitConnection(ctx, cfg)
// 	storeTran := storeTransaction.NewStore(CoreConn)
// 	args := make(map[string]interface{})
// 	args["account_id"] = 1
// 	args["automat_id"] = 20
// 	transaction, err := storeTran.AddByParams(args)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	args2 := make(map[string]interface{})
// 	args2["id"] = transaction.Id
// 	args2["automat_id"] = 21
// 	args2["date"] = date.Now()
// 	fmt.Printf("\n%+v_____________\n", args2)
// 	transactionset, err := storeTran.SetByParams(args2)
// 	if err != nil {
// 		fmt.Printf("\n%v\n ______________________", err)
// 		t.Error(err)
// 	}
// 	assert.Equal(t, transaction.Id, transactionset.Id, "they should be equal")
// }
