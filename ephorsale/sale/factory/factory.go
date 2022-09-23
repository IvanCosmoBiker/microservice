package factory

import (
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal"
	"ephorservices/ephorsale/payment"
	automatPostPaid "ephorservices/ephorsale/sale/automatpostpaid"
	automatPrePaid "ephorservices/ephorsale/sale/automatprepaid"
	coolerPrePaid "ephorservices/ephorsale/sale/coolerprepaid"
	"ephorservices/ephorsale/sale/interfaceSale"
	transaction "ephorservices/ephorsale/transaction"
	rabbitMq "ephorservices/pkg/rabbitmq"
)

var automatPrePay automatPrePaid.NewSaleAutomatPrePaid
var automatPostPay automatPostPaid.NewSaleAutomatPostPaid
var coolerPrePay coolerPrePaid.NewSaleCoolerPrePaid

func GetSale(Type, PayType uint8, Rabbit *rabbitMq.Manager, conf *config.Config, fiscalM *fiscal.FiscalManager, paymentM *payment.PaymentManager, dispether *transaction.TransactionDispetcher) interfaceSale.Sale {
	switch Type {
	case interfaceSale.TypeCoffee,
		interfaceSale.TypeSnack,
		interfaceSale.TypeHoreca,
		interfaceSale.TypeSodaWater,
		interfaceSale.TypeMechanical,
		interfaceSale.TypeComb:
		if PayType == interfaceSale.Type_Prepayment {
			return automatPrePay.New(conf, Rabbit, fiscalM, paymentM, dispether)
		}
		if PayType == interfaceSale.Type_PostPaid {
			return automatPostPay.New(conf, Rabbit, fiscalM, paymentM, dispether)
		}
		fallthrough
	case interfaceSale.TypeCooler:
		if PayType == interfaceSale.Type_Prepayment {
			return coolerPrePay.New(conf, Rabbit, fiscalM, paymentM, dispether)
		}
	}
	return nil
}
