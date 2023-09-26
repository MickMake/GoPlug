package Plugin

import (
	"log"

	"github.com/MickMake/GoPlug/utils"
)

func DEBUG() {
	log.Println(utils.GetCaller(1))
}
