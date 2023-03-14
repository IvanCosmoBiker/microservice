package randgeneratestring

import (
	criptoRand "crypto/rand"
	"fmt"
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-")
var orderNumberSize int = 32

type GenerateString struct {
	String string
}

func Init() *GenerateString {
	return &GenerateString{}
}

func (gen *GenerateString) getDataOfTimeString() string {
	today := time.Now()
	hour := today.Hour()
	minute := today.Minute()
	second := today.Second()
	TimeString := fmt.Sprintf("%d%d%d", hour, minute, second)
	return TimeString
}

func (gen *GenerateString) RandStringRunes() {
	stringResult := ""
	b := make([]rune, orderNumberSize)
	for k := 0; k < len(b); k++ {
		if k == 0 && k < 3 {
			stringResult += string(letterRunes[rand.Intn(len(letterRunes))])
		}
	}
	stringResult += "-"
	Time := gen.getDataOfTimeString()
	stringResult += Time
	stringResult += "-"
	for k := 0; k < len(b); k++ {
		if k == 0 && k < 3 {
			stringResult += string(letterRunes[rand.Intn(len(letterRunes))])
		}
	}
	gen.String = stringResult

}

func (gen *GenerateString) RandGiud() {
	b := make([]byte, 16)
	criptoRand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	gen.String = uuid
}

func (gen *GenerateString) Strimwidth(str string, start, width int, trim_marker string) {
	result := []byte(str)
	if len(result) > width {
		result = result[start:width]
	}
	result = append(result, trim_marker...)
	gen.String = string(result[:])
}
