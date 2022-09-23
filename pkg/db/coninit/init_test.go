package coninit

import (
	"fmt"
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var stringConnect = "user=postgres password=123 host=127.0.0.1 port=5432 dbname=postgres sslmode=disable pool_max_conns=10"

func TestGetString(t *testing.T) {
	ctx := context.Background()
	Conn := NewConnectionPool(ctx)
	stringConn := Conn.GetStringConnect("postgres","123","127.0.0.1","postgres",uint16(5432),uint16(10))
	assert.Equal(t, stringConn, stringConnect, "they should be equal")
}
func TestSetConfig(t *testing.T) {
	ctx := context.Background()
	Conn := NewConnectionPool(ctx)
	result := Conn.SetConfig("postgres","123","127.0.0.1","postgres",uint16(5432),uint16(10),true,true)
	assert.Equal(t, result, nil, "they should be equal")
	config := Conn.GetConfig()
	t.Errorf("%+v",config)
	fmt.Sprintf("%+v",config)
	if config.ConnConfig.PreferSimpleProtocol != true {
		t.Error("not set PreferSimpleProtocol")
	} 
}

func TestConnect(t *testing.T) {
	ctx := context.Background()
	Conn := NewConnectionPool(ctx)
	err := Conn.Init("postgres","123","127.0.0.1","postgres",uint16(5432),uint16(10),true,true)
	if err != nil {
		require.Truef(t, false, "Should be able to set Camry : %v.", "not connected")
	}
}