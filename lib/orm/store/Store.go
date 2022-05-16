package store 

import (

)
type Store interface {
    Get()
    Set(map[string]interface{})
    GetWithOptions(options map[string]interface{})
    SetWithOptions(options map[string]interface{})
    AddByParams(parametrs map[string]interface{})
    SetByParams(parametrs map[string]interface{})
    GetOneById(id interface{})
}

