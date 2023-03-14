package factory

import (
	config "ephorservices/config"
	"ephorservices/ephorpayment/payment/interface/payment"
	sberPay "ephorservices/ephorpayment/payment/sberpay"
	modulSbp "ephorservices/ephorpayment/payment/sbp/modul"
	vendpay "ephorservices/ephorpayment/payment/vendpay"
)

// instance of type banks
var SberPay sberPay.NewSberPayStruct
var vendPay vendpay.NewVendStruct
var sbpModul modulSbp.NewSbpModul

func NewPayment(TypePayment uint8, conf *config.Config) payment.Payment {
	switch TypePayment {
	case payment.TypeSber:
		return SberPay.New(conf)
	case payment.TypeVendPay:
		return vendPay.New(conf)
	case payment.TypeModulSbp:
		return sbpModul.New(conf)
	}
	return nil
}
