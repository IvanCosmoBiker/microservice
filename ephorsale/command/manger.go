package command

import (
	service_command_http "ephorservices/ephorsale/command/api/service_http"
	transaction_dispetcher "ephorservices/ephorsale/transaction"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"ephorservices/ephorsale/transport"
	storeEvent "ephorservices/internal/model/schema/account/automatevent/store"
	commandModel "ephorservices/internal/model/schema/main/command/model"
	storeCommand "ephorservices/internal/model/schema/main/command/store"
	logger "ephorservices/pkg/logger"
)

type CommandManager struct {
}

var ManagerCommand *CommandManager

func New() *CommandManager {
	commandManager := &CommandManager{}
	ManagerCommand = commandManager
	return commandManager
}

func (cm *CommandManager) InitApi() {
	serviceHttpHandler := service_command_http.New()
	serviceHttpHandler.CommandDeviceHandler = cm.SendCommandToDevice
	serviceHttpHandler.InitApi(transport.TransportManager.RequestManager)
}

func (cm *CommandManager) SendCommandToDevice(tran *transaction.Transaction) {
	logger.Log.Infof("%+v", tran)
	reqModem := tran.NewRequest()
	reqModem.AddFilterParam("imei", reqModem.Operator.OperatorEqual, true, tran.Config.Imei)
	modem, errModem := transaction_dispetcher.Dispetcher.StoreModem.GetOneBy(reqModem)
	if errModem != nil {
		logger.Log.Errorf("%v", errModem)
		return
	}
	reqCommand := tran.NewRequest()
	reqCommand.AddFilterParam("modem_id", reqCommand.Operator.OperatorEqual, true, modem.Id)
	reqCommand.AddFilterParam("sended", reqCommand.Operator.OperatorEqual, true, storeCommand.SendUnSuccess)
	models, _ := tran.Stores.StoreCommand.Get(reqCommand)
	for _, command := range models {
		logger.Log.Infof("%+v", command)
		requestModem := make(map[string]interface{})
		requestModem["a"] = command.Command.Int32
		requestModem["m"] = 2
		if command.Command_param1.Int32 > int32(0) {
			requestModem["sum"] = command.Command_param1.Int32
		}
		command.Sended.Scan(storeCommand.SendSuccess)
		command, err := tran.Stores.StoreCommand.Set(command)
		if err != nil {
			logger.Log.Error(err.Error())
		}
		logger.Log.Infof("MODELCOMMAND:: %+v", command)
		cm.AddEventAndSetCommand(command, tran, true)
		transport.TransportManager.QueueManager.SendMessage(requestModem, tran.Config.Imei)
	}
}

func (cm *CommandManager) AddEventAndSetCommand(command *commandModel.CommandModel, tran *transaction.Transaction, err bool) {
	reqAutomat := tran.NewRequest()
	reqAutomat.AddFilterParam("modem_id", reqAutomat.Operator.OperatorEqual, true, command.Modem_id.Int32)
	automat, errAutomat := tran.Stores.StoreAutomat.GetOneBy(reqAutomat)
	if errAutomat != nil {
		logger.Log.Errorf("%v", errAutomat)
		return
	}
	eventCommand := make(map[string]interface{})
	name := tran.Stores.StoreCommand.MakeResponseCommand(err, command.Command.Int32, command.Command_param1.Int32)
	date := transaction_dispetcher.Dispetcher.Date.Now()
	eventCommand["name"] = name
	eventCommand["date"] = date
	eventCommand["automat_id"] = automat.Id
	eventCommand["type"] = storeEvent.Type_ServerCommand
	logger.Log.Infof("%v", eventCommand)
	tran.Stores.StoreAutomatEvent.AddByParams(eventCommand)
}
