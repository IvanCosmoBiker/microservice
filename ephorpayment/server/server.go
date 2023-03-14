package server

import (
	"context"
	config "ephorservices/config"
	"ephorservices/ephorpayment/payment/factory"
	"ephorservices/ephorpayment/payment/interface/payment"
	pb "ephorservices/ephorpayment/service"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var (
	SBER     uint8 = 1
	VENDPAY  uint8 = 2
	SBPMODUL uint8 = 5
)

var TypePayment = [...]uint8{SBER, VENDPAY, SBPMODUL}

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	Config   *config.Config
	Server   *grpc.Server
	Status   int
	Payments map[uint8]payment.Payment
	Ctx      context.Context
	rmutex   sync.RWMutex
}

func Init(cfg *config.Config) *PaymentServer {
	paymentServer := &PaymentServer{
		Config: cfg,
	}
	paymentServer.InitPayment(cfg)
	return paymentServer
}

func (p *PaymentServer) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", p.Config.Services.EphorPayment.Transport.Grpc.Address, p.Config.Services.EphorPayment.Transport.Grpc.Port))
	if err != nil {
		log.Fatalf("failed connection: %v", err)
		return err
	}
	s := grpc.NewServer()
	p.Server = s
	pb.RegisterPaymentServiceServer(s, p)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
		return err
	}
	return nil
}

func (p *PaymentServer) InitPayment(conf *config.Config) {
	p.Payments = make(map[uint8]payment.Payment, len(TypePayment))
	for _, item := range TypePayment {
		p.Payments[item] = factory.NewPayment(item, conf)
	}
}

func (p *PaymentServer) GetPaymentOfType(tp uint8) (payment.Payment, error) {
	payment, ok := p.Payments[tp]
	if !ok {
		return nil, errors.New("No Avalable Type Payment System")
	}
	return payment, nil
}

func (p *PaymentServer) SetPayment(mapPayments map[uint8]payment.Payment) {
	p.Payments = mapPayments
}

func (p *PaymentServer) Hold(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	res := &pb.Response{}
	payment, err := p.GetPaymentOfType(uint8(req.Type))
	if err != nil {
		res.Status = uint32(transaction.TransactionState_Error)
		res.Order = "none"
		res.Desc = err.Error()
		res.Error = err.Error()
		return res, err
	}
	resultHold := payment.HoldMoney(req)
	if resultHold["status"] == false {
		res.Status = uint32(transaction.TransactionState_Error)
		res.Order = "none"
		res.Desc = parserTypes.ParseTypeInString(resultHold["description"])
		res.Error = parserTypes.ParseTypeInString(resultHold["message"])
		//errors.New("Err hold money")
		return res, nil
	}
	// if tid, exist := resultHold["tid"]; !exist {
	// 	result["ps_tid"] = tid
	// }
	res.Order = parserTypes.ParseTypeInString(resultHold["orderId"])
	res.Desc = parserTypes.ParseTypeInString(resultHold["description"])
	res.Error = ""
	//result["ps_invoice_id"] = resultHold["invoiceId"]
	res.Status = uint32(transaction.TransactionState_MoneyHoldWait)
	return res, nil
}

func (p *PaymentServer) Debit(context.Context, *pb.Request) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (p *PaymentServer) Satus(context.Context, *pb.Request) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (p *PaymentServer) Payment(context.Context, *pb.Request) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (p *PaymentServer) Return(context.Context, *pb.Request) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (p *PaymentServer) ShutDown() error {
	p.Server.GracefulStop()
	return nil
}
