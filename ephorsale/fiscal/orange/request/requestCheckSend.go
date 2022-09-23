package request

type RequestSendCheck struct {
	Id      string
	Inn     string
	Group   string
	Content struct {
		FfdVersion uint8
		Type       uint8
		Positions  []struct {
			quantity           int64
			price              int64
			tax                uint8
			text               string
			paymentMethodType  uint8
			paymentSubjectType uint8
		}
		CheckClose struct {
			taxationSystem uint8
		}
		CustomerContact   string
		AgentType         uint8
		SettlementAddress string
		SettlementPlace   string
	}
	Key                 string
	CallbackUrl         string
	Meta                string
	IgnoreItemCodeCheck bool
}
