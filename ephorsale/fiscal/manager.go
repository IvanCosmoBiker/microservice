package fiscal

import (
	"context"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal/factory"
	"ephorservices/ephorsale/fiscal/interfaceFiscal"
	transaction "ephorservices/ephorsale/transaction"
	connectDb "ephorservices/pkg/db"
	fiscalStore "ephorservices/pkg/model/schema/account/fr/store"

	fiscalModel "ephorservices/pkg/model/schema/account/fr/model"
	"fmt"
)

type FiscalManager struct {
	StoreFiscal *fiscalStore.StoreFr
	cfg         *config.Config
	ctx         context.Context
}

func Init(conn *connectDb.Manager, conf *config.Config, contx context.Context) (*FiscalManager, error) {
	fr := fiscalStore.NewStore(conn)
	fiscal := &FiscalManager{
		StoreFiscal: fr,
		cfg:         conf,
		ctx:         contx,
	}
	return fiscal, nil
}

func (fm *FiscalManager) setFiscalConfig(tran *transaction.Transaction, frModel *fiscalModel.ReturningStruct) {
	tran.Fiscal.Config.Name = frModel.Name
	tran.Fiscal.Config.Type = uint8(frModel.Type)
	tran.Fiscal.Config.Dev_interface = int(frModel.Dev_interface)
	tran.Fiscal.Config.Login = frModel.Login
	tran.Fiscal.Config.Password = frModel.Password
	tran.Fiscal.Config.Phone = frModel.Phone
	tran.Fiscal.Config.Email = frModel.Email
	tran.Fiscal.Config.Dev_addr = frModel.Dev_addr
	tran.Fiscal.Config.Dev_port = int(frModel.Dev_port)
	tran.Fiscal.Config.Ofd_addr = frModel.Ofd_addr
	tran.Fiscal.Config.Ofd_port = int(frModel.Ofd_port)
	tran.Fiscal.Config.Inn = frModel.Inn
	tran.Fiscal.Config.Param1 = frModel.Param1
	tran.Fiscal.Config.Use_sn = int(frModel.Use_sn)
	tran.Fiscal.Config.Add_fiscal = int(frModel.Add_fiscal)
	tran.Fiscal.Config.Id_shift = frModel.Id_shift
	tran.Fiscal.Config.Fr_disable_cash = int(frModel.Fr_disable_cash)
	tran.Fiscal.Config.Fr_disable_cashless = int(frModel.Fr_disable_cashless)
	tran.Fiscal.Config.Ffd_version = int(frModel.Ffd_version)
	tran.Fiscal.Config.Auth_public_key = frModel.Auth_public_key
	tran.Fiscal.Config.Auth_private_key = frModel.Auth_private_key
	tran.Fiscal.Config.Sign_private_key = frModel.Sign_private_key
}

func (fm *FiscalManager) StartFiscal(tran *transaction.Transaction) (map[string]interface{}, interfaceFiscal.Fiscal) {
	fiscalResult := make(map[string]interface{})
	if tran.Fiscal.Config.Id == 0 {
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = interfaceFiscal.Status_None
		fiscalResult["message"] = "нет активной кассы"
		return fiscalResult, nil
	}
	frModel, err := fm.GetFr(tran)
	if err != nil {
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = interfaceFiscal.Status_None
		fiscalResult["message"] = "касса не привязана к автомату"
		return fiscalResult, nil
	}
	fm.setFiscalConfig(tran, frModel)
	if tran.Payment.TypeSale == transaction.Sale_Cash && tran.Fiscal.Config.Fr_disable_cash != interfaceFiscal.Fr_Disable_Cash {
		fiscalResult["status"] = false
		fiscalResult["message"] = "отключение фискализации со стороны клиента"
		fiscalResult["fr_status"] = interfaceFiscal.Status_OFF_CA
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	if tran.Payment.TypeSale == transaction.Sale_Cashless && tran.Fiscal.Config.Fr_disable_cashless != interfaceFiscal.Fr_Disable_Cashless {
		fiscalResult["status"] = false
		fiscalResult["message"] = "отключение фискализации со стороны клиента"
		fiscalResult["fr_status"] = interfaceFiscal.Status_OFF_DA
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	if tran.Point.Id == 0 {
		fiscalResult["status"] = false
		fiscalResult["message"] = "нет торговой точки"
		fiscalResult["fr_status"] = interfaceFiscal.Status_None
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	checkMaxSum := fm.CheckMaxSumm(tran)
	if checkMaxSum == false {
		fiscalResult["status"] = false
		fiscalResult["message"] = "превышен лимит суммы по чеку"
		fiscalResult["fr_status"] = interfaceFiscal.Status_MAX_CHECK
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	fr := factory.GetFiscal(tran.Fiscal.Config.Type)
	if fr == nil {
		fiscalResult["status"] = false
		fiscalResult["message"] = "тип кассы не поддерживается"
		fiscalResult["fr_status"] = interfaceFiscal.Status_None
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	fiscalResult["status"] = true
	return fiscalResult, fr
}

func (fm *FiscalManager) TimeOut() {

}

func (fm *FiscalManager) Send() {

}

func (fm *FiscalManager) GetFr(tran *transaction.Transaction) (*fiscalModel.ReturningStruct, error) {
	return fm.StoreFiscal.GetOneById(tran.Fiscal.Config.Id, tran.Config.AccountId)
}

func (fm *FiscalManager) CalcSumProducts(tran *transaction.Transaction) int {
	sum := 0
	if len(tran.Products) < 1 {
		sum = tran.Payment.Sum
		return sum
	}
	for _, product := range tran.Products {
		fmt.Println(int(product["price"].(float64) * product["quantity"].(float64)))
		sum += int(product["price"].(float64) * product["quantity"].(float64))
	}
	return sum
}

func (fm *FiscalManager) CheckMaxSumm(tran *transaction.Transaction) bool {
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
