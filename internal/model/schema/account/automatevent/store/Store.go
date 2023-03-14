package store

import (
	automat_event_model "ephorservices/internal/model/schema/account/automatevent/model"
	model_interface "ephorservices/pkg/orm/model"
	store_parent "ephorservices/pkg/orm/store"
	"math"
)

var (
	Type_OnlineStart           = 0x0000 // начало периода неприрывной работы модема
	Type_OnlineEnd             = 0x0001 // конец периода неприрывной работы модема
	Type_OnlineLast            = 0x0002 // связь с сервером
	Type_Sale                  = 0x0003 // продажа
	Type_PaymentBlocked        = 0x0004 // платежи заблокированы (присылается, если платежи заблокированы слишком долго)
	Type_PaymentUnblocked      = 0x0005 // платежи разблокированы
	Type_PowerUp               = 0x0006 // модем включен/перезагружен
	Type_PowerDown             = 0x0007 // модем выключен
	Type_BadSignal             = 0x0008 // плохой уровень сигнала
	Type_CashlessIdNotFound    = 0x0009 // ТА запросил некорректный номер продукта (STRING:<cashlessId>)
	Type_PriceListNotFound     = 0x000A // Прайс лист не найден (STRING:<deviceId><priceListNumber>)
	Type_SyncConfigError       = 0x000B // Планограммы не совпадают (обратная совместимость)
	Type_PriceNotEqual         = 0x000C // ТА запросил некорректную цену продукта (STRING:<selectId>*<expectedPrice>*<actualPrice>)
	Type_SaleDisabled          = 0x000D // Продажи заблокированы
	Type_SaleEnabled           = 0x000E // Продажи включены
	Type_ConfigEdited          = 0x000F // Конфигурация изменена локально
	Type_ConfigLoaded          = 0x0010 // Конфигурация загружена с сервера
	Type_ConfigLoadFailed      = 0x0011 // Ошибка загрузки конфигурации
	Type_ConfigParseFailed     = 0x0012 // Ошибка формата конфигурации
	Type_ConfigUnknowProduct   = 0x0013 // Неожиданный номер продукта (STRING:<selectId>)
	Type_ConfigUnknowPriceList = 0x0014 // Неожиданный прайс-лист (STRING:<deviceId><priceListNumber>)
	Type_ModemReboot           = 0x0015 // Перезагрузка модема (STRING:<rebootReason>)
	Type_CashCanceled          = 0x0016 // Оплата наличными отменена автоматом
	Type_SaleFailed            = 0x0017 // Ошибка выдачи товара (STRING:<selectId>)
	Type_ProductLow            = 0x0019 // В ячейке закончился товар
	Type_LoadProduct           = 0x001A // Событие о загрузке товаров
	Type_ContainerLow          = 0x001B // В контейнере закончился товар
)

// server
var (
	Type_ServerNoSaleLong       = 0x0101 // Серверное событие, давно не было продаж
	Type_ServerNoAuditLong      = 0x0102 // Серверное событие, давно не было аудита
	Type_ServerNoEncashLong     = 0x0103 // Серверное событие, давно не было инкассации
	Type_ServerNotBillValidator = 0x0104 // Серверное событие, не обнаружен купюроприёмник
	Type_ServerNotCoinChanger   = 0x0105 // Серверное событие, не обнаружен монетоприёмник
	Type_ServerNotCashless      = 0x0106 // Серверное событие, не обнаружен картридер
	Type_ServerNotCoin          = 0x0107 // Серверное событие, давно не было монет
	Type_ServerNotBill          = 0x0108 // Серверное событие, давно не было купюр
	Type_ServerNotCashlessSale  = 0x0109 // Серверное событие, давно не было карт
	Type_ServerConfigEdit       = 0x010A // С момента последней загрузки конфигурации она была изменена
	Type_ServerCommand          = 0x010B // Серверное событие, отправка команды
)

