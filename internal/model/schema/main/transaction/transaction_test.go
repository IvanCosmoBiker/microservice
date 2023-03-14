package transaction

import (
	"context"
	storeTransaction "ephorservices/internal/model/schema/main/transaction/store"
	dateTime "ephorservices/pkg/datetime"
	"ephorservices/pkg/logger"
	connectionPostgresql "ephorservices/pkg/orm/db"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ConnectionDb *connectionPostgresql.Manager

func init() {
	logger.New("log.txt", "Test", []int{1, 2, 3, 4, 5, 6}, false)
	var err error
	//login, password, address, databaseName string, port, pgConnectionPool, pgConnectionMin, pgConnectionMax uint16, healthCheckPeriod int, preferSimpleProtocol bool, debug bool
	ConnectionDb, err = connectionPostgresql.Init("postgres", "123", "127.0.0.1", "local", uint16(5432), uint16(10), uint16(10), uint16(10), 30, false, true, context.Background())
	if err != nil {
		panic(err)
	}
}

func TestAddTransaction(t *testing.T) {
	storeTran := storeTransaction.New(1)
	args := make(map[string]interface{})
	args["automat_id"] = 20
	transaction, err := storeTran.AddByParams(args)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(20), transaction.Automat_id, "they should be equal")
}

func TestSetTransaction(t *testing.T) {
	date, _ := dateTime.Init()
	storeTran := storeTransaction.New(1)
	args := make(map[string]interface{})
	args["automat_id"] = 20
	transaction, err := storeTran.AddByParams(args)
	fmt.Printf("\n%+v_____________\n", transaction)
	if err != nil {
		t.Error(err)
	}
	args2 := make(map[string]interface{})
	args2["id"] = transaction.Id
	args2["automat_id"] = 21
	args2["date"] = date.Now()
	transaction, err = storeTran.SetByParams(args2)
	fmt.Printf("\n%+v_____________\n", transaction)
	if err != nil {
		fmt.Printf("\n%v\n ______________________", err)
		t.Error(err)
	}
	assert.Equal(t, nil, err, "they should be equal")
}
