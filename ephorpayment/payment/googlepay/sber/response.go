package sber

type Response struct {
	StatusCode int
	Data       map[string]interface{}
}

func (r Response) GetData(field string) interface{} {
	value, _ := r.Data[field]
	return value
}
