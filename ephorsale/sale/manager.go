package sale

import (
	"context"
	"encoding/json"
	config "ephorservices/config"
	"ephorservices/ephorsale/fiscal"
	"ephorservices/ephorsale/payment"
	"ephorservices/ephorsale/sale/factory"
	"ephorservices/ephorsale/sale/interfaceSale"
	requestPay "ephorservices/ephorsale/sale/request"
	responseQueueManager "ephorservices/ephorsale/sale/responseQueueManager"
	transaction "ephorservices/ephorsale/transaction"
	rabbit "ephorservices/pkg/rabbitmq"
	transportHttp "ephorservices/pkg/transportprotocol/http/v1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type PayConsumer struct {
	Dispetcher *transaction.TransactionDispetcher
}

func (pc *PayConsumer) Consume(ctx context.Context, data []byte) error {
	req := responseQueueManager.ResponseQueue{}
	err := json.Unmarshal(data, &req)
	log.Printf("%+v", req)
	if err != nil {
		return errors.New("Fail parse json")
	}
	result := pc.Dispetcher.Send(req.Tid, data)
	if !result {
		return errors.New("Fail send")
	}
	return nil
}

func (pc *PayConsumer) Shutdown(ctx context.Context) error {
	log.Println("end pay")
	return errors.New("end pay")
}

type SaleManager struct {
	Status     int
	Rabbit     *rabbit.Manager
	Dispetcher *transaction.TransactionDispetcher
	Fiscal     *fiscal.FiscalManager
	Payment    *payment.PaymentManager
	Pay        *PayConsumer
	cfg        *config.Config
	ctx        context.Context
}

var saleManager *SaleManager

func Init(ctx context.Context, conf *config.Config, RabbitMq *rabbit.Manager, errGroup *errgroup.Group, dispether *transaction.TransactionDispetcher, fiscalM *fiscal.FiscalManager, paymentM *payment.PaymentManager, Transport *transportHttp.ServerHttp) (*SaleManager, error) {
	pay, err := initConsumer(ctx, conf, errGroup, dispether, RabbitMq)
	if err != nil {
		return &SaleManager{}, err
	}
	sale := &SaleManager{
		cfg:        conf,
		Dispetcher: dispether,
		Rabbit:     RabbitMq,
		Pay:        pay,
		Fiscal:     fiscalM,
		Payment:    paymentM,
		ctx:        ctx,
	}
	initHttpUrl(Transport)
	saleManager = sale
	return sale, nil
}

func initConsumer(ctx context.Context, conf *config.Config, errorGroup *errgroup.Group, dispether *transaction.TransactionDispetcher, RabbitMq *rabbit.Manager) (*PayConsumer, error) {
	consumer, err := RabbitMq.AddConsumer(ctx, conf.Services.EphorPay.NameQueue, "Pay", "", conf.Services.EphorPay.NameQueue)
	if err != nil {
		return &PayConsumer{}, err
	}
	pay := &PayConsumer{}
	err = consumer.Subscribe(ctx, errorGroup, pay)
	if err != nil {
		return pay, err
	}
	pay.Dispetcher = dispether
	log.Println("PayStart")
	return pay, nil
}

func handlerFiscal(w http.ResponseWriter, req *http.Request) {
	return
}

