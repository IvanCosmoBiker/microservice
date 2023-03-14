package client

import (
	config "ephorservices/config"
	pb "ephorservices/ephorpayment/service"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	Config *config.Config
	Conn   *grpc.ClientConn
	Client pb.PaymentServiceClient
}

func Init(cfg *config.Config) (client *PaymentClient) {
	client = &PaymentClient{
		Config: cfg,
	}
	client.Connect()
	client.Client = pb.NewPaymentServiceClient(client.Conn)
	return
}

func (pc *PaymentClient) Connect() (err error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", pc.Config.Services.EphorPayment.Transport.Grpc.Address, pc.Config.Services.EphorPayment.Transport.Grpc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}
	pc.Conn = conn
	return
}

func (pc *PaymentClient) Shutdown() {
	pc.Conn.Close()
}
