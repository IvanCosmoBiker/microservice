package request

type Position struct {
	Quantity           int64
	Price              float64
	Tax                uint8
	Text               string
	PaymentMethodType  uint8
	PaymentSubjectType uint8
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
	IgnoreItemCodeCheck bool
	Content             struct {
		FfdVersion    uint8
		Type          uint8
		AutomatNumber string
		Positions     []Position
		CheckClose    struct {
			Payments       []Payment
			TaxationSystem uint8
		}
		SettlementAddress string
		SettlementPlace   string
	}
}
