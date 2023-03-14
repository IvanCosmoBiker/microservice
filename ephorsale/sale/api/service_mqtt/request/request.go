package request

import (
	"encoding/json"
)

type RequestPay struct {
}

type ResponsePay struct {
	D   string
	A   int
	Tid int
	St  int
	Wid int
	Sum int64
	Err int32
}

func (rp *ResponsePay) JsonToStruct(data []byte) error {
	err := json.Unmarshal(data, rp)
	return err
}
