package transaction

import (
	"context"
	config "ephorservices/config"

	"ephorservices/ephorsale/transaction/transaction_struct"
	storeLog "ephorservices/internal/model/schema/main/log/store"
	storeModem "ephorservices/internal/model/schema/main/modem/store"
	storeTransaction "ephorservices/internal/model/schema/main/transaction/store"
	storeTransactionProduct "ephorservices/internal/model/schema/main/transactionproduct/store"
	dateTime "ephorservices/pkg/datetime"
	parserTypes "ephorservices/pkg/parser/typeParse"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type TransactionDispetcher struct {
	mutex                   sync.Mutex
	rmutex                  sync.RWMutex
	Ctx                     context.Context
	Config                  *config.Config
	transactions            map[int]*transaction_struct.Transaction
	replayProtection        map[int]interface{}
	Date                    *dateTime.DateTime
	StoreLog                *storeLog.StoreLog
	StoreTransaction        *storeTransaction.StoreTransaction
	StoreTransactionProduct *storeTransactionProduct.StoreTransactionProduct
	StoreModem              *storeModem.StoreModem
}

var Dispetcher *TransactionDispetcher

func New(conf *config.Config, ctx context.Context) *TransactionDispetcher {
	var newMutex sync.Mutex
	var newRmutex sync.RWMutex
	var newTransactions = make(map[int]*transaction_struct.Transaction)
	var newReplayProtection = make(map[int]interface{})
	date, _ := dateTime.Init()
	TransactionDispetcherNew := &TransactionDispetcher{
		mutex:            newMutex,
		rmutex:           newRmutex,
		transactions:     newTransactions,
		replayProtection: newReplayProtection,
		Date:             date,
		Ctx:              ctx,
		Config:           conf,
	}
	TransactionDispetcherNew.initStore()
	Dispetcher = TransactionDispetcherNew
	return TransactionDispetcherNew
}

func (t *TransactionDispetcher) initStore() {
	t.StoreLog = storeLog.New()
	t.StoreTransaction = storeTransaction.New()
	t.StoreTransactionProduct = storeTransactionProduct.New()
	t.StoreModem = storeModem.New()
}

func (t *TransactionDispetcher) NewTransaction() *transaction_struct.Transaction {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	transactionStruct := transaction_struct.InitTransaction()
	return transactionStruct
}

func (t *TransactionDispetcher) StartTransaction(tran *transaction_struct.Transaction) error {
	err := t.AddTransaction(tran)
	return err
}

func (t *TransactionDispetcher) AddTransaction(tran *transaction_struct.Transaction) error {
	rand.Seed(time.Now().UnixNano())
	randNoise := t.Random(10000000, 20000000)
	parametrs := make(map[string]interface{})
	parametrs["automat_id"] = tran.Config.AutomatId
	parametrs["account_id"] = tran.Config.AccountId
	parametrs["token_id"] = tran.Payment.Token
	parametrs["status"] = transaction_struct.TransactionState_MoneyHoldStart
	parametrs["date"] = t.Date.Now()
	parametrs["pay_type"] = tran.Payment.PayType
	parametrs["sum"] = tran.Payment.Sum
	parametrs["ps_type"] = tran.Payment.Type
	parametrs["token_type"] = tran.Payment.TokenType
	parametrs["qr_format"] = tran.Fiscal.Config.QrFormat
	model, err := t.StoreTransaction.AddByParams(parametrs)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	noise := randNoise + model.Id
	updateNoise := make(map[string]interface{})
	updateNoise["id"] = model.Id
	updateNoise["noise"] = strconv.Itoa(noise)
	log.Printf("%+v", updateNoise)
	_, err = t.StoreTransaction.SetByParams(updateNoise)
	tran.Config.Tid = model.Id
	tran.Config.Noise = noise
	t.AddChannel(tran.Config.Tid, tran)
	t.AddTransactionProduct(tran)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func (t *TransactionDispetcher) Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func (t *TransactionDispetcher) AddTransactionProduct(tran *transaction_struct.Transaction) {
	params := make(map[string]interface{})
	for _, product := range tran.Products {
		var value float64 = 0
		if product.Value == value {
			value = product.Price
		} else {
			value = product.Value
		}
		tran.Payment.DebitSum += parserTypes.ParseTypeInterfaceToInt(value * float64(product.Quantity/1000))
		params["transaction_id"] = tran.Config.Tid
		params["name"] = product.Name
		params["select_id"] = product.Select_id
		params["ware_id"] = product.Ware_id
		params["value"] = value
		params["tax_rate"] = product.Tax_rate
		params["quantity"] = product.Quantity / 1000
		transactionProduct, err := t.StoreTransactionProduct.AddByParams(params)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		product.Payment_device = "DA"
		tran.Sum += transactionProduct.Value.Int32 * transactionProduct.Quantity.Int32
	}
}

func (t *TransactionDispetcher) RemoveTransaction(key int) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	tran, exist := t.transactions[key]
	if !exist {
		return false
	}
	close(tran.ChannelMessage)
	close(tran.TimeOut)
	close(tran.Close)
	delete(t.transactions, key)
	return true
}

func (t *TransactionDispetcher) CheckDuplicate(automat, account int) bool {
	key := automat + account
	return t.GetReplayProtection(key)
}

/*
key is composed by accountId and automatId
*/
func (t *TransactionDispetcher) AddReplayProtection(key, automat int) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.replayProtection[key] = automat
	return true
}

func (t *TransactionDispetcher) GetReplayProtection(key int) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_, exist := t.replayProtection[key]
	return exist
}

func (t *TransactionDispetcher) RemoveReplayProtection(key int) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_, exist := t.replayProtection[key]
	if !exist {
		return false
	}
	delete(t.replayProtection, key)
	return true
}

func (t *TransactionDispetcher) AddChannel(key int, tran *transaction_struct.Transaction) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.transactions[key] = tran
}

func (t *TransactionDispetcher) Send(key int, message []byte) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	tran, exist := t.transactions[key]
	if !exist {
		return false
	}
	tranSend := *tran
	tranSend.ChannelMessage <- message
	return true
}

func (t *TransactionDispetcher) GetTransactions() map[int]*transaction_struct.Transaction {
	t.rmutex.RLock()
	defer t.rmutex.RUnlock()
	return t.transactions
}

func (t *TransactionDispetcher) GetOneTransaction(key int) (*transaction_struct.Transaction, bool) {
	t.rmutex.RLock()
	defer t.rmutex.RUnlock()
	tran, exist := t.transactions[key]
	return tran, exist
}
