package core

import (
	"context"
	connection "ephorservices/pkg/db/coninit"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Connect *connection.Connection

func init() {
	ctx := context.Background()
	Conn := connection.NewConnectionPool(ctx)
	err := Conn.Init("postgres", "123", "127.0.0.1", "postgres", uint16(5432), uint16(10), uint16(10), true, true)
	if err != nil {
		panic(err)
	}
	Connect = Conn
}

func TestMakeCore(t *testing.T) {
	var ctx context.Context
	backOffPolicyAcquireConn := setBackOffPolicyAcquireConn()
	coreTest := &CoreConnection{
		Status:                   Status_Work,
		Connection:               Connect,
		ctx:                      ctx,
		BackOffPolicyAcquireConn: backOffPolicyAcquireConn,
	}
	core := NewCoreConnection(ctx, Connect)
	assert.Equal(t, coreTest, core)
}

type EmptyStruct struct {
}

func TestPrepareGetEmpty(t *testing.T) {
	var ctx context.Context
	core := NewCoreConnection(ctx, Connect)
	option := make(map[string]interface{})
	model := EmptyStruct{}
	sqlFields, _, _ := core.PrepareGet(option, model)
	assert.Equal(t, "", sqlFields)
}

func TestPrepareUpdateEmpty(t *testing.T) {
	var ctx context.Context
	core := NewCoreConnection(ctx, Connect)
	option := make(map[string]interface{})
	model := EmptyStruct{}
	sqlValues, _, _ := core.PrepareGet(option, model)
	assert.Equal(t, "", sqlValues)
}

type StructModel struct {
	Name string
}

func TestPrepareGet(t *testing.T) {
	var ctx context.Context
	core := NewCoreConnection(ctx, Connect)
	option := make(map[string]interface{})
	option["name"] = "test"
	model := &StructModel{}
	sqlFields, _, _ := core.PrepareGet(option, model)
	assert.Equal(t, "name", sqlFields)
}

func TestPrepareUpdate(t *testing.T) {
	var ctx context.Context
	core := NewCoreConnection(ctx, Connect)
	option := make(map[string]interface{})
	paraments := make(map[string]interface{})
	paraments["name"] = "test"
	model := &StructModel{}
	sqlValues, _, _ := core.PrepareUpdate(paraments, option, model)
	assert.Equal(t, "name = $1", sqlValues)
}
