package response

type Response interface {
	GetStatus() int
}
type ResponseCreateCheck struct {
	StatusCode int
	Errors     []string
}

func (rcc *ResponseCreateCheck) GetStatus() int {
	return rcc.StatusCode
}

type ResponseStatusCheck struct {
	StatusCode     int
	Id             string
	DeviceSN       string
	DeviceRN       string
	FsNumber       string
	OfdName        string
	OfdWebsite     string
	Ofdinn         string
	FnsWebsite     string
	CompanyINN     string
	CompanyName    string
	DocumentNumber int64
	ShiftNumber    int64
	DocumentIndex  int64
	ProcessedAt    string
	Content        struct {
		FfdVersion uint8
		Type       uint8
		Positions  []struct {
			Quantity           int64
			Price              float64
			Tax                uint8
			Text               string
			PaymentMethodType  uint8
			PaymentSubjectType uint8
		}
		CheckClose struct {
			Payments []struct {
				Type   uint8
				Amount float64
			}
			taxationSystem uint8
		}
		CustomerContact   string
		AgentType         uint8
		SettlementAddress string
		SettlementPlace   string
	}
	Change float64
	Fp     string
	Errors []string
}

func (rsc *ResponseStatusCheck) GetStatus() int {
	return rsc.StatusCode
}