func handlerPay(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "POST":
		var byteData []byte
		var errResp error
		response := make(map[string]interface{})
		json_data, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		request := requestPay.RequestPay{}
		err := request.Init(json_data)
		if err != nil {
			response["message"] = fmt.Sprintf("%v", err)
			byteData, errResp = MakeResponse(response)
			if errResp != nil {
				log.Printf("%+v", response)
				w.Write([]byte(errResp.Error()))
				return
			}
			log.Printf("%+v", response)
			w.Write(byteData)
			return
		}
		if request.Sum == 0 && len(request.Products) < 1 {
			response["message"] = "Сумма 0 или массив продуктов пустой"
			byteData, errResp = MakeResponse(response)
			if errResp != nil {
				log.Printf("%+v", response)
				w.Write([]byte(errResp.Error()))
				return
			}
			log.Printf("%+v", response)
			w.Write(byteData)
			return
		}
		check := saleManager.Dispetcher.CheckDuplicate(request.Config.AutomatId, request.Config.AccountId)
		if check != false {
			response["message"] = "Действие над автоматом производится другим пользователем, пожалуйста подождите"
			byteData, errResp = MakeResponse(response)
			if errResp != nil {
				log.Printf("%+v", response)
				w.Write([]byte(errResp.Error()))
			}
			log.Printf("%+v", response)
			w.Write(byteData)
			return
		}
		SalePay := saleManager.InitSale(uint8(request.Config.DeviceType), uint8(request.Config.PayType))
		log.Printf("%+v", SalePay)
		if SalePay == nil {
			response["message"] = "Данный способ оплаты не реализован"
			byteData, errResp = MakeResponse(response)
			if errResp != nil {
				w.Write([]byte(errResp.Error()))
				log.Printf("%+v", response)
				return
			}
			log.Printf("%+v", response)
			w.Write(byteData)
			return
		}
		tran := saleManager.InitTransactionPay(&request)
		err = saleManager.Dispetcher.StartTransaction(tran)
		if err != nil {
			response["message"] = fmt.Sprintf("%v", err)
			byteData, errResp = MakeResponse(response)
			if errResp != nil {
				log.Printf("%+v", response)
				w.Write([]byte(errResp.Error()))
				return
			}
			log.Printf("%+v", response)
			w.Write(byteData)
			return
		}
		go SalePay.Start(tran)
		response["message"] = "ok"
		response["tid"] = tran.Config.Noise
		byteData, errResp = MakeResponse(response)
		if errResp != nil {
			log.Printf("%+v", response)
			w.Write([]byte(errResp.Error()))
		}
		log.Printf("%+v", response)
		w.Write(byteData)
		return
	case "GET":
		fmt.Fprintf(w, "%s: Running\n", "Servirce Payment")
		log.Println("Running")
	}
}

func MakeResponse(response map[string]interface{}) ([]byte, error) {
	return json.Marshal(response)
}

func initHttpUrl(transport *transportHttp.ServerHttp) {
	routePay := transport.SetHandlerListener("/pay", handlerPay)
	routeFiscal := transport.SetHandlerListener("/fiscal", handlerFiscal)
	routeFiscal.Methods("GET", "POST")
	routePay.Methods("GET", "POST")
}

func (sm *SaleManager) InitSale(DeviceType, PayType uint8) interfaceSale.Sale {
	return factory.GetSale(DeviceType, PayType, sm.Rabbit, sm.cfg, sm.Fiscal, sm.Payment, sm.Dispetcher)
}

func (sm *SaleManager) InitTransactionPay(request *requestPay.RequestPay) *transaction.Transaction {
	tran := sm.Dispetcher.NewTransaction()
	tran.Config.AccountId = request.Config.AccountId
	tran.Config.AutomatId = request.Config.AutomatId
	tran.Config.TokenType = request.Config.TokenType
	tran.Config.DeviceType = request.Config.DeviceType
	tran.Config.CurrensyCode = request.Config.CurrensyCode
	tran.Config.Imei = request.Imei
	tran.Payment.PayType = request.Config.PayType
	tran.Payment.TokenType = request.Config.TokenType
	tran.Payment.Token = request.PaymentToken
	tran.Payment.UserPhone = request.Config.UserPhone
	tran.Payment.ReturnUrl = request.Config.ReturnUrl
	tran.Payment.DeepLink = request.Config.DeepLink
	tran.Payment.Login = request.Config.Login
	tran.Payment.Password = request.Config.Password
	tran.Payment.Type = uint8(request.Config.BankType)
	tran.Payment.CurrensyCode = request.Config.CurrensyCode
	tran.Payment.Language = request.Config.Language
	tran.Payment.GateWay = request.GateWay
	tran.Payment.MerchantId = request.MerchantId
	tran.Payment.SbolBankInvoiceId = request.SbolBankInvoiceId
	tran.Payment.HostPaymentSystem = request.HostPaymentSystem
	tran.Payment.SbpPoint = request.Config.SbpPoint
	tran.Payment.Service_id = request.Service_id
	tran.Payment.SecretKey = request.SecretKey
	tran.Payment.KeyPayment = request.KeyPayment
	tran.Payment.TidPaymentSystem = request.TidPaymentSystem
	tran.Payment.Sum = request.Sum
	tran.Fiscal.Config.QrFormat = request.Config.QrFormat
	tran.Products = request.Products
	return tran
}
