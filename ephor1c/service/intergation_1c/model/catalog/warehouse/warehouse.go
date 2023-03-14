package warehouse

type Warehouse struct {
	Ref_Key                                           string
	DataVersion                                       string
	DeletionMark                                      bool
	Parent_Key                                        string
	IsFolder                                          bool
	Description                                       string
	UseAddresStore                                    bool `json:"ИспользоватьАдресноеХранение"`
	UseAddresStoreForRef                              bool `json:"ИспользоватьАдресноеХранениеСправочно"`
	UseTheOrderSchemeWhenShipping                     bool `json:"ИспользоватьОрдернуюСхемуПриОтгрузке"`
	UseTheOrderSchemeWhenReflectingSurplusesShortages bool `json:"ИспользоватьОрдернуюСхемуПриОтраженииИзлишковНедостач"`
	UseTheOrderOfTheSchemeInTheAccess                 bool `json:"ИспользоватьОрдернуюСхемуПриПоступлении"`
	UseTheSeriesOfTheNomenclature                     bool `json:"ИспользоватьСерииНоменклатуры"`
	UseStorageFacilities                              bool `json:"ИспользоватьСкладскиеПомещения"`
}
