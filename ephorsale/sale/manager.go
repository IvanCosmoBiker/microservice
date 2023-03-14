package sale

import (
	"context"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal"
	"ephorservices/ephorsale/fiscal/interface/fr"
	"ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/payment/interface/manager"
	sale_http_service "ephorservices/ephorsale/sale/api/service_http"
	sale_mqtt_service "ephorservices/ephorsale/sale/api/service_mqtt"
	automatPostPaid "ephorservices/ephorsale/sale/automatpostpaid"
	automatPrePaid "ephorservices/ephorsale/sale/automatprepaid"
	coolerPrePaid "ephorservices/ephorsale/sale/coolerprepaid"
	"ephorservices/ephorsale/sale/interfaceSale"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	transaction_struct "ephorservices/ephorsale/transaction/transaction_struct"
	"ephorservices/ephorsale/transport"
	automatEventStore "ephorservices/internal/model/schema/account/automatevent/store"
	wareStore "ephorservices/internal/model/schema/account/ware/store"
	logger "ephorservices/pkg/logger"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"fmt"
	"runtime"

	"golang.org/x/sync/errgroup"
)

// Sale
var (
	automatPrePay  automatPrePaid.NewSaleAutomatPrePaid
	automatPostPay automatPostPaid.NewSaleAutomatPostPaid
	coolerPrePay   coolerPrePaid.NewSaleCoolerPrePaid
)

// Fiscal

// Payment

type SaleManager struct {
	Status  int
	Fiscal  *fiscal.FiscalManager
	Payment manager.ManagerPayment
	ctx     context.Context
	cfg     *config.Config
}

var Manager *SaleManager

func New(ctx context.Context, conf *config.Config) (sale *SaleManager, err error) {
	sale = &SaleManager{
		ctx: ctx,
		cfg: conf,
	}
	sale.Fiscal, err = fiscal.New(ctx, conf.Services.EphorFiscal.ExecuteMinutes, conf.Services.EphorFiscal.SleepSecond, conf.Services.EphorFiscal.PathCert, conf.Services.EphorFiscal.ResponseUrl, conf.Debug)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	sale.Payment = payment.New(conf.Debug, ctx)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	Manager = sale
	return sale, nil
}

func (sm *SaleManager) getSale(typeDevice, payType uint8) interfaceSale.Sale {
	switch typeDevice {
	case interfaceSale.TypeCoffee,
		interfaceSale.TypeSnack,
		interfaceSale.TypeHoreca,
		interfaceSale.TypeSodaWater,
		interfaceSale.TypeMechanical,
		interfaceSale.TypeComb:
		if payType == interfaceSale.Type_Prepayment {
			return automatPrePay.New(sm.cfg.Transport.Mqtt.ExecuteTimeSeconds, sm.cfg.Debug)
		}
		if payType == interfaceSale.Type_PostPaid {
			return automatPostPay.New(sm.cfg.Transport.Mqtt.ExecuteTimeSeconds, sm.cfg.Debug)
		}
	case interfaceSale.TypeCooler:
		if payType == interfaceSale.Type_Prepayment {
			return coolerPrePay.New(sm.cfg.Transport.Mqtt.ExecuteTimeSeconds, sm.cfg.Debug)
		}
	}
	return nil
}

func (sm *SaleManager) CancelTransaction(tran *transaction_struct.Transaction) {

}

/*
*
The method StartSale triggers sale with include payment, conversation with device and fiscalization of sale
*
*/
func (sm *SaleManager) StartSale(tran *transaction_struct.Transaction) {
	defer sm.FinishSale(tran)
	sale := sm.getSale(uint8(tran.Config.DeviceType), uint8(tran.Payment.PayType))
	sale.SetFiscalisation(sm.StartFiscal)
	resultSale := sale.Sale(tran)
	fmt.Printf("SALE IS - %+v", sale)
	fmt.Printf("SALE IS - %+v", resultSale)
	if resultSale["status"].(int) == transaction_struct.TransactionState_Error {
		tran.Status = int(resultSale["status"].(int))
	}
}

