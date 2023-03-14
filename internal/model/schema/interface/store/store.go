package store 

import (
	"ephorservices/pkg/model/schema/interface/model"
)

type Store interface {
	GetModel() model.Model
}