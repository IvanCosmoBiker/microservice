package response

type Response interface {
	GetStatusCode() int
}

type ResponseAuth struct {
	StatusCode int
	Status     string
	Error      struct {
		Code    int
		Message string
	}
	Data struct {
		AuthToken, ExpirationDateUtc string
	}
}

func (ra *ResponseAuth) GetStatusCode() int {
	return ra.StatusCode
}

type ResponseSendCheck struct {
	StatusCode int
	Status     string
	Error      struct {
		Code    int
		Message string
	}
	Data struct {
		ReceiptId string
	}
}

func (rsc *ResponseSendCheck) GetStatusCode() int {
	return rsc.StatusCode
}

type ResponseStatusCheck struct {
	StatusCode int
	Status     string
	Error      struct {
		Code    int
		Message string
	}
	Data struct {
		StatusCode          int
		StatusName          string
		StatusMessage       string
		Description         string
		ModifiedDateUtc     string
		ReceiptDateUtc      string
		ModifiedDateTimeIso string
		ReceiptDateTimeIso  string
		Device              struct {
			DeviceId              string
			RNM, ZN, FN, FDN, FPD string
			ReceiptNumInShift     int64
			OfdReceiptUrl         string
		}
	}
}

func (rsc *ResponseStatusCheck) GetStatusCode() int {
	return rsc.StatusCode
}