/*
*
The method only triggers fiscalization. The sale was on the device side
*
*/
func (sm *SaleManager) StartFiscal(tran *transaction_struct.Transaction) {
	fmt.Printf("___%+v\n", tran)
	if tran.Fiscal.OnlyFiscal && tran.Fiscal.Send {
		fmt.Println(" SEND!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		defer func() {
			if err := recover(); err != nil {
				er := fmt.Errorf("%v", err)
				tran.Status = int(fr.Status_Error)
				tran.Fiscal.Status = int(fr.Status_Error)
				tran.Error = er.Error()
				tran.Fiscal.Message = tran.Error
				return
			}
		}()
		defer sm.Fiscal.Send("http", tran)
	}
	if !tran.Fiscal.NeedFiscal {
		tran.Status = int(fr.Status_None)
		tran.Fiscal.Status = int(fr.Status_None)
		tran.Error = "Фискализация будет через устройство Ephor"
		return
	}
	// Check all parametrs for fiscalisation sale(s)
	resultStartFiscal, frKass := sm.Fiscal.DataVerification(tran)
	fmt.Printf("%+v", tran)
	if resultStartFiscal["status"].(bool) == false {
		tran.Status = parserTypes.ParseTypeInterfaceToInt(resultStartFiscal["fr_status"])
		tran.Fiscal.Status = parserTypes.ParseTypeInterfaceToInt(resultStartFiscal["fr_status"])
		tran.Error = parserTypes.ParseTypeInString(resultStartFiscal["message"])
		tran.Fiscal.Message = tran.Error
		return
	}
	statusSendCheck := frKass.SendCheck(tran)
	fmt.Printf("SEND_CHECK!!!***%+v\n", statusSendCheck)
	if parserTypes.ParseTypeInterfaceToUint8(statusSendCheck["status"]) != fr.Status_InQueue {
		tran.Status = int(fr.Status_Error)
		tran.Fiscal.Status = int(fr.Status_Error)
		tran.Error = parserTypes.ParseTypeInString(statusSendCheck["f_desc"])
		tran.Fiscal.Message = tran.Error
		return
	}
	resultStatusCheck := frKass.GetStatus(tran)
	fmt.Printf("STATUS_CHECK!!!***%+v\n", resultStatusCheck)
	if parserTypes.ParseTypeInterfaceToUint8(resultStatusCheck["status"]) != fr.Status_Complete {
		tran.Status = int(fr.Status_Error)
		tran.Fiscal.Status = int(fr.Status_Error)
		tran.Error = parserTypes.ParseTypeInString(resultStatusCheck["f_desc"])
		tran.Fiscal.Message = tran.Error
		return
	}
	tran.Status = int(fr.Status_Complete)
	tran.Fiscal.Status = int(fr.Status_Complete)
	tran.Error = "Нет ошибок"
	tran.Fiscal.Message = tran.Error
	if tran.Fiscal.Config.QrFormat == 1 {
		tran.Fiscal.QrCode = frKass.GetQrPicture(tran)
	} else {
		tran.Fiscal.QrCode = frKass.GetQrUrl(tran)
	}
	return
}

/*
*
The method only triggers payment without fiscalisation. The sale was on the device side
*
*/
func (sm *SaleManager) StartPayment(tran *transaction_struct.Transaction) {
	//Sale := sm.getSale(tran.Config.DeviceType, tran.Payment.PayType)
	//go Sale.Payment(tran)
	//sm.AddSale(tran)
}

func (sm *SaleManager) AddSale(tran *transaction_struct.Transaction) error {
	sm.UpdateSaleAutomat(tran)
	sm.DeductionProducts(tran)
	sm.AddEvent(tran)
	return nil
}

func (sm *SaleManager) FinishSale(tran *transaction_struct.Transaction) {
	if tran.Status != transaction_struct.TransactionState_Error && uint8(tran.Config.DeviceType) == interfaceSale.TypeCooler {
		sm.AddSale(tran)
	}
	logger.Log.Info("REMOVE TRANSACTION")
	result := transaction_dispetcher.Dispetcher.RemoveTransaction(tran.Config.Tid)
	if !result {
		logger.Log.Error("REMOVE TRANSACTION IS FALSE")
	}
	result = transaction_dispetcher.Dispetcher.RemoveReplayProtection(tran.KeyReplay)
	if !result {
		logger.Log.Error("REMOVE PROTECTION IS FALSE")
	}
	transaction_struct.Destroy(tran)
	runtime.Goexit()
}

