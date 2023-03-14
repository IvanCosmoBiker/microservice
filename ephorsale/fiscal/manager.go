package fiscal

import (
	"context"
	"encoding/json"
	ofdFerma "ephorservices/ephorsale/fiscal/fermaOfd"
	"ephorservices/ephorsale/fiscal/interface/fr"
	ofdNanokass "ephorservices/ephorsale/fiscal/nanokass"
	ofdOrange "ephorservices/ephorsale/fiscal/orange"
	transportHttp "ephorservices/ephorsale/fiscal/transport"
	response "ephorservices/ephorsale/fiscal/transport/response"
	transaction_struct "ephorservices/ephorsale/transaction/transaction_struct"
	logger "ephorservices/pkg/logger"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"errors"
	"fmt"
	"runtime"
)

var (
	Error_Not_Automat  = errors.New("Not avalable automat")
	Error_Not_Location = errors.New("Automat is not installed at the point")
	Error_Not_Point    = errors.New("Not avalable point")
)

// instance of type fiscal
var ofd ofdFerma.NewFermaOfdStruct
var Orange ofdOrange.NewOrangeStruct
var Nanokass ofdNanokass.NewNanokassaStruct

var TypeFr = [7]uint8{fr.Fr_EphorOrangeData,
	fr.Fr_ServerOrangeData,
	fr.Fr_EphorServerOrangeData,
	fr.Fr_OrangeData,
	fr.Fr_NanoKassa,
	fr.Fr_ServerNanoKassa,
	fr.Fr_OFD}

type FiscalManager struct {
	UrlToSendHttp  string
	TransportHttp  *transportHttp.Http
	ctx            context.Context
	ExecuteMinutes int
	SleepSecond    int
	PathCert       string
	Debug          bool
}

var Fiscal *FiscalManager

func New(contx context.Context, executeMinutes, sleepSecond int, pathCert, urlSendHttp string, debug bool) (*FiscalManager, error) {
	fiscal := &FiscalManager{
		UrlToSendHttp:  urlSendHttp,
		TransportHttp:  transportHttp.New(debug),
		ctx:            contx,
		ExecuteMinutes: executeMinutes,
		SleepSecond:    sleepSecond,
		PathCert:       pathCert,
		Debug:          debug,
	}
	Fiscal = fiscal
	return fiscal, nil
}

func (fm *FiscalManager) setFiscalConfig(tran *transaction_struct.Transaction) error {
	frModel, err := tran.Stores.StoreFr.GetOneById(int(tran.Fiscal.Config.Id))
	if err != nil {
		return err
	}
	fmt.Printf("FR MODEL !!!!!!!!++++________________%+v\n", *frModel)
	tran.Fiscal.Config.Name = frModel.Name.String
	tran.Fiscal.Config.Type = uint8(frModel.Type.Int32)
	tran.Fiscal.Config.Dev_interface = int(frModel.Dev_interface.Int32)
	tran.Fiscal.Config.Login = frModel.Login.String
	tran.Fiscal.Config.Password = frModel.Password.String
	tran.Fiscal.Config.Phone = frModel.Phone.String
	tran.Fiscal.Config.Email = frModel.Email.String
	tran.Fiscal.Config.Dev_addr = frModel.Dev_addr.String
	tran.Fiscal.Config.Dev_port = int(frModel.Dev_port.Int32)
	tran.Fiscal.Config.Ofd_addr = frModel.Ofd_addr.String
	tran.Fiscal.Config.Ofd_port = int(frModel.Ofd_port.Int32)
	tran.Fiscal.Config.Inn = frModel.Inn.String
	tran.Fiscal.Config.Param1 = frModel.Param1.String
	tran.Fiscal.Config.Use_sn = int(frModel.Use_sn.Int32)
	tran.Fiscal.Config.Add_fiscal = int(frModel.Add_fiscal.Int32)
	tran.Fiscal.Config.Id_shift = frModel.Id_shift.String
	tran.Fiscal.Config.Fr_disable_cash = int(frModel.Fr_disable_cash.Int32)
	tran.Fiscal.Config.Fr_disable_cashless = int(frModel.Fr_disable_cashless.Int32)
	tran.Fiscal.Config.Ffd_version = int(frModel.Ffd_version.Int32)
	tran.Fiscal.Config.Auth_public_key = frModel.Auth_public_key.String
	tran.Fiscal.Config.Auth_private_key = frModel.Auth_private_key.String
	tran.Fiscal.Config.Sign_private_key = frModel.Sign_private_key.String
	frModel = nil
	return nil
}

