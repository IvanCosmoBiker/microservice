package response

type Response interface {
	GetStatusCode() int
}

type ResponseCreateOrder struct {
	StatusCode     int
	ErrorCode      string
	ErrorMessage   string
	OrderId        string
	FormUrl        string
	ExternalParams struct {
		SbolBankInvoiceId string
		SbolDeepLink      string
		SbolInactive      bool
	}
}

func (rco *ResponseCreateOrder) GetStatusCode() int {
	return rco.StatusCode
}

type ResponseStatusOrder struct {
	StatusCode            int
	ErrorCode             string
	ErrorMessage          string
	OrderNumber           string
	OrderStatus           uint8
	ActionCode            uint8
	ActionCodeDescription string
	Amount                int64
	Currency              uint8
	Date                  string
	DepositedDate         int64
	OrderDescription      string
	Ip                    string
	AuthRefNum            string
	RefundedDate          string
	PaymentWay            string
	AvsCode               string
}

func (rso *ResponseStatusOrder) GetStatusCode() int {
	return rso.StatusCode
}

type ResponseDebitOrder struct {
	StatusCode   int
	ErrorCode    string
	ErrorMessage string
}

func (rdo *ResponseDebitOrder) GetStatusCode() int {
	return rdo.StatusCode
}

type ResponseReturnOrder struct {
	StatusCode   int
	ErrorCode    string
	ErrorMessage string
}

func (rro *ResponseReturnOrder) GetStatusCode() int {
	return rro.StatusCode
}
