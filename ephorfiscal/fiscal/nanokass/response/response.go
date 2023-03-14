package response

type Response interface {
	GetStatusCode() int
}

type ResponseSendCheck struct {
	StatusCode int
	Status     string
	Error      string
	Code       string
	Success    string
	Nuid       string
	Qnuid      string
}

func (rsc *ResponseSendCheck) GetStatusCode() int {
	return rsc.StatusCode
}

type ResponseStatusCheck struct {
	StatusCode             int
	Error                  string
	Code                   string
	Check_nuid             string
	Check_qnuid            string
	Check_status           int
	Check_status_info      string
	Check_name             string
	Check_type             string
	Check_kkt_operator     string
	Check_sno              string
	Check_vend_address     string
	Check_vend_mesto       string
	Check_vend_num_avtovat string
	Check_dt_unixtime      string
	Check_dt_ofdtime       string
	Check_num_fd           float64
	Check_num_fp           int64
	Check_fn_num           string
	Check_site_fns         string
	Check_qr_code          string
	Check_qr_code_nano_url string
	Check_qr_code_img_url  string
	Check_qr_code_img_b64  string
	Check_itog             string
	Status_code            string
	Error_code             string
}

func (rsc *ResponseStatusCheck) GetStatusCode() int {
	return rsc.StatusCode
}
