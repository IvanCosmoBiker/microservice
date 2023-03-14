package store

import (
	command_model "ephorservices/internal/model/schema/main/command/model"
	model_interface "ephorservices/pkg/orm/model"
	request "ephorservices/pkg/orm/request"
	store_parent "ephorservices/pkg/orm/store"
	"fmt"
)

var (
	Command_None                    = 0x00
	Command_Setting                 = 0x01
	Command_ReloadModem             = 0x02
	Command_ReloadAutomat           = 0x03
	Command_ResetErrors             = 0x04
	Command_ChargeCash              = 0x05
	Command_LoadAudit               = 0x06
	Command_LoadSoftWare            = 0x07
	Command_SverkaItogov            = 0x08
	Command_SbrosSchetchikaAudita   = 0x09
	Command_LoadStaticAudit         = 0x000A
	Command_Remote_Flush            = 0x000B
	Command_Update_Cashless_Fimware = 0x000C
	Command_Service_Mode            = 0x000D
	Command_Off_Payment             = 0x000E
	Command_CashlessLog             = 0x000F
)

var (
	SendUnSuccess = 0
	SendSuccess   = 1
)

type StoreCommand struct {
	*store_parent.Store
}

func New(accountNumber ...int) *StoreCommand {
	AccountNumber := 0
	if len(accountNumber) != 0 {
		AccountNumber = accountNumber[0]
	}
	Model := command_model.New()
	store := &StoreCommand{
		store_parent.New(Model),
	}
	store.Store.SetAccountNumber(AccountNumber)
	return store
}

func (sc *StoreCommand) GetStructModel(model model_interface.Model) *command_model.CommandModel {
	if model != nil {
		return model.(*command_model.CommandModel)
	}
	return &command_model.CommandModel{}
}

func (sc *StoreCommand) Get(req *request.Request) (Models []*command_model.CommandModel, err error) {
	models, err := sc.Store.Get(req)
	if len(models) > 0 {
		Models = make([]*command_model.CommandModel, 0, len(models))
		for _, model := range models {
			Models = append(Models, sc.GetStructModel(model))
		}
	}
	return
}

func (sc *StoreCommand) AddByParams(params map[string]interface{}) (*command_model.CommandModel, error) {
	model, err := sc.Store.AddByParams(params)
	Model := sc.GetStructModel(model)
	return Model, err
}

func (sc *StoreCommand) SetByParams(params map[string]interface{}) (*command_model.CommandModel, error) {
	model, err := sc.Store.SetByParams(params)
	Model := sc.GetStructModel(model)
	return Model, err
}

func (sc *StoreCommand) GetOneById(id int) (*command_model.CommandModel, error) {
	model, err := sc.Store.GetOneById(id)
	Model := sc.GetStructModel(model)
	return Model, err
}

func (sc *StoreCommand) GetOneBy(req *request.Request) (*command_model.CommandModel, error) {
	model, err := sc.Store.GetOneBy(req)
	Model := sc.GetStructModel(model)
	return Model, err
}

func (sc *StoreCommand) Set(model model_interface.Model) (*command_model.CommandModel, error) {
	model, err := sc.Store.Set(model)
	Model := sc.GetStructModel(model)
	fmt.Printf("%+v\n", Model)
	return Model, err
}

func (sc *StoreCommand) GetNameCommand(command, command_param1 int32) string {
	var returningString string
	switch int(command) {
	case Command_Setting:
		returningString = "загрузка настроек"
		return returningString
	case Command_ReloadModem:
		returningString = "перезагрузка модема"
		return returningString
	case Command_ReloadAutomat:
		returningString = "перезагрузка автомата"
		return returningString
	case Command_ResetErrors:
		returningString = "сброс ошибок"
		return returningString
	case Command_ChargeCash:
		returningString = fmt.Sprintf("Кредит в размере %2.f", float64(command_param1/100))
		return returningString
	case Command_LoadAudit:
		returningString = "выгрузка аудита"
		return returningString
	case Command_LoadSoftWare:
		returningString = "обновление прошивки"
		return returningString
	case Command_SverkaItogov:
		returningString = "сверка итогов"
		return returningString
	case Command_SbrosSchetchikaAudita:
		returningString = "сброс счетчиков аудита"
		return returningString
	case Command_Remote_Flush:
		returningString = "промывка"
		return returningString
	case Command_Update_Cashless_Fimware:
		returningString = "обновление прошивки безнала"
		return returningString
	case Command_Service_Mode:
		returningString = "режим обслуживания"
		return returningString
	case Command_Off_Payment:
		returningString = "отключение платежек"
		return returningString
	case Command_CashlessLog:
		returningString = "выгрузка логов"
		return returningString
	}
	return returningString
}

func (sc *StoreCommand) MakeResponseCommand(err bool, command, command_param1 int32) string {
	name := sc.GetNameCommand(command, command_param1)
	if !err {
		if command == int32(Command_ChargeCash) {
			return fmt.Sprintf("%s %s", name, "не начислен")
		}
		return fmt.Sprintf("%s %s %s", "Команда", name, "не отправлена")
	}
	if command == int32(Command_ChargeCash) {
		return fmt.Sprintf("%s %s", name, "начислен")
	}
	return fmt.Sprintf("%s %s %s", "Команда", name, "отправлена")
}
