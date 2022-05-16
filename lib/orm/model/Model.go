package model

import (

)
type Model interface {
	InitData(json []byte)
	JsonSerialize()
	GetNameTable()
	GetNameSchema()
}



