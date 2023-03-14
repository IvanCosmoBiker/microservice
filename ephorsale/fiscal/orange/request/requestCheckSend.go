package request

type IndustryAttribute struct {
	FoivId              string
	CauseDocumentDate   string
	CauseDocumentNumber string
	Value               string
}

type Position struct {
	ItemCode                string `json:"itemCode,omitempty"`
	Quantity                int64
	Price                   float64
	Tax                     uint8
	Text                    string
	PaymentMethodType       uint8
	PaymentSubjectType      uint8
	QuantityMeasurementUnit string             `json:"quantityMeasurementUnit,omitempty"`
	PlannedStatus           int                `json:"plannedStatus,omitempty"`
	IndustryAttribute       *IndustryAttribute `json:"IndustryAttribute,omitempty"`
}

type Payment struct {
	Type   uint8
	Amount float64
}
type RequestSendCheck struct {
	Id                  string
	Inn                 string
	Group               string
	Key                 string
	IgnoreItemCodeCheck bool `json:"ignoreItemCodeCheck,omitempty"`
	Content             struct {
		FfdVersion    uint8
		Type          uint8
		AutomatNumber string
		Positions     []*Position
		CheckClose    struct {
			Payments       []*Payment
			TaxationSystem uint8
		}
		SettlementAddress string
		SettlementPlace   string
		FsItemCodeType    string `json:"fsItemCodeType,omitempty"`
	}
}