func (fm *FiscalManager) CheckAutomat(tran *transaction_struct.Transaction) error {
	automat, _ := tran.Stores.StoreAutomat.GetOneById(tran.Config.AutomatId)
	if automat.Id == 0 {
		automat = nil
		tran.Fiscal.Message = Error_Not_Automat.Error()
		return Error_Not_Automat
	}
	req := tran.NewRequest()
	req.AddFilterParam("automat_id", req.Operator.OperatorEqual, true, automat.Id)
	req.AddFilterParam("to_date", req.Operator.OperatorEqual, true)
	locationAutomat, err := tran.Stores.StoreAutomatLocation.GetOneBy(req)
	if err != nil {
		locationAutomat = nil
		err = nil
		tran.Fiscal.Message = Error_Not_Location.Error()
		return Error_Not_Location
	}
	if locationAutomat.Id == 0 {
		locationAutomat = nil
		tran.Fiscal.Message = Error_Not_Location.Error()
		return Error_Not_Location
	}
	tran.Fiscal.Config.Id = int64(automat.Fr_id.Int32)
	tran.Point.Id = parserTypes.ParseTypeInterfaceToInt(locationAutomat.Company_point_id.Int32)
	entryPoint, _ := tran.Stores.StoreCompanyPoint.GetOneById(int(locationAutomat.Company_point_id.Int32))
	if entryPoint.Id == 0 {
		entryPoint = nil
		tran.Fiscal.Message = Error_Not_Point.Error()
		return Error_Not_Point
	}
	tran.Point.Address = entryPoint.Address.String
	tran.Point.PointName = entryPoint.Name.String
	entryPoint = nil
	locationAutomat = nil
	automat = nil
	return nil
}

func (fm *FiscalManager) DataVerification(tran *transaction_struct.Transaction) (map[string]interface{}, fr.Fiscal) {
	fiscalResult := make(map[string]interface{})
	if tran.Fiscal.OnlyFiscal {
		registrator := fm.GetFiscal(uint8(tran.Fiscal.Config.Type))
		if registrator == nil {
			fiscalResult["status"] = false
			fiscalResult["message"] = "тип кассы не поддерживается"
			fiscalResult["fr_status"] = fr.Status_None
			fiscalResult["type_fr"] = tran.Fiscal.Config.Type
			logger.Log.Errorf("%+v", fiscalResult)
			return fiscalResult, nil
		}
		fiscalResult["status"] = true
		return fiscalResult, registrator
	}
	err := fm.CheckAutomat(tran)
	if err != nil {
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = fr.Status_None
		fiscalResult["message"] = err.Error()
		logger.Log.Errorf("%+v", fiscalResult)
		return fiscalResult, nil
	}
	checkMaxSum := fm.CheckMaxSumm(tran)
	if !checkMaxSum {
		fiscalResult["status"] = false
		fiscalResult["message"] = "превышен лимит суммы по чеку"
		fiscalResult["fr_status"] = fr.Status_MAX_CHECK
		fiscalResult["type_fr"] = tran.Fiscal.Config.Type
		logger.Log.Errorf("%+v", fiscalResult)
		return fiscalResult, nil
	}
	if tran.Fiscal.Config.Id == int64(0) {
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = fr.Status_None
		fiscalResult["message"] = "нет активной кассы"
		logger.Log.Errorf("%+v", fiscalResult)
		return fiscalResult, nil
	}
	err = fm.setFiscalConfig(tran)
	if err != nil {
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = fr.Status_Error
		fiscalResult["message"] = err.Error()
		logger.Log.Errorf("%+v", fiscalResult)
		return fiscalResult, nil
	}
	checkDisable := fm.CheckFiscalOfTypePayment(tran)
	if checkDisable == 0 {
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = fr.Status_OFF_FR
		fiscalResult["message"] = "отключение фискализации со стороны клиента"
		return fiscalResult, nil
	}
	registrator := fm.GetFiscal(tran.Fiscal.Config.Type)
	if registrator == nil {
		fiscalResult["status"] = false
		fiscalResult["message"] = "тип кассы не поддерживается"
		fiscalResult["fr_status"] = fr.Status_None
		fiscalResult["type_fr"] = tran.Fiscal.Config.Type
		logger.Log.Errorf("%+v", fiscalResult)
		return fiscalResult, nil
	}
	fiscalResult["status"] = true
	return fiscalResult, registrator
}

