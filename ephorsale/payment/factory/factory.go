package factory

import (
	config "ephorservices/config"
	"ephorservices/ephorsale/payment/interfacePayment"
	sberPay "ephorservices/ephorsale/payment/sberpay"
	modulSbp "ephorservices/ephorsale/payment/sbp/modul"
	vendpay "ephorservices/ephorsale/payment/vendpay"
)

// instance of type banks
var SberPay sberPay.NewSberPayStruct
var vendPay vendpay.NewVendStruct
var sbpModul modulSbp.NewSbpModul

func NewPayment(TypePayment uint8, conf *config.Config) interfacePayment.Payment {
	switch TypePayment {
	case interfacePayment.TypeSber:
		return SberPay.New(conf)
	case interfacePayment.TypeVendPay:
		return vendPay.New(conf)
	case interfacePayment.TypeModulSbp:
		return sbpModul.New(conf)
	}
	return nil
}
