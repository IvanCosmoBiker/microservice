package config

import (
	"context"
	config "ephorservices/config"
	connectionPostgresql "ephorservices/pkg/db"
	storeConfig "ephorservices/pkg/model/schema/account/automat/config/store"
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

func TestGetAutomatConfig(t *testing.T) {
	cfg := InitConfig()
	ctx := context.Background()
	CoreConn := InitConnection(ctx, cfg)
	storeAutomatConfig := storeConfig.NewStore(CoreConn)
	args := make(map[string]interface{})
	args["account_id"] = 1
	args["automat_id"] = 17
	args["to_date"] = nil
	automatConfig, err := storeAutomatConfig.GetOneWithOptions(args)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int32(500), automatConfig.Config_id, "they should be equal")
}
