package fiscalinterface 

import (
)
var (
     Fr_None = 0;
	 Fr_PayOnlineFA = 1;
	 Fr_KaznachejFA = 2;
	 Fr_RPSystem1FA = 3;
	 Fr_TerminalFA = 4;
	 Fr_OrangeData = 5;
	 Fr_ChekOnline = 6;
	 Fr_EphorOrangeData = 7;
	 Fr_EphorOnline = 8;
	 Fr_NanoKassa = 9;
	 Fr_ServerOrangeData = 10;
	 Fr_EphorServerOrangeData = 11;
	 Fr_OFD = 12;
	 Fr_ServerNanoKassa = 13;
	 Fr_ProSistem = 14;
	 Fr_CheckBox = 15;
)
// -- Status Fiscal -- //
var (
	 Status_Fr_None		 	 = 0; // продажа за 0 рублей или касса отключена
	 Status_Fr_Complete		 = 1; // чек создан успешно, реквизиты получены
	 Status_Fr_Manual			 = 2; // чек фискализирован вручную
	 Status_Fr_Need			 = 3; // чек необходимо фискализировать
	 Status_Fr_OFF_FR		 	 = 5; // оключение фискализации со стороны клиента
	 Status_Fr_OFF_DA		 	 = 6; // отключение фискализации для безналичной оплаты
	 Status_Fr_OFF_CA		 	 = 7; // отключение фискализации для оплаты за наличку
	 Status_Fr_Check_Cancel	 = 8; // отмена чека
    )
// -- Error Fiscal -- //
var (
	 Status_Fr_Fr_InQueue		 = 101; // чек добавлен в очередь, реквизиты не получены
	 Status_Fr_Unknown		 = 102; // результат постановки чека в очередь не известен
	 Status_Fr_Error		 = 103; // ошибка создания чека
	 Status_Fr_Overflow	 = 104; // очередь отложенной регистрации переполнена
	 Status_Fr_MAX			 = 105; // превышена максимальная сумма чека
    )
type Fiscal interface {
    InitData(transactionStruct.Transaction)
    SendCheck() map[string]interface{}
    GetStatus() map[string]interface{}
}