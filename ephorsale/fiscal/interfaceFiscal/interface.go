package interfaceFiscal

import (
	transaction "ephorservices/ephorsale/transaction"
)

const (
	Fr_None                  uint8 = 0
	Fr_PayOnlineFA           uint8 = 1
	Fr_KaznachejFA           uint8 = 2
	Fr_RPSystem1FA           uint8 = 3
	Fr_TerminalFA            uint8 = 4
	Fr_OrangeData            uint8 = 5
	Fr_ChekOnline            uint8 = 6
	Fr_EphorOrangeData       uint8 = 7
	Fr_EphorOnline           uint8 = 8
	Fr_NanoKassa             uint8 = 9
	Fr_ServerOrangeData      uint8 = 10
	Fr_EphorServerOrangeData uint8 = 11
	Fr_OFD                   uint8 = 12
	Fr_ServerNanoKassa       uint8 = 13
	Fr_ProSistem             uint8 = 14
	Fr_CheckBox              uint8 = 15
)

// -- Status Fiscal -- //
const (
	Status_None         uint8 = 0 // продажа за 0 рублей или касса отключена
	Status_Complete     uint8 = 1 // чек создан успешно, реквизиты получены
	Status_InQueue      uint8 = 2 // чек добавлен в очередь, реквизиты не получены
	Status_Unknown      uint8 = 3 // результат постановки чека в очередь не известен
	Status_Error        uint8 = 4 // ошибка создания чека
	Status_Overflow     uint8 = 5 // очередь отложенной регистрации переполнена
	Status_Manual       uint8 = 6 // чек фискализирован вручную
	Status_Need         uint8 = 7 // чек необходимо фискализировать
	Status_MAX_CHECK    uint8 = 8
	Status_OFF_FR       uint8 = 9  // оключение фискализации со стороны клиента
	Status_OFF_DA       uint8 = 10 // отключение фискализации для безналичной оплаты
	Status_OFF_CA       uint8 = 11 // отключение фискализации для оплаты за наличку
	Status_Check_Cancel uint8 = 12 // отмена чека
)

const (
	Fr_Disable_Cash     = 1
	Fr_Disable_Cashless = 1
)

// -- Error Fiscal -- //
const (
	Status_Fr_Fr_InQueue = 101 // чек добавлен в очередь, реквизиты не получены
	Status_Fr_Unknown    = 102 // результат постановки чека в очередь не известен
	Status_Fr_Error      = 103 // ошибка создания чека
	Status_Fr_Overflow   = 104 // очередь отложенной регистрации переполнена
	Status_Fr_MAX        = 105 // превышена максимальная сумма чека
)

type Fiscal interface {
	SendCheck(data *transaction.Transaction) map[string]interface{}
	GetStatus(data *transaction.Transaction) map[string]interface{}
	GetQrPicture(data *transaction.Transaction) string
	GetQrUrl(data *transaction.Transaction) string
}
