package transaction

import (
	"time"
)

var (
	Prepayment = 1 // предоплата товара
	Postpaid   = 2 // постоплата товара
)

var (
	IceboxStatus_Drink  = 1 // выдача напитков
	IceboxStatus_Icebox = 2 // дверь открыта
	IceboxStatus_End    = 3 // выдача завершена
)

var (
	VendState_Session   = 1 //[11]PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
	VendState_Approving = 2 //[14] Продукт выбран. Ожидание оплаты.
	VendState_Vending   = 3 //[12]PAY_OK_AUTOMAT_PREPARE Оплата успешна, ТА готовит продукт
	VendState_VendOk    = 4 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
	VendState_VendError = 5 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
)

var (
	VendError_VendFailed       = 769 //769 Ошибка выдачи продукта
	VendError_SessionCancelled = 770 //770
	VendError_SessionTimeout   = 771 //771
	VendError_WrongProduct     = 772 //772
	VendError_VendCancelled    = 773 //773
	VendError_ApprovingTimeout = 774 //774
)

var (
	TransactionState_Idle             = 0 // Transaction Idle
	TransactionState_MoneyHoldStart   = 1 // создали транзакцию банка
	TransactionState_MoneyHoldWait    = 2 // ожидает ответ от банка
	TransactionState_VendSessionStart = 3 // PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
	TransactionState_VendSession      = 4 //[11]PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
	TransactionState_VendApproving    = 5 //[14] Продукт выбран. Ожидание оплаты.
	TransactionState_Vending          = 6 //[12]PAY_OK_AUTOMAT_PREPARE Оплата успешна, ТА готовит продукт
	TransactionState_MoneyDebitStart  = 8
	TransactionState_MoneyDebitWait   = 9
	TransactionState_MoneyDebitOk     = 10
	TransactionState_VendOk           = 11 // приготовил продукт
	TransactionState_ReturnMoney      = 12
	TransactionState_SberPay          = 13 // инициализация платежа Сбербанк онлайн
	TransactionState_WaitFiscal       = 14
	TransactionState_Error            = 120
	TransactionState_ErrorTimeOut     = 121
	TransactionState_EndClient        = 122
)
var (
	TypeTokenGooglePay      = 1
	TypeTokenApplePay       = 2
	TypeTokenSamsungPay     = 3
	TypeTokenSberPayWeb     = 4
	TypeTokenSberPayAndroid = 5
	TypeTokenSberPayiOS     = 6
)

var (
	Sale_Cash     = 1
	Sale_Cashless = 2
)

var (
	TaxSystem_OSN    = 0x01 // Общая ОСН
	TaxSystem_USND   = 0x02 // Упрощенная доход
	TaxSystem_USNDMR = 0x04 // Упрощенная доход минус расход
	TaxSystem_ENVD   = 0x08 // Единый налог на вмененный доход
	TaxSystem_ESN    = 0x10 // Единый сельскохозяйственный налог
	TaxSystem_Patent = 0x20 // Патентная система налогообложения
)
var (
	TaxRate_NDSNone = 0
	TaxRate_NDS0    = 1
	TaxRate_NDS10   = 2
	TaxRate_NDS18   = 3
)

type Config struct {
	Tid          int
	AccountId    int
	AutomatId    int
	Imei         string
	TokenType    int
	Noise        int
	DeviceType   int
	CurrensyCode int
}

type ConfigFR struct {
	Id                  int64
	Name                string
	Type                uint8
	Dev_interface       int
	AutomatNumber       int
	Login               string
	Password            string
	Phone               string
	Email               string
	Dev_addr            string
	Dev_port            int
	Ofd_addr            string
	Ofd_port            int
	Inn                 string
	Auth_public_key     string
	Auth_private_key    string
	Sign_private_key    string
	Param1              string
	Use_sn              int
	Add_fiscal          int
	Id_shift            string
	Fr_disable_cash     int
	Fr_disable_cashless int
	Ffd_version         int
	MaxSum              int
	QrFormat            int
}

type Payment struct {
	Status, PayType, TokenType, PaymentId, Sum, CurrensyCode, TypeSale, DebitSum                                                    int
	Type                                                                                                                            uint8
	OrderId, OperationId, InvoiceId, Tid, SecretKey, KeyPayment, Service_id, TidPaymentSystem, HostPaymentSystem, SbolBankInvoiceId string
	Login, Password, MerchantId, GateWay, Token, UserPhone, ReturnUrl, DeepLink, SbpPoint, Language, Description                    string
}

