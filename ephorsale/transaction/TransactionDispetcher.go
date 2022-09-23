package transaction

import (
	dateTime "ephorservices/pkg/datetime"
	db "ephorservices/pkg/db"
	storeAutomatConfig "ephorservices/pkg/model/schema/account/automat/config/store"
	storeAutomat "ephorservices/pkg/model/schema/account/automat/store"
	storeLocation "ephorservices/pkg/model/schema/account/automatlocation/store"
	storePoint "ephorservices/pkg/model/schema/account/companypoint/store"
	storeConfigProduct "ephorservices/pkg/model/schema/account/config/product/store"
	storeConfig "ephorservices/pkg/model/schema/account/config/store"
	storeWare "ephorservices/pkg/model/schema/account/ware/store"
	storeModem "ephorservices/pkg/model/schema/main/modem/store"
	storeTransaction "ephorservices/pkg/model/schema/main/transaction/store"
	storeTransactionProduct "ephorservices/pkg/model/schema/main/transactionproduct/store"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type TransactionDispetcher struct {
	mutex                   sync.Mutex
	rmutex                  sync.RWMutex
	transactions            map[int]*Transaction
	replayProtection        map[int]interface{}
	Date                    *dateTime.DateTime
	StoreTransaction        *storeTransaction.StoreTransaction
	StoreTransactionProduct *storeTransactionProduct.StoreTransactionProduct
	StoreAutomat            *storeAutomat.StoreAutomat
	StoreModem              *storeModem.StoreModem
	StoreLocation           *storeLocation.StoreAutomatLocation
	StorePoint              *storePoint.StoreCompanyPoint
	StoreWare               *storeWare.StoreWare
	StoreAutomatConfig      *storeAutomatConfig.StoreAutomatConfig
	StoreConfig             *storeConfig.StoreConfig
	StoreConfigProduct      *storeConfigProduct.StoreConfigProduct
}

func New(conn *db.Manager) *TransactionDispetcher {
	var newMutex sync.Mutex
	var newRmutex sync.RWMutex
	var newTransactions = make(map[int]*Transaction)
	var newReplayProtection = make(map[int]interface{})
	date, _ := dateTime.Init()
	TransactionDispetcherNew := &TransactionDispetcher{
		mutex:            newMutex,
		rmutex:           newRmutex,
		transactions:     newTransactions,
		replayProtection: newReplayProtection,
		Date:             date,
	}
	initStore(conn, TransactionDispetcherNew)
	return TransactionDispetcherNew
}

func initStore(conn *db.Manager, tran *TransactionDispetcher) {
	tran.StoreTransaction = storeTransaction.NewStore(conn)
	tran.StoreTransactionProduct = storeTransactionProduct.NewStore(conn)
	tran.StoreAutomat = storeAutomat.NewStore(conn)
	tran.StoreModem = storeModem.NewStore(conn)
	tran.StoreLocation = storeLocation.NewStore(conn)
	tran.StorePoint = storePoint.NewStore(conn)
	tran.StoreWare = storeWare.NewStore(conn)
	tran.StoreAutomatConfig = storeAutomatConfig.NewStore(conn)
	tran.StoreConfig = storeConfig.NewStore(conn)
	tran.StoreConfigProduct = storeConfigProduct.NewStore(conn)
}

func (t *TransactionDispetcher) NewTransaction() *Transaction {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	transactionStruct := InitTransaction()
	return transactionStruct
}

func (t *TransactionDispetcher) StartTransaction(tran *Transaction) error {
	err := t.AddTransaction(tran)
	return err
}

func (t *TransactionDispetcher) AddTransaction(tran *Transaction) error {
	rand.Seed(time.Now().UnixNano())
	randNoise := Random(10000000, 20000000)
	parametrs := make(map[string]interface{})
	parametrs["automat_id"] = tran.Config.AutomatId
	parametrs["account_id"] = tran.Config.AccountId
	parametrs["token_id"] = tran.Payment.Token
	parametrs["status"] = TransactionState_MoneyHoldStart
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
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	tran.Config.Tid = model.Id
	tran.Config.Noise = noise
	log.Printf("%+v", tran)
	t.AddChannel(tran.Config.Tid, tran)
	t.AddTransactionProduct(tran)
	return nil
}

func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func (t *TransactionDispetcher) AddTransactionProduct(tran *Transaction) {
	params := make(map[string]interface{})
	for _, product := range tran.Products {
		params["transaction_id"] = tran.Config.Tid
		params["name"] = product["name"]
		params["select_id"] = product["select_id"]
		params["ware_id"] = product["ware_id"]
		params["value"] = product["price"]
		params["tax_rate"] = product["tax_rate"]
		params["quantity"] = product["quantity"]
		t.StoreTransactionProduct.AddByParams(params)
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

func (t *TransactionDispetcher) AddChannel(key int, tran *Transaction) {
	log.Println("lock")
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.transactions[key] = tran
}

func (t *TransactionDispetcher) Send(key int, message []byte) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	tran, exist := t.transactions[key]
	log.Printf("\n [x] %+v", tran)
	if !exist {
		return false
	}
	tranSend := *tran
	tranSend.ChannelMessage <- message
	return true
}

func (t *TransactionDispetcher) GetTransactions() map[int]*Transaction {
	t.rmutex.RLock()
	defer t.rmutex.RUnlock()
	return t.transactions
}

func (t *TransactionDispetcher) GetOneTransaction(key int) bool {
	t.rmutex.RLock()
	defer t.rmutex.RUnlock()
	_, exist := t.transactions[key]
	return exist
}
