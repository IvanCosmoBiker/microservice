package payment

import (
	pb "ephorservices/ephorpayment/service"
)

var (
	TypeSber     uint8 = 1
	TypeVendPay  uint8 = 2
	TypeLifePay  uint8 = 4
	TypeModulSbp uint8 = 5
	TypeSkbSbp   uint8 = 6
)

type Payment interface {
	HoldMoney(payment *pb.Request) map[string]interface{}
	GetStatusHoldMoney(payment *pb.Request) map[string]interface{}
	DebitHoldMoney(payment *pb.Request) map[string]interface{}
	ReturnMoney(payment *pb.Request) map[string]interface{}
	Timeout()
}
