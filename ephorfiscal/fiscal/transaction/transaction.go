package transaction

type Config struct {
	Id                  int64
	Name                string
	Type                uint8
	Dev_interface       int
	AutomatNumber       int
	Login               string
	Password            string
	Phone               string
	Email               string
	Dev_addr            string
	Dev_port            int
	Ofd_addr            string
	Ofd_port            int
	Inn                 string
	Auth_public_key     string
	Auth_private_key    string
	Sign_private_key    string
	Param1              string
	Use_sn              int
	Add_fiscal          int
	Id_shift            string
	Fr_disable_cash     int
	Fr_disable_cashless int
	Ffd_version         int
	MaxSum              int
	QrFormat            int
	CancelCheck         int
}

type TransactionFiscal struct {
	Config
}
