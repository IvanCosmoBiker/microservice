package response

type ResponseCreateQr struct {
	StatusCode int
}

func (r *ResponseCreateQr) GEtStatus() int {
	return r.StatusCode
}
