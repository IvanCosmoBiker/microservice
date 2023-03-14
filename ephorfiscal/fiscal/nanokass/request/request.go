package request

type Position struct {
	Name_tovar                string
	Stavka_nds                int
	Priznak_agenta            string
	Kolvo                     int64
	Priznak_sposoba_rascheta  uint8
	Priznak_predmeta_rascheta uint8
	Price_piece_bez_skidki    float64
	Price_piece               float64
	Summa                     float64
}

type Payment struct {
	Dop_rekvizit_1192 string
	Inn_pokupatel     string
	Name_pokupatel    string
	Rezhim_nalog      string
	Kassir_inn        string
	Kassir_fio        string
	Client_email      string
	Money_nal         float64
	Money_electro     float64
	Money_predoplata  float64
	Money_postoplata  float64
	Money_vstrecha    float64
}

type RequestSendCheck struct {
	Kassaid                string
	Kassatoken             string
	Cms                    string
	Check_send_type        string
	Check_vend_address     string
	Check_vend_mesto       string
	Check_vend_num_avtovat string
	Products_arr           []Position
	Oplata_arr             []Payment
	Itog_arr               struct {
		Priznak_rascheta uint8
		Itog_cheka       float64
	}
}
