package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	config "ephorservices/config"
	request "ephorservices/ephorsale/fiscal/nanokass/request"
	transaction "ephorservices/ephorsale/transaction/transaction_struct"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
)

var RSA_PUB_FIRST = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwFXHnzc5YKj8e3tlNzST
CkA8Tq4gjTH0VMuhJhg5QWpFjFKwtnK3u4EOaQGmjqDtzyffVHmKuGikg9jE20sG
nJN4hTtySihOiUWRd4zhJVMevBQmsEQS33bg26UzzKCeO12mbM/Q4ip7YXEfWM/F
Tq2l94psQgmIDh/LtHVf3OBlz8I6u5VaP3AS0Hv9RBUin0RBkRUC+5tgURm382XT
nJ2GzZ8cEGJm3C+s0+W1N2igjV0X3MihylHGDyl+8FpbFIlXsaJOYQ0//JIgnaBz
MV2JyNTHBzPJrcIMHIbKBVAmDLfgeDNKug7wIadEcqoJaCz74yG9l9nJWISWQkI6
Ed8nDVsoaIkMQBuWWxfHjQEU8R8OVjRzhOGHPG2ka6y1/jcOS5JWPzS5YVXRPbrh
QYcoNebsOBaFxJYZ2E7VhVdrGWlBqhANFba7umZXVOvmDXIsH974Yv4awAaP70VP
SLFIdjiNy/SB8w0O8PJOUPznpMhvi1clBgp3PvtYmhUqmdHWPwjcjy0JmY9KrWz0
0Im1yDTTybtV3uYnwR677TmsLmR9c6T7EHlT3gG6Y0bM3w9tyrGqVKy1jIkyUZPV
f0dmXTfbh+hcC5kYal+M7lcn7wSSLHTUk+C/YWE1e5TvTBK6teU2VNmz80Yt2IS2
mcXlfKlZXilMmPJCdUI7nNMCAwEAAQ==
-----END PUBLIC KEY-----`)

var RSA_PUB_SECOND = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA+fu+NGlnWAXqIVgEL37v
eatlyooYi+iHLiBmCDowNZUBAiQ+pvbnzkowUKdr86lGrzQLCAvVyXWG0U4kdixA
X0GTkIR/3g3h2/8hRx0x3K0umT+tcZC3iJytKzP+EM/B6sDdw6/URbykwvrAlbQs
G9d6eCqq0F/6muOM3gQazy8CuHyx4iFQpml4E1/IQgp3tZJOX5I9xieHTUwct2Ok
URCKYnHJZrRIN9rwXQkNG1q+M8HDqI1Mwq88wieVC+SUuoPc8F0MlIWs2zwDhLcX
84OQTRFqlW3NFR/6kUn3TIC1JZD1Ft/8fWukZzAFsAmdXmFzhBUuBPvjIzzLafY3
f8IszADMnloJ0BW3iGVRGj6hygX7Jpr/86LPHu6PBJzHzCp9bnfOiSjRENzzy55f
DdVbYpVgWDt4+UEkl9qNRNuiSMDpKeVNy6jxbihZneYCR8alnH8Olh6lL7bmGdww
qI9LSyq/qFfIMDV8onit/dLxzypFJofRfjZ1Dc8ZEqh2sab8qEMNPGQwTM/FVFWM
bq0hmjjY+BFWGY/h0z1NZMX75Uzyd9OdXaRoTlHPfOxxAIfclP2XY2K8f5PQ37g/
fX2R8bw/fXQd2ndi/+uPCGK92Xw4/3/osJKpm3QSYhSda53T9Ddned7BtWDQJqdV
Y/SUskwLLyjtSb0LqsSKBHkCAwEAAQ==
-----END PUBLIC KEY-----`)

const HMAC_FIRST = "BBuXaXBdHg+wLPjRJpf3N/NmLq5kuvzGQx3II15/j8o="
const HMAC_SECOND = "aFZP3PbvrMZNNxxqJxaCnCLama5L8H1/YGO3UYsoCVQ="
const URL_TO_SEND_TO_NANOKASSA = "http://q.nanokassa.ru/srv/igd.php"

type Core struct {
	cfg *config.Config
}

func InitCore(conf *config.Config) *Core {
	return &Core{
		cfg: conf,
	}
}

