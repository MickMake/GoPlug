package utils

import (
	"log"
)

func DEBUG() {
	log.Println(GetCallerDebug(1))
}