func (sm *SaleManager) InitApi(conf *config.Config, ctx context.Context, errorGroup *errgroup.Group, Transpoprt *transport.Transport) {
	serviceMqttSale := sale_mqtt_service.New(ctx)
	serviceMqttSale.InitApi(conf, ctx, errorGroup)

	serviceHttpSale := sale_http_service.New()
	serviceHttpSale.SaleHandler = sm.StartSale
	serviceHttpSale.FiscalHandler = sm.StartFiscal
	serviceHttpSale.InitApi(Transpoprt.RequestManager)
}

func (sm *SaleManager) UpdateSaleAutomat(tran *transaction_struct.Transaction) {
	automat, err := tran.Stores.StoreAutomat.GetOneById(tran.Config.AutomatId)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if automat.Id == 0 {
		return
	}
	automat.Type_nosale.Scan(int32(automatEventStore.TYPE_NORMAL))
	automat.Last_sale.Scan(tran.Date)
	for _, product := range tran.Products {
		fmt.Printf("PRODUCT %+v", product)
		var value float64 = 0
		if product.Value == value {
			value = product.Price
		} else {
			value = product.Value
		}
		if product.Payment_device == "DA" || product.Payment_device == "DB" {
			automat.Last_cashless.Scan(tran.Date)
			automat.Now_cashless_val.Int64 += int64(value)
			automat.Now_cashless_num.Int32 += int32(1)
		}
		if product.Payment_device == "TA" {
			automat.Now_token_val.Int64 += int64(value)
			automat.Now_token_num.Int32 += int32(1)
		}
		if product.Payment_device == "CA" {
			automat.Now_cash_val.Int64 += int64(value)
			automat.Now_cash_num.Int32 += int32(1)
		}
	}
	_, err = tran.Stores.StoreAutomat.Set(automat)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
}

func (sm *SaleManager) DeductionProducts(tran *transaction_struct.Transaction) {
	var eventType int
	wareflow, err := tran.Stores.StoreWareFlow.GetLatest(tran.Config.AutomatId)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	if wareflow.Id == 0 {
		return
	}
	for _, product := range tran.Products {
		if product.Type == int32(wareStore.Type_Snack) {
			eventType, err = tran.Stores.StoreWareFlowProduct.Deduction(wareflow.Id, product.Select_id, int(product.Ware_id), int32(product.Quantity/1000))
			if err != nil {
				logger.Log.Error(err.Error())
				continue
			}
			if eventType != 0 {
				event := make(map[string]interface{})
				event["automat_id"] = tran.Config.AutomatId
				event["type"] = automatEventStore.Type_ProductLow
				event["date"] = tran.DateTime.Now()
				event["name"] = fmt.Sprintf("%s*%s", product.Select_id, product.Name)
				tran.Stores.StoreAutomatEvent.AddByParams(event)
			}
		} else if product.Type == int32(wareStore.Type_Recipe) {
			sm.DeductResipe(wareflow.Id, product, tran)
		}
	}
}

