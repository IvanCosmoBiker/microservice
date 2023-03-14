package server

import (
	"context"
	config "ephorservices/config"
	grpcPayment "ephorservices/ephorpayment/server"
	pb "ephorservices/ephorpayment/service"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

	"google.golang.org/grpc"
)

var serverPayment *grpcPayment.PaymentServer
var Config *config.Config

func init() {
	Config = &config.Config{}
	Config.Services.EphorPayment.Transport.Grpc.Address = "127.0.0.1"
	Config.Services.EphorPayment.Transport.Grpc.Port = "8060"
	Config.Services.EphorPayment.Config.ExecuteMinutes = 5
	Config.Services.EphorPayment.Config.IntervalTime = 5
	serverPayment = grpcPayment.Init(Config)
	go serverPayment.Serve()
}

func TestStartServer(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:8060", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPaymentServiceClient(conn)
	resp, err := client.Hold(ctx, &pb.Request{
		PayType: 1,
		Sum:     300,
		Type:    1,
	})
	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", resp)
}