type Fiscal struct {
	Config                                                                ConfigFR
	Message, Status                                                       string
	Code, StatusCode, Fiscalization, Qr, FfDisableCashless, FfDisableCash int
	Fields                                                                struct {
		Fp, Fn, DateFisal string
		Fd                float64
	}
}

type Point struct {
	Address, PointName string
	Id                 int
}

type TaxSystem struct {
	Type    int
	TaxRate int
}

type Transaction struct {
	Date           string
	Status         int
	ChannelMessage chan []byte
	TimeOut        chan bool
	Config         Config
	Payment        Payment
	Fiscal         Fiscal
	Point          Point
	TaxSystem      TaxSystem
	Products       []map[string]interface{}
}

func InitTransaction() *Transaction {
	return &Transaction{
		Status:         0,
		TimeOut:        make(chan bool),
		ChannelMessage: make(chan []byte),
	}
}

func (t *Transaction) GateDateTimeMoscow() string {
	timeNow := time.Now()
	timeNow.Add(-72 * time.Hour)
	resultTime := timeNow.Format("2006-01-02 15:04:05")
	return resultTime
}

func (t Transaction) GetDescriptionCodeCooler(code int) string {
	stringCode := "Нет"
	switch code {
	case IceboxStatus_Drink:
		stringCode = `Подойдите к кофе машине. Подставьте стакан и выберите напиток`
		return stringCode
	case IceboxStatus_Icebox:
		stringCode = `Замок открыт, заберите продукты`
		return stringCode
	case IceboxStatus_End:
		stringCode = `Замок закрыт. Счастливого пути`
		return stringCode
	}
	return stringCode
}

func (t Transaction) GetStatusServerCooler(status int) int {
	switch status {
	case IceboxStatus_Drink:
		return TransactionState_VendSession
	case IceboxStatus_Icebox:
		return TransactionState_VendSession
	case IceboxStatus_End:
		return TransactionState_VendOk
	}
	return TransactionState_Error
}

func (t Transaction) GetDescriptionCode(code int) string {
	stringCode := "Нет"
	switch code {
	case VendState_Session:
		stringCode = `Оплата успешна, ожидание нажатия пользователем кнопки на ТА`
		return stringCode
	case VendState_Approving:
		stringCode = `Продукт выбран. Ожидание оплаты.`
		return stringCode
	case VendState_Vending:
		stringCode = `Оплата успешна, ТА готовит продукт`
		return stringCode
	case VendState_VendOk:
		stringCode = `Оплата успешна, ТА приготовил продукт`
		return stringCode
	case VendState_VendError:
		stringCode = `Ошибка`
		return stringCode
	case TransactionState_ErrorTimeOut:
		stringCode = `Время ответа от автомата истекло`
		return stringCode
	case IceboxStatus_Drink:
		stringCode = `Подойдите к кофе машине. Подставьте стакан и выберите напиток`
		return stringCode
	}
	return stringCode
}

func (t Transaction) GetStatusServer(status int) int {
	switch status {
	case VendState_Session:
		return TransactionState_VendSession
	case VendState_Approving:
		return TransactionState_VendApproving
	case VendState_Vending:
		return TransactionState_Vending
	case VendState_VendOk:
		return TransactionState_VendOk
	case VendState_VendError:
		return TransactionState_Error
	}
	return TransactionState_Error
}

func (t Transaction) GetDescriptionErr(err int) string {
	stringErr := "Нет"
	switch err {
	case VendError_VendFailed:
		stringErr = `Ошибка выдачи товара на автомате`
		return stringErr
	case VendError_SessionCancelled:
		stringErr = `Продажа отменена автоматом`
		return stringErr
	case VendError_SessionTimeout:
		stringErr = `Время ожидание выбора товара на автомате истекло`
		return stringErr
	case VendError_WrongProduct:
		stringErr = `Выбранный на автомате товар не совпадает с оплаченым`
		return stringErr
	case VendError_VendCancelled:
		stringErr = `Выдача товара отменена автоматом`
		return stringErr
	case VendError_ApprovingTimeout:
		stringErr = `Время ожидание оплаты истекло`
		return stringErr
	case TransactionState_ErrorTimeOut:
		stringErr = `Время ответа от автомата истекло`
		return stringErr
	}
	return stringErr
}
