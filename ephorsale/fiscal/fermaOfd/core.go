package fermaOfd

import (
	"encoding/base64"
	"encoding/json"
	request "ephorservices/ephorsale/fiscal/fermaOfd/request"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	logger "ephorservices/pkg/logger"
	randString "ephorservices/pkg/randgeneratestring"
	"fmt"
	"log"
	"math"
)

var layoutQr = "20060102T150405"

type Core struct{}

func InitCore() *Core {
	return &Core{}
}

func (c *Core) MakeRequestAuth(tran *transaction.Transaction) (jsonData []byte, err error) {
	resuestAuth := request.RequestAuthToken{}
	resuestAuth.Login = tran.Fiscal.Config.Login
	resuestAuth.Password = tran.Fiscal.Config.Password
	jsonData, err = json.Marshal(resuestAuth)
	if err != nil {
		return []byte(""), err
	}
	return jsonData, nil
}

func (c *Core) MakeRequestSendCheck(tran *transaction.Transaction) ([]byte, error) {
	if tran.Fiscal.OnlyFiscal {
		return tran.Fiscal.FiscalRequest.Request, nil
	}
	requestSendCheck := &request.RequestSendCheck{}
	requestSendCheck.Request.Inn = tran.Fiscal.Config.Inn
	requestSendCheck.Request.Type = "Income"
	if len(tran.Fiscal.ResiptId) < 1 {
		orderString := randString.Init()
		orderString.Strimwidth("23", 0, 64, fmt.Sprintf("%v%v%v", tran.Config.AccountId, tran.Config.AutomatId, tran.DateTime.StringToUnix(tran.Date)))
		requestSendCheck.Request.InvoiceId = orderString.String
		orderString = nil
	} else {
		requestSendCheck.Request.InvoiceId = tran.Fiscal.ResiptId
	}
	requestSendCheck.Request.CustomerReceipt.KktFA = true
	requestSendCheck.Request.CustomerReceipt.TaxationSystem = c.ConvertTaxationSystem(tran.TaxSystem.Type)
	requestSendCheck.Request.CustomerReceipt.Email = tran.Fiscal.Config.Email
	requestSendCheck.Request.CustomerReceipt.PaymentType = 1
	requestSendCheck.Request.CustomerReceipt.AutomaticDeviceNumber = int32(tran.Config.AutomatId)
	requestSendCheck.Request.CustomerReceipt.BillAddress = tran.Point.Address
	Payments, Positions := c.GenerateDataForCheck(tran)
	requestSendCheck.Request.CustomerReceipt.Items = Positions
	requestSendCheck.Request.CustomerReceipt.PaymentItems = append(requestSendCheck.Request.CustomerReceipt.PaymentItems, Payments...)
	jsonData, err := json.Marshal(requestSendCheck)
	if err != nil {
		logger.Log.Error(err.Error())
		jsonData = nil
		requestSendCheck = nil
		return []byte(""), err
	}
	requestSendCheck = nil
	return jsonData, nil
}

func (C *Core) MakeRequestStatusCheck(tran *transaction.Transaction) (jsonData []byte, err error) {
	requestStatus := &request.RequestStatusCheck{}
	requestStatus.Request.ReceiptId = tran.Fiscal.ResiptId
	jsonData, err = json.Marshal(requestStatus)
	if err != nil {
		requestStatus = nil
		jsonData = nil
		return []byte(""), err
	}
	requestStatus = nil
	return jsonData, nil
}

func (c *Core) GenerateDataForCheck(tran *transaction.Transaction) (payments []*request.Payment, positions []*request.Position) {
	var summ float64
	paymentCA := &request.Payment{}
	paymentCA.PaymentType = 0
	paymentCA.Sum = 0
	paymentDA := &request.Payment{}
	paymentDA.PaymentType = 1
	paymentDA.Sum = 0
	payments = make([]*request.Payment, 0, 1)
	positions = make([]*request.Position, 0, len(tran.Products))

	for _, product := range tran.Products {
		if !product.Fiscalization {
			continue
		}
		entryPosition := &request.Position{}

		var value float64 = 0
		if product.Value == value {
			value = product.Price
		} else {
			value = product.Value
		}
		if product.Marking != 0 {
			if len(product.Barcode) != 0 {
				entryPosition.MarkingCodeData = &request.MarkingCodeData{
					Type:          "UNKNOWN_PRODUCT_CODE",
					Code:          fmt.Sprintf("0%s", product.Barcode),
					PlannedStatus: "PRODUCT_STATUS_NOT_CHANGED",
				}
			}
		}
		quantity := float64(product.Quantity / 1000)
		price := float64(value / 100)
		if product.Payment_device == "CA" {
			paymentCA.Sum += math.Round(quantity * price)
		} else if product.Payment_device == "DA" {
			paymentDA.Sum += math.Round(quantity * price)
		}

		entryPosition.PaymentMethod = 4
		entryPosition.PaymentType = 1
		entryPosition.Quantity = quantity
		entryPosition.Price = math.Round(price)
		entryPosition.Amount = math.Round(quantity * price)
		entryPosition.Vat = c.ConvertTax(int(product.Tax_rate))
		entryPosition.Label = product.Name
		positions = append(positions, entryPosition)

	}
	if paymentCA.Sum != float64(0) {
		payments = append(payments, paymentCA)
	}
	if paymentDA.Sum != float64(0) {
		payments = append(payments, paymentDA)
	}
	tran.Fiscal.SumCheck = summ
	return payments, positions
}

func (c *Core) ConvertPaymentDevice(paymentDevice string) int {
	switch paymentDevice {
	case "CA":
		return 0
		fallthrough
	case "DA":
		return 1
	}
	return 0
}

func (c *Core) ConvertTax(tax int) string {
	switch tax {
	case transaction.TaxRate_NDSNone:
		return "VatNo"
		fallthrough
	case transaction.TaxRate_NDS0:
		return "Vat0"
		fallthrough
	case transaction.TaxRate_NDS10:
		return "Vat10"
		fallthrough
	case transaction.TaxRate_NDS18:
		return "Vat20"
	}
	return "VatNo"
}

func (c *Core) ConvertTaxationSystem(taxsystem int) string {
	switch taxsystem {
	case transaction.TaxSystem_OSN:
		return "Common"
		fallthrough
	case transaction.TaxSystem_USND:
		return "SimpleIn"
		fallthrough
	case transaction.TaxSystem_USNDMR:
		return "SimpleInOut"
		fallthrough
	case transaction.TaxSystem_ENVD:
		return "Unified"
		fallthrough
	case transaction.TaxSystem_ESN:
		return "UnifiedAgricultural"
		fallthrough
	case transaction.TaxSystem_Patent:
		return "Patent"
	}
	return "Common"
}

func (c *Core) EncodeUrlToBase64(tran *transaction.Transaction) string {
	stringUrl := c.MakeUrlQr(tran)
	return base64.StdEncoding.EncodeToString([]byte(stringUrl))
}

func (c *Core) MakeUrlQr(tran *transaction.Transaction) string {
	stringResult := fmt.Sprintf("t=%s&s=%v&fn=%v&i=%v&fp=%v&n=1", tran.Fiscal.Fields.DateFisal, fmt.Sprintf("%.2f", tran.Fiscal.SumCheck), tran.Fiscal.Fields.Fn, tran.Fiscal.Fields.Fd, tran.Fiscal.Fields.Fp)
	log.Println(stringResult)
	return stringResult
}
