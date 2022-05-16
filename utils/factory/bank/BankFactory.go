package bank

import (
    sberBank "bank/sber"
    vendpay "bank/vendpay"
    interfaceBank "interface/bankinterface"
)
// instance of type banks
var bankSber sberBank.NewSberStruct
var bankVendPay vendpay.NewVendStruct

func GetBank(bank int) (interfaceBank.Bank) {
  switch bank {
        case interfaceBank.TypeSber:
            return bankSber.NewBank()
            fallthrough
        case interfaceBank.TypeVendPay:
            return bankVendPay.NewBank()
  }
  return nil
}



