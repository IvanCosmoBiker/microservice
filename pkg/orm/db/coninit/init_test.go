package coninit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var stringConnect = "user=goadmin password=go2021 host=188.225.18.140 port=6432 dbname=cardtest sslmode=disable pool_max_conns=10"

func TestGetString(t *testing.T) {
	ctx := context.Background()
	Conn := NewConnectionPool(ctx)
	stringConn := Conn.GetStringConnect("goadmin", "go2021", "188.225.18.140", "cardtest", uint16(6432), uint16(10))
	assert.Equal(t, stringConn, stringConnect, "they should be equal")
}
func TestSetConfig(t *testing.T) {
	ctx := context.Background()
	Conn := NewConnectionPool(ctx)
	result := Conn.SetConfig("goadmin", "go2021", "188.225.18.140", "cardtest", uint16(6432), uint16(10), uint16(2), true, true)
	assert.Equal(t, result, nil, "they should be equal")
	config := Conn.GetConfig()
	if config.ConnConfig.PreferSimpleProtocol != true {
		t.Error("not set PreferSimpleProtocol")
	}
}

func TestConnect(t *testing.T) {
	ctx := context.Background()
	Conn := NewConnectionPool(ctx)
	err := Conn.Init("postgres", "123", "127.0.0.1", "postgres", uint16(5432), uint16(10), uint16(10), true, true)
	if err != nil {
		require.Truef(t, false, "Should be connect: %s.", "not")
	}
}

func TestCloseConnection(t *testing.T) {
	ctx := context.Background()
	Conn := NewConnectionPool(ctx)
	err := Conn.Init("postgres", "123", "127.0.0.1", "postgres", uint16(5432), uint16(10), uint16(10), true, true)
	if err != nil {
		require.Truef(t, false, "Should be connect: %s.", "not")
	}
	Conn.Close()
	err = Conn.Ping()
	assert.Equal(t, err.Error(), "closed pool", "they should be equal")
}
