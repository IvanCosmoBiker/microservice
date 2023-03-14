package request

type RequestAuthToken struct {
	Login, Password string
}

type MarkingCodeData struct {
	Type          string //UNKNOWN_PRODUCT_CODE
	Code          string
	PlannedStatus string //PRODUCT_STATUS_NOT_CHANGED
}

type Position struct {
	Label           string
	Price           float64
	Quantity        float64
	Amount          float64
	Vat             string
	PaymentMethod   uint8
	PaymentType     uint8
	MarkingCodeData *MarkingCodeData `json:"MarkingCodeData,omitempty"`
}

type Payment struct {
	PaymentType uint8
	Sum         float64
}

type RequestSendCheck struct {
	Request struct {
		Inn             string
		Type            string
		InvoiceId       string
		CallbackUrl     string
		CustomerReceipt struct {
			KktFA                 bool
			TaxationSystem        string
			Email                 string
			Phone                 string
			PaymentType           uint8
			AutomaticDeviceNumber int32
			BillAddress           string
			Items                 []*Position
			PaymentItems          []*Payment
		}
	}
}

type RequestStatusCheck struct {
	Request struct {
		ReceiptId string
	}
}
