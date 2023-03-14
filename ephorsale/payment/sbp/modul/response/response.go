package modul

type Response interface {
	GetStatusCode() int
}

type ResponseGetQR struct {
	StatusCode int
	QrcId      string
	Payload    string
	Massage    string
}

func (r *ResponseGetQR) GetStatusCode() int {
	return r.StatusCode
}

type ResponseGetStatus struct {
	StatusCode     int
	QrcId          string
	LocalQrcId     string
	Amount         float64
	Type           string
	Status         string
	PaymentPurpose string
	Payload        string
	OperationId    string
	Massage        string
}

func (r *ResponseGetStatus) GetStatusCode() int {
	return r.StatusCode
}

type ResponseReturnQr struct {
	StatusCode int
	RequestId  string
	Massage    string
}

func (r *ResponseReturnQr) GetStatusCode() int {
	return r.StatusCode
}
