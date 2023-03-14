package transaction_struct

import (
	"encoding/json"
	storeAutomatConfig "ephorservices/internal/model/schema/account/automat/config/store"
	storeAutomat "ephorservices/internal/model/schema/account/automat/store"
	storeAutomatEvent "ephorservices/internal/model/schema/account/automatevent/store"
	storeLocation "ephorservices/internal/model/schema/account/automatlocation/store"
	storeCompanyPoint "ephorservices/internal/model/schema/account/companypoint/store"
	storeConfigPriceList "ephorservices/internal/model/schema/account/config/pricelist/store"
	storeConfigProduct "ephorservices/internal/model/schema/account/config/product/store"
	storefr "ephorservices/internal/model/schema/account/fr/store"
	storeRecipeIngredient "ephorservices/internal/model/schema/account/ingredient/recipe/ingredient/store"
	storeRecipe "ephorservices/internal/model/schema/account/ingredient/recipe/store"
	storeIngredient "ephorservices/internal/model/schema/account/ingredient/store"
	storeWareBarcode "ephorservices/internal/model/schema/account/ware/barcode/store"
	storeWare "ephorservices/internal/model/schema/account/ware/store"
	storeWareFlowIngredient "ephorservices/internal/model/schema/account/wareflow/ingredient/store"
	storeWareFlowProduct "ephorservices/internal/model/schema/account/wareflow/product/store"
	storeWareFlow "ephorservices/internal/model/schema/account/wareflow/store"
	storeCommand "ephorservices/internal/model/schema/main/command/store"
	dateTime "ephorservices/pkg/datetime"
	"ephorservices/pkg/orm/request"
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
	Command_id   int
}
type Payment struct {
	Message                                                                                                                         string
	Status, PayType, TokenType, PaymentId, Sum, CurrensyCode, TypeSale, DebitSum                                                    int
	Type                                                                                                                            uint8
	OrderId, OperationId, InvoiceId, Tid, SecretKey, KeyPayment, Service_id, TidPaymentSystem, HostPaymentSystem, SbolBankInvoiceId string
	Login, Password, MerchantId, GateWay, Token, UserPhone, ReturnUrl, DeepLink, SbpPoint, Language, Description                    string
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
	CancelCheck         int
}
type Fiscal struct {
	QrCode        string
	Events        []string
	FiscalRequest struct {
		Request json.RawMessage
	}
	NeedFiscal, OnlyFiscal, Send                                                  bool
	Config                                                                        ConfigFR
	SumCheck                                                                      float64
	Message, Signature, ResiptId, AuthToken                                       string
	Status, Code, StatusCode, Fiscalization, Qr, FfDisableCashless, FfDisableCash int
	Fields                                                                        struct {
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
type Product struct {
	Name           string
	Payment_device string
	Price_list     int32
	Type           int32
	Ware_id        int32
	Select_id      string
	Value          float64
	Price          float64
	Tax_rate       int32
	Quantity       int64
	Marking        int
	Barcode        string
	Fiscalization  bool
}

type Transaction struct {
	Status         int
	Error          string
	KeyReplay      int
	DateTime       *dateTime.DateTime
	Date           string
	ChannelMessage chan []byte
	TimeOut        chan bool
	Close          chan []byte
	Config         Config
	Payment        Payment
	Fiscal         Fiscal
	Point          Point
	TaxSystem      TaxSystem
	Sum            int32
	Products       []*Product
	Stores         struct {
		StoreCommand            *storeCommand.StoreCommand
		StoreAutomat            *storeAutomat.StoreAutomat
		StoreAutomatLocation    *storeLocation.StoreAutomatLocation
		StoreCompanyPoint       *storeCompanyPoint.StoreCompanyPoint
		StoreAutomatEvent       *storeAutomatEvent.StoreAutomatEvent
		StoreWare               *storeWare.StoreWare
		StoreWareBarcode        *storeWareBarcode.StoreWareBarcode
		StoreAutomatConfig      *storeAutomatConfig.StoreAutomatConfig
		StoreConfigProduct      *storeConfigProduct.StoreConfigProduct
		StoreConfigPriceList    *storeConfigPriceList.StoreConfigPriceList
		StoreFr                 *storefr.StoreFr
		StoreWareFlow           *storeWareFlow.StoreWareFlow
		StoreWareFlowProduct    *storeWareFlowProduct.StoreWareFlowProduct
		StoreWareFlowIngredient *storeWareFlowIngredient.StoreWareFlowIngredient
		StoreRecipe             *storeRecipe.StoreRecipe
		StoreRecipeIngredient   *storeRecipeIngredient.StoreRecipeIngredient
		StoreIngredient         *storeIngredient.StoreIngredient
	}
}

func InitTransaction() *Transaction {
	date, _ := dateTime.Init()
	return &Transaction{
		Status:         0,
		TimeOut:        make(chan bool),
		Close:          make(chan []byte),
		ChannelMessage: make(chan []byte),
		Products:       make([]*Product, 0, 1),
		DateTime:       date,
	}
}

func (t *Transaction) InitStores(accountId int) {
	t.Stores.StoreCommand = storeCommand.New()
	t.Stores.StoreAutomat = storeAutomat.New(accountId)
	t.Stores.StoreAutomatLocation = storeLocation.New(accountId)
	t.Stores.StoreCompanyPoint = storeCompanyPoint.New(accountId)
	t.Stores.StoreAutomatEvent = storeAutomatEvent.New(accountId)
	t.Stores.StoreWare = storeWare.New(accountId)
	t.Stores.StoreWareBarcode = storeWareBarcode.New(accountId)
	t.Stores.StoreAutomatConfig = storeAutomatConfig.New(accountId)
	t.Stores.StoreFr = storefr.New(accountId)
	t.Stores.StoreWareFlow = storeWareFlow.New(accountId)
	t.Stores.StoreWareFlowProduct = storeWareFlowProduct.New(accountId)
	t.Stores.StoreWareFlowIngredient = storeWareFlowIngredient.New(accountId)
	t.Stores.StoreRecipe = storeRecipe.New(accountId)
	t.Stores.StoreRecipeIngredient = storeRecipeIngredient.New(accountId)
	t.Stores.StoreConfigProduct = storeConfigProduct.New(accountId)
	t.Stores.StoreIngredient = storeIngredient.New(accountId)
	t.Stores.StoreConfigPriceList = storeConfigPriceList.New(accountId)
}

func (t *Transaction) NewRequest() *request.Request {
	return request.New()
}

func (t Transaction) GetDescriptionCodeCooler(code int) string {
	stringCode := "Нет"
	switch code {
	case IceboxStatus_Drink:
		stringCode = `Замок открыт, заберите продукты`
		return stringCode
	case IceboxStatus_Icebox:
		stringCode = `Замок открыт, заберите продукты`
		return stringCode
	case IceboxStatus_End:
		stringCode = `Замок закрыт.Спасибо за покупку`
		return stringCode
	case TransactionState_WaitFiscal:
		stringCode = `Фискализируем продажу`
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
	case TransactionState_WaitFiscal:
		stringCode = `Фискализируем продажу`
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

func Destroy(tran *Transaction) {
	tran = nil
}
