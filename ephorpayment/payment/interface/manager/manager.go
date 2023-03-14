package manager

import (
	"context"
	config "ephorservices/config"
	"ephorservices/ephorpayment/payment/interface/payment"
	pb "ephorservices/ephorpayment/service"
)

type ManagerPayment interface {
	Hold(ctx context.Context, req *pb.Request) (res *pb.Response, err error)
	Debit(ctx context.Context, req *pb.Request) (res *pb.Response, err error)
	Satus(ctx context.Context, req *pb.Request) (res *pb.Response, err error)
	Payment(ctx context.Context, req *pb.Request) (res *pb.Response, err error)
	Return(ctx context.Context, req *pb.Request) (res *pb.Response, err error)
	InitPayment(conf *config.Config)
	SetPayment(mapPayments map[uint8]payment.Payment)
	GetPaymentOfType(tp uint8) (payment.Payment, error)
}
