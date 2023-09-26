package utils

import (
	"log"
)

func DEBUG() {
	log.Println(GetCaller(1))
}