// fiscal
var (
	Type_FiscalUnknownError   = 0x0300 // Необрабатываемая ошибка ФР (дополнительный параметр - код ошибки ФР)
	Type_FiscalLogicError     = 0x0301 // Неожиданное поведение ФР (дополнительный параметр - строка в файле)
	Type_FiscalConnectError   = 0x0302 // Нет связи с ФР
	Type_FiscalPassword       = 0x0303 // Неправильный пароль ФР
	Type_PrinterNotFound      = 0x0304 // Принтер не найден
	Type_PrinterNoPaper       = 0x0305 // В принтере закончилась бумага
	Type_FiscalWrongResponse  = 0x0307 // Некорректный формат ответа
	Type_FiscalFormatError    = 0x0308 // Повреждённый ответ с сервера
	Type_FiscalCompleteNoData = 0x0309 // Чек создан, но реквизиты не получены
	Type_FiscalDiffSetting    = 0x030A // Настройки на сервере не совпадают с настройками в модеме

)

var (
	Category_ERROR    uint8 = 0 // - ошибки (только ошибки)
	Category_SALE     uint8 = 1 // - продажи (только продажи)
	Category_INFO     uint8 = 2 // - автомат (перезагрузка автомата, загрузка настроек и прочее)
	Category_CASHFLOW uint8 = 3 // - движение денег (аудиты, инкассации, размен - в событиях нет, зарезервировано для отчетов)
	Category_TASK     uint8 = 4 // - задачи (в событиях нет, зарезервировано для отчетов)
	Category_MONEY    uint8 = 5 // - прием денег (купюра принята, монет принята, ошибка приема монеты и прочее, сейчас этих событий нет)
	Category_WARNING  uint8 = 6 // - прием денег (купюра принята, монет принята, ошибка приема монеты и прочее, сейчас этих событий нет)
	Category_DEBUG    uint8 = 255
)

var (
	PaymentMethod_None     uint8 = 0
	PaymentMethod_Cash     uint8 = 1
	PaymentMethod_Cashless uint8 = 2
	PaymentMethod_Token    uint8 = 3
)

var (
	TaxRate_NDSNone = 0
	TaxRate_NDS0    = 1
	TaxRate_NDS10   = 2
	TaxRate_NDS18   = 3
	TaxRate_NDS12   = 4
)

var (
	TYPE_NORMAL  uint8 = 0
	TYPE_WARNING uint8 = 1
	TYPE_ALERT   uint8 = 2
)

type StoreAutomatEvent struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreAutomatEvent {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := automat_event_model.New()
	store := &StoreAutomatEvent{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sae *StoreAutomatEvent) GetStructModel(model model_interface.Model) *automat_event_model.AutomatEventModel {
	if model != nil {
		return model.(*automat_event_model.AutomatEventModel)
	}
	return &automat_event_model.AutomatEventModel{}
}

func (sae *StoreAutomatEvent) AddByParams(params map[string]interface{}) (*automat_event_model.AutomatEventModel, error) {
	model, err := sae.Store.AddByParams(params)
	Model := sae.GetStructModel(model)
	return Model, err
}

func (sae *StoreAutomatEvent) SetByParams(params map[string]interface{}) (*automat_event_model.AutomatEventModel, error) {
	model, err := sae.Store.SetByParams(params)
	Model := sae.GetStructModel(model)
	return Model, err
}

func (sae *StoreAutomatEvent) Set(m model_interface.Model) (*automat_event_model.AutomatEventModel, error) {
	model, err := sae.Store.Set(m)
	Model := sae.GetStructModel(model)
	return Model, err
}

func (sae *StoreAutomatEvent) GetOneById(id int) (*automat_event_model.AutomatEventModel, error) {
	model, err := sae.Store.GetOneById(id)
	Model := sae.GetStructModel(model)
	return Model, err
}

func (sae *StoreAutomatEvent) GetNds(tax_rate int, value float64) float64 {
	var nds float64 = 0
	switch tax_rate {
	case TaxRate_NDS18:
		nds = math.Floor((value * float64(20) / float64(120)) * float64(100))
	case TaxRate_NDS10:
		nds = math.Floor((value * float64(10) / float64(110)) * float64(100))
	case TaxRate_NDS12:
		nds = math.Floor((value * float64(12) / float64(112)) * float64(100))
	}
	return nds
}
