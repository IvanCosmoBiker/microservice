package server

import (
	"context"
	"ephorfiscal/fiscal/factory"
	"ephorfiscal/fiscal/interface/fr"
	pb "ephorfiscal/service"
	"ephorservices/config"
	"fmt"
	"log"
)

type ServerFiscal struct {
	Config *config.Config
	pb.UnimplementedFiscalServiceServer
}

func Init(cfg *config.Config) {

}

func (sf *ServerFiscal) Fiscal(ctx context.Context, req *pb.Request) (*pb.Response, error) {

}

func (sf *ServerFiscal) Refund(ctx context.Context, req *pb.Request) (*pb.Response, error) {

}

func (sf *ServerFiscal) StatusKkt(ctx context.Context, req *pb.RequestStatus) (*pb.ResponseStatus, error) {

}

func (sf *ServerFiscal) DataVerification(req *pb.Request) {
	fiscalResult := make(map[string]interface{})
	log.Printf("%+v", tran)
	if tran.Fiscal.InCome.Imei != "" {
		registrator := factory.GetFiscal(tran.Fiscal.InCome.TypeFr, fm.Cfg)
		if registrator == nil {
			fiscalResult["status"] = false
			fiscalResult["message"] = "тип кассы не поддерживается"
			fiscalResult["fr_status"] = fr.Status_None
			fiscalResult["type_fr"] = tran.Fiscal.InCome.TypeFr
			return fiscalResult, nil
		}
		fiscalResult["status"] = true
		return fiscalResult, registrator
	}
	log.Printf("%+v", tran)
	if tran.Fiscal.Config.Id == int64(0) {
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = fr.Status_None
		fiscalResult["message"] = "нет активной кассы"
		return fiscalResult, nil
	}
	frModel, err := fm.GetFr(tran)
	log.Printf("%+v", frModel)
	if err != nil {
		log.Printf("%v", err)
		fiscalResult["status"] = false
		fiscalResult["fr_status"] = fr.Status_None
		fiscalResult["message"] = "касса не привязана к автомату"
		return fiscalResult, nil
	}
	log.Printf("%+v", frModel)
	fm.setFiscalConfig(tran, frModel)
	log.Printf("%+v", frModel)
	if tran.Payment.TypeSale == transaction.Sale_Cash && tran.Fiscal.Config.Fr_disable_cash != fr.Fr_Disable_Cash {
		fiscalResult["status"] = false
		fiscalResult["message"] = "отключение фискализации со стороны клиента"
		fiscalResult["fr_status"] = fr.Status_OFF_CA
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	if tran.Payment.TypeSale == transaction.Sale_Cashless && tran.Fiscal.Config.Fr_disable_cashless != fr.Fr_Disable_Cashless {
		fiscalResult["status"] = false
		fiscalResult["message"] = "отключение фискализации со стороны клиента"
		fiscalResult["fr_status"] = fr.Status_OFF_DA
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	if tran.Point.Id == 0 {
		fiscalResult["status"] = false
		fiscalResult["message"] = "нет торговой точки"
		fiscalResult["fr_status"] = fr.Status_None
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	log.Printf("%+v", frModel)
	checkMaxSum := fm.CheckMaxSumm(tran)
	if checkMaxSum == false {
		fiscalResult["status"] = false
		fiscalResult["message"] = "превышен лимит суммы по чеку"
		fiscalResult["fr_status"] = fr.Status_MAX_CHECK
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	registrator := factory.GetFiscal(tran.Fiscal.Config.Type, fm.Cfg)
	log.Printf("%+v", registrator)
	if registrator == nil {
		fiscalResult["status"] = false
		fiscalResult["message"] = "тип кассы не поддерживается"
		fiscalResult["fr_status"] = fr.Status_None
		fiscalResult["type_fr"] = frModel.Type
		return fiscalResult, nil
	}
	fiscalResult["status"] = true
	return fiscalResult, registrator
}

func (sf *ServerFiscal) CalcSumProducts(req *pb.Request) int {
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

func (sf *ServerFiscal) CheckMaxSumm(req *pb.Request) bool {
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