func (sm *SaleManager) DeductResipe(wareflowId int, product *transaction_struct.Product, tran *transaction_struct.Transaction) {
	var eventType int
	var err error
	reqAutomatConfig := tran.NewRequest()
	reqAutomatConfig.AddFilterParam("automat_id", reqAutomatConfig.Operator.OperatorEqual, true, tran.Config.AutomatId)
	reqAutomatConfig.AddFilterParam("to_date", reqAutomatConfig.Operator.OperatorEqual, true)
	automatConfig, errAc := tran.Stores.StoreAutomatConfig.GetOneBy(reqAutomatConfig)
	if errAc != nil {
		return
	}
	reqProductConfig := tran.NewRequest()
	reqProductConfig.AddFilterParam("config_id", reqProductConfig.Operator.OperatorEqual, true, automatConfig.Config_id.Int32)
	reqProductConfig.AddFilterParam("select_id", reqProductConfig.Operator.OperatorEqual, true, product.Select_id)
	configProduct, errCp := tran.Stores.StoreConfigProduct.GetOneBy(reqProductConfig)
	if errCp != nil {
		return
	}
	recipe, errR := tran.Stores.StoreRecipe.GetOneById(int(configProduct.Recipe_id.Int32))
	if errR != nil {
		return
	}
	if recipe.Id == 0 {
		return
	}
	reqRecipeIngredient := tran.NewRequest()
	reqRecipeIngredient.AddFilterParam("recipe_id", reqRecipeIngredient.Operator.OperatorEqual, true, recipe.Id)
	recipeIngredients, errRi := tran.Stores.StoreRecipeIngredient.Get(reqRecipeIngredient)
	if errRi != nil {
		return
	}
	for _, ingredientRecipe := range recipeIngredients {
		if ingredientRecipe.Id == 0 {
			continue
		}
		eventType, err = tran.Stores.StoreWareFlowIngredient.Deduction(wareflowId, int(ingredientRecipe.Ingredient_id.Int32), int32(ingredientRecipe.Count.Int32))
		if err != nil {
			logger.Log.Error(err.Error())
			continue
		}
		ingredient, errIn := tran.Stores.StoreIngredient.GetOneById(int(ingredientRecipe.Ingredient_id.Int32))
		if errIn != nil {
			continue
		}
		if ingredient.Id == 0 {
			continue
		}
		if eventType != 0 {
			event := make(map[string]interface{})
			event["automat_id"] = tran.Config.AutomatId
			event["type"] = automatEventStore.Type_ContainerLow
			event["date"] = tran.DateTime.Now()
			event["name"] = ingredient.Name.String
			tran.Stores.StoreAutomatEvent.AddByParams(event)
		}
	}
}

func (sm *SaleManager) AddEvent(tran *transaction_struct.Transaction) {
	dateFiscal, errDateFiscal := tran.DateTime.ParseDateAndSubtractHour(tran.Fiscal.Fields.DateFisal)
	var date string
	for _, item := range tran.Products {
		for i := 0; i < int(item.Quantity/1000); i++ {
			if i == 0 {
				date = tran.Date
			} else {
				date, _ = tran.DateTime.AddSeconds(date, 1)
			}
			var value float64 = 0
			if item.Value == value {
				value = item.Price
			} else {
				value = item.Value
			}
			entry := make(map[string]interface{})
			entry["account_id"] = tran.Config.AccountId
			entry["date"] = date
			entry["automat_id"] = tran.Config.AutomatId
			entry["modem_date"] = date
			if errDateFiscal == nil {
				entry["fiscal_date"] = dateFiscal
			}
			entry["update_date"] = date
			entry["type"] = automatEventStore.Type_Sale
			entry["category"] = automatEventStore.Category_SALE
			entry["select_id"] = item.Select_id
			entry["ware_id"] = item.Ware_id
			entry["name"] = item.Name
			entry["payment_device"] = item.Payment_device
			entry["price_list"] = item.Price_list
			entry["value"] = value
			entry["credit"] = value
			entry["tax_system"] = tran.TaxSystem.Type
			entry["tax_rate"] = item.Tax_rate
			entry["tax_value"] = tran.Stores.StoreAutomatEvent.GetNds(int(item.Tax_rate), value)
			entry["fn"] = tran.Fiscal.Fields.Fn
			entry["fd"] = tran.Fiscal.Fields.Fd
			entry["fp"] = tran.Fiscal.Fields.Fp
			entry["id_fr"] = tran.Fiscal.ResiptId
			entry["status"] = tran.Fiscal.Status
			if tran.Point.Id != 0 {
				entry["point_id"] = tran.Point.Id
			}
			entry["loyality_type"] = nil
			entry["loyality_code"] = nil
			entry["error_detail"] = tran.Fiscal.Message
			entry["warehouse_id"] = nil
			entry["type_fr"] = tran.Fiscal.Config.Type
			_, err := tran.Stores.StoreAutomatEvent.AddByParams(entry)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}
}
