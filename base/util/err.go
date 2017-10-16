package util

import (
	"log"
)

func PanicOnErr(e error) {
	if e != nil {
		log.Panic(e)
	}
}