func (c *Core) MakeRequestSendCheck(tran *transaction.Transaction) (jsonData []byte, url string, err error) {
	if tran.Fiscal.InCome.Imei != "" {
		jsonData, _ = json.Marshal(tran.Fiscal.InCome.Fields.Request)
		url = URL_TO_SEND_TO_NANOKASSA
		return
	}
	requestSendCheck := request.RequestSendCheck{}
	paramsNano := tran.Fiscal.Config.Sign_private_key
	dataNano := strings.Split(paramsNano, ":")
	log.Printf("%v", dataNano)
	requestSendCheck.Kassaid = dataNano[0]
	requestSendCheck.Kassatoken = dataNano[1]
	requestSendCheck.Cms = "wordpress"
	requestSendCheck.Check_send_type = "email"
	requestSendCheck.Check_vend_address = tran.Point.Address
	requestSendCheck.Check_vend_mesto = tran.Point.PointName
	requestSendCheck.Check_vend_num_avtovat = fmt.Sprintf("%v", tran.Config.AutomatId)
	c.GenerateDataForCheck(&requestSendCheck, tran)
	requestSendCheck.Itog_arr.Itog_cheka = tran.Fiscal.SumCheck
	requestSendCheck.Itog_arr.Priznak_rascheta = 1
	request1, _ := json.Marshal(requestSendCheck)
	firstcrypt := c.Crypt_nanokassa_first(request1)
	returnDataAB := firstcrypt["ad"]
	returnDataDE := firstcrypt["de"]
	content2 := make(map[string]interface{})
	content2["ab"] = fmt.Sprintf("'%v'", returnDataAB)
	content2["de"] = fmt.Sprintf("'%v'", returnDataDE)
	content2["kassaid"] = fmt.Sprintf("'%s'", dataNano[0])
	content2["kassatoken"] = fmt.Sprintf("'%s'", dataNano[1])
	content2["check_type"] = "standart"
	content2["test"] = "1"
	request2, _ := json.Marshal(content2)
	secondcrypt := c.Crypt_nanokassa_second(request2)
	returnDataAAB := secondcrypt["aab"]
	returnDataDE2 := secondcrypt["dde"]
	content3 := make(map[string]interface{})
	content3["aab"] = fmt.Sprintf("'%v'", returnDataAAB)
	content3["aab"] = fmt.Sprintf("'%v'", returnDataDE2)
	content3["test"] = "0"
	jsonData, _ = json.Marshal(content3)
	url = URL_TO_SEND_TO_NANOKASSA
	return
}

func (c *Core) MakeRequestStatusCheck(tran *transaction.Transaction) (url string) {
	dataNano := strings.Split(tran.Fiscal.ResiptId, ":")
	if len(dataNano) < 2 {
		return
	}
	url = fmt.Sprintf("https://fp.nanokassa.com/getfp?nuid=%s&qnuid=%s&auth=base", dataNano[0], dataNano[1])
	return
}

func (c *Core) GenerateRandomBytesInString(n int) (b []byte) {
	b = make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return
	}
	return
}

func (c *Core) EncryptAES(plaintext []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Режим CTR заполнять не нужно, возвращает режим счетчика, нижний уровень использует блок для генерации интерфейса srtream ключевого потока
	stream := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(plaintext))
	// Операция шифрования
	stream.XORKeyStream(cipherText, plaintext)
	return cipherText, nil
}

