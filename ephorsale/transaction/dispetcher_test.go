package transaction

import (
	"context"
	config "ephorservices/config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func init() {
	conf := &config.Config{}
	ctx := context.Background()
	New(conf, ctx)
}

func TestTwoTransaction(t *testing.T) {
	var wg sync.WaitGroup
	tran1 := Dispetcher.NewTransaction()
	noise1 := Dispetcher.Random(10000, 500000)
	tran1.Config.Tid = noise1
	tran1.Config.Noise = noise1
	Dispetcher.AddChannel(tran1.Config.Tid, tran1)
	wg.Add(1)
	tran2 := Dispetcher.NewTransaction()
	noise2 := Dispetcher.Random(10000, 500000)
	tran2.Config.Tid = noise2
	tran2.Config.Noise = noise2
	Dispetcher.AddChannel(tran2.Config.Tid, tran2)
	wg.Add(1)
	go func() {
		message, ok := <-tran1.ChannelMessage
		if ok {
			fmt.Printf("%s", string(message))
			assert.Equal(t, "transaction1", string(message))
		}
		wg.Done()
		return
	}()
	go func() {
		message, ok := <-tran2.ChannelMessage
		if ok {
			fmt.Printf("%s", string(message))
			assert.Equal(t, "transaction2", string(message))
		}
		wg.Done()
		return
	}()
	Dispetcher.Send(noise1, []byte("transaction1"))
	Dispetcher.Send(noise2, []byte("transaction2"))
	wg.Wait()
}
