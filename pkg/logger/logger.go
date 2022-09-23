package logger

import (
    "log"
)

func LoggingRecoverError() {
    if r := recover(); r != nil {
        log.Println("Recovered error", r)
    }
}