func (c *Core) RsaEncrypt(origData, publickey []byte) ([]byte, error) {
	rng := rand.Reader
	block, _ := pem.Decode(publickey)
	if block == nil {
		return nil, errors.New("is not pem")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	signature, err := rsa.EncryptPKCS1v15(rng, pub, origData)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (c *Core) HMAC512(ciphertext, key []byte) []byte {
	mac := hmac.New(sha512.New, key)
	io.WriteString(mac, string(ciphertext))
	return mac.Sum(nil)
}

func (c *Core) Crypt_nanokassa_first(data []byte) map[string]interface{} {
	resultMap := make(map[string]interface{})
	IVdata := c.GenerateRandomBytesInString(16)
	pw := c.GenerateRandomBytesInString(32)
	mk := HMAC_FIRST
	dataAES, err := c.EncryptAES(data, pw, IVdata)
	if err != nil {
		log.Println(err)
		return resultMap
	}
	keyHMAC, _ := base64.StdEncoding.DecodeString(mk)
	hmac512 := make([]byte, 0)
	hmac512 = append(hmac512, IVdata...)
	hmac512 = append(hmac512, dataAES...)
	hmac := c.HMAC512(hmac512, keyHMAC)
	DataDe := make([]byte, 0)
	DataDe = append(DataDe, hmac...)
	DataDe = append(DataDe, IVdata...)
	DataDe = append(DataDe, dataAES...)
	returnDataDE := base64.StdEncoding.EncodeToString(DataDe)
	ab_rsa, errRsa := c.RsaEncrypt([]byte(pw), RSA_PUB_FIRST)
	if errRsa != nil {
		log.Println(errRsa)
		return resultMap
	}
	returnDataAB := base64.StdEncoding.EncodeToString(ab_rsa)
	resultMap["ab"] = returnDataAB
	resultMap["de"] = returnDataDE
	return resultMap
}

func (c *Core) Crypt_nanokassa_second(data []byte) map[string]interface{} {
	resultMap := make(map[string]interface{})
	IVdata := c.GenerateRandomBytesInString(16)
	pw := c.GenerateRandomBytesInString(32)
	mk := HMAC_SECOND
	dataAES, err := c.EncryptAES(data, pw, IVdata)
	if err != nil {
		log.Println(err)
		return resultMap
	}
	keyHMAC, _ := base64.StdEncoding.DecodeString(mk)
	hmac512 := make([]byte, 0)
	hmac512 = append(hmac512, IVdata...)
	hmac512 = append(hmac512, dataAES...)
	hmac := c.HMAC512(hmac512, keyHMAC)
	DataDee := make([]byte, 0)
	DataDee = append(DataDee, hmac...)
	DataDee = append(DataDee, IVdata...)
	DataDee = append(DataDee, dataAES...)
	returnDataDEE := base64.StdEncoding.EncodeToString(DataDee)
	aab_rsa, _ := c.RsaEncrypt([]byte(pw), RSA_PUB_SECOND)
	returnDataAAB := base64.StdEncoding.EncodeToString(aab_rsa)
	resultMap["aab"] = returnDataAAB
	resultMap["dde"] = returnDataDEE
	return resultMap
}

func (c *Core) GenerateDataForCheck(req *request.RequestSendCheck, tran *transaction.Transaction) {
	var summ int
	entryPayments := request.Payment{}
	entryPayments.Dop_rekvizit_1192 = ""
	entryPayments.Inn_pokupatel = ""
	entryPayments.Name_pokupatel = ""
	entryPayments.Rezhim_nalog = c.ConvertTaxationSystem(tran.TaxSystem.Type)
	entryPayments.Kassir_inn = ""
	entryPayments.Kassir_fio = ""
	entryPayments.Client_email = "none"
	entryPayments.Money_nal = 0
	entryPayments.Money_predoplata = 0
	entryPayments.Money_postoplata = 0
	entryPayments.Money_vstrecha = 0
	entryPayments.Money_electro = 0
	for _, product := range tran.Products {
		entryPositions := request.Position{}
		var value float64 = 0
		if product.Value == value {
			value = product.Price
		} else {
			value = product.Value
		}
		quantity := float64(product.Quantity / 1000)
		price := value
		summ += int(math.Round(quantity * price))
		entryPayments.Money_electro = math.Round(quantity * price)
		entryPositions.Summa = math.Round(quantity * price)
		entryPositions.Price_piece_bez_skidki = math.Round(price)
		entryPositions.Priznak_sposoba_rascheta = 4
		entryPositions.Priznak_predmeta_rascheta = 1
		entryPositions.Kolvo = int64(quantity)
		entryPositions.Name_tovar = product.Name
		entryPositions.Price_piece = math.Round(price)
		entryPositions.Stavka_nds = c.ConvertTax(int(product.Tax_rate))
		entryPositions.Priznak_agenta = "none"
		req.Products_arr = append(req.Products_arr, entryPositions)
	}
	entryPayments.Money_electro = float64(summ)
	tran.Fiscal.SumCheck = float64(summ)
	req.Oplata_arr = append(req.Oplata_arr, entryPayments)
}

func (c *Core) EncodeUrlToBase64(tran *transaction.Transaction) string {
	stringUrl := c.MakeUrlQr(tran)
	return base64.StdEncoding.EncodeToString([]byte(stringUrl))
}

func (c *Core) MakeUrlQr(tran *transaction.Transaction) string {
	stringResult := fmt.Sprintf("t=%s&s=%v&fn=%v&i=%v&fp=%v&n=1", tran.Fiscal.Fields.DateFisal, fmt.Sprintf("%.2f", tran.Fiscal.SumCheck), tran.Fiscal.Fields.Fn, tran.Fiscal.Fields.Fd, tran.Fiscal.Fields.Fp)
	log.Println(stringResult)
	return stringResult
}

func (c *Core) ConvertTax(tax int) int {
	switch tax {
	case transaction.TaxRate_NDSNone:
		return 6
		fallthrough
	case transaction.TaxRate_NDS0:
		return 5
		fallthrough
	case transaction.TaxRate_NDS10:
		return 2
		fallthrough
	case transaction.TaxRate_NDS18:
		return 1
	}
	return 6
}

func (c *Core) ConvertTaxationSystem(taxsystem int) string {
	switch taxsystem {
	case transaction.TaxSystem_OSN:
		return "0"
		fallthrough
	case transaction.TaxSystem_USND:
		return "1"
		fallthrough
	case transaction.TaxSystem_USNDMR:
		return "2"
		fallthrough
	case transaction.TaxSystem_ENVD:
		return "3"
		fallthrough
	case transaction.TaxSystem_ESN:
		return "4"
		fallthrough
	case transaction.TaxSystem_Patent:
		return "5"
	}
	return "0"
}