func (fm *FiscalManager) GetFiscal(typeFiscal uint8) fr.Fiscal {
	switch typeFiscal {
	case fr.Fr_EphorOrangeData,
		fr.Fr_ServerOrangeData,
		fr.Fr_EphorServerOrangeData,
		fr.Fr_OrangeData:
		return Orange.New(fm.ExecuteMinutes, fm.SleepSecond, fm.PathCert, fm.Debug)
	case fr.Fr_NanoKassa,
		fr.Fr_ServerNanoKassa:
		return Nanokass.New(fm.ExecuteMinutes, fm.SleepSecond, fm.Debug)
	case fr.Fr_OFD:
		return ofd.New(fm.ExecuteMinutes, fm.SleepSecond, fm.Debug)
	}
	return nil
}

func (fm *FiscalManager) MakeResponseHttp(tran *transaction_struct.Transaction) *response.OutCome {
	res := response.OutCome{}
	res.Imei = tran.Config.Imei
	res.Data.Message = tran.Fiscal.Message
	if tran.Status == int(fr.Status_Complete) {
		res.Data.Status = "success"
	} else {
		res.Data.Status = "unsuccess"
	}
	res.Data.Code = tran.Fiscal.Code
	res.Data.StatusCode = tran.Fiscal.StatusCode
	res.SetEventId(tran.Fiscal.Events)
	res.Data.Fields.Fp = tran.Fiscal.Fields.Fp
	res.Data.Fields.Fd = parserTypes.ParseTypeInString(tran.Fiscal.Fields.Fd)
	res.Data.Fields.Fn = tran.Fiscal.Fields.Fn
	res.Data.Fields.DateFisal = tran.Fiscal.Fields.DateFisal
	fmt.Printf("\n !!!!!!!!-----------------RESPONSE SET: %+v\n", res)
	return &res
}

func (fm *FiscalManager) Send(protocol string, tran *transaction_struct.Transaction) {
	fmt.Printf("%+v", *tran)
	if protocol == "http" {
		res := fm.MakeResponseHttp(tran)
		jsonRequest, _ := json.Marshal(res.Data)
		headers := make(map[string]interface{})
		headers["Content-Length"] = len(jsonRequest)
		urlSend := fmt.Sprintf("%s&login=%s&password=%s&_dc=%s", fm.UrlToSendHttp, tran.Config.Imei, "12345678", tran.DateTime.UnixNano())
		funcSend := fm.TransportHttp.Send(fm.TransportHttp.Call, "POST", urlSend, headers, jsonRequest, 900000, 60)
		code, errResp := funcSend("POST", urlSend, headers, jsonRequest)
		if fm.Debug {
			logger.Log.Infof("%v", string(jsonRequest))
			logger.Log.Infof("%v", errResp)
			logger.Log.Infof("%v", code)
		}
		transaction_struct.Destroy(tran)
		res = nil
	}
	runtime.Goexit()
}

func (fm *FiscalManager) CalcSumProducts(tran *transaction_struct.Transaction) int {
	sum := 0
	if len(tran.Products) < 1 {
		sum = tran.Payment.Sum
		return sum
	}
	for _, product := range tran.Products {
		var value float64 = 0
		if product.Value == value {
			value = product.Price
		} else {
			value = product.Value
		}
		fmt.Println(int(value * float64(product.Quantity)))
		sum += int(value * float64(product.Quantity))
	}
	return sum
}

func (fm *FiscalManager) CheckMaxSumm(tran *transaction_struct.Transaction) bool {
	maxSum := fm.CalcSumProducts(tran)
	fmt.Println(maxSum)
	if tran.Fiscal.Config.MaxSum == 0 {
		return true
	}
	if maxSum > tran.Fiscal.Config.MaxSum {
		return false
	}
	return true
}

func (fm *FiscalManager) CheckFiscalOfTypePayment(tran *transaction_struct.Transaction) int {
	result := 0
	for _, product := range tran.Products {
		fmt.Printf("FR_DISABLE : %v", tran.Fiscal.Config.Fr_disable_cashless)
		fmt.Printf("FR_DISABLE CONST : %v", fr.Fr_Disable_Cashless)
		if product.Payment_device == "DA" && tran.Fiscal.Config.Fr_disable_cashless == fr.Fr_Disable_Cashless {
			result++
			product.Fiscalization = true
		}
		if product.Payment_device == "CA" && tran.Fiscal.Config.Fr_disable_cash == fr.Fr_Disable_Cash {
			result++
			product.Fiscalization = true
		}
	}

	return result
}
