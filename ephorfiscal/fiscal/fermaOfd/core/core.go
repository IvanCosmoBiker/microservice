package fermaOfd

import (
	"encoding/base64"
	"encoding/json"
	"ephorfiscal/fiscal/transaction"
	pb "ephorfiscal/service"
	config "ephorservices/config"
	request "ephorservices/ephorsale/fiscal/fermaOfd/request"
	randString "ephorservices/pkg/randgeneratestring"
	"fmt"
	"log"
	"math"
)

var layoutQr = "20060102T150405"

type Core struct {
	cfg *config.Config
}

func InitCore(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}

func (c *Core) MakeRequestAuth(login, password string) (jsonData []byte, err error) {
	resuestAuth := request.RequestAuthToken{}
	resuestAuth.Login = login
	resuestAuth.Password = password
	jsonData, err = json.Marshal(resuestAuth)
	if err != nil {
		return []byte(""), err
	}
	return jsonData, nil
}

func (c *Core) MakeRequestSendCheck(tran *pb.Request, fiscalModel *transaction.TransactionFiscal) (jsonData []byte, err error) {
	requestSendCheck := request.RequestSendCheck{}
	requestSendCheck.Request.Inn = fiscalModel.Inn
	requestSendCheck.Request.Inn = "Income"
	if len(tran.Fiscal.ResiptId) < 1 {
		var orderString randString.GenerateString
		orderString.RandGiud()
		requestSendCheck.Request.InvoiceId = orderString.String
	} else {
		requestSendCheck.Request.InvoiceId = fiscalModel.Id_shift
	}
	requestSendCheck.Request.CustomerReceipt.KktFA = true
	requestSendCheck.Request.CustomerReceipt.TaxationSystem = c.ConvertTaxationSystem(tran.TaxSystem)
	requestSendCheck.Request.CustomerReceipt.Email = fiscalModel.Email
	requestSendCheck.Request.CustomerReceipt.PaymentType = 1
	requestSendCheck.Request.CustomerReceipt.AutomaticDeviceNumber = int32(tran.Automat_number)
	requestSendCheck.Request.CustomerReceipt.BillAddress = tran.Point_addr
	Payment, Positions := c.GenerateDataForCheck(tran)
	requestSendCheck.Request.CustomerReceipt.Items = Positions
	requestSendCheck.Request.CustomerReceipt.PaymentItems = append(requestSendCheck.Request.CustomerReceipt.PaymentItems, Payment)
	jsonData, err = json.Marshal(requestSendCheck)
	if err != nil {
		return []byte(""), err
	}
	return jsonData, nil
}

func (C *Core) MakeRequestStatusCheck(tran *pb.Request, fiscalModel *transaction.TransactionFiscal) (jsonData []byte, err error) {
	requestStatus := request.RequestStatusCheck{}
	requestStatus.Request.ReceiptId = tran.Fiscal.ResiptId
	jsonData, err = json.Marshal(requestStatus)
	if err != nil {
		return []byte(""), err
	}
	return jsonData, nil
}

func (c *Core) GenerateDataForCheck(tran *pb.Request) (payment request.Payment, positions []request.Position) {
	var summ float64
	payment = request.Payment{}
	positions = make([]request.Position, 1)
	payment.PaymentType = 1
	payment.Sum = 0
	for _, product := range tran.Item {
		entryPosition := request.Position{}
		var value float64 = product.Price
		quantity := float64(product.Amount / 1000)
		price := float64(value / 100)
		summ += math.Round(quantity * price)

		entryPosition.PaymentMethod = 4
		entryPosition.PaymentType = 1
		entryPosition.Quantity = quantity
		entryPosition.Price = math.Round(price)
		entryPosition.Amount = math.Round(quantity * price)
		entryPosition.Vat = c.ConvertTax(int(product.Tax_rate))
		entryPosition.Label = product.Name
		positions = append(positions, entryPosition)
	}
	payment.Sum = summ
	tran.Fiscal.SumCheck = summ
	return payment, positions
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
