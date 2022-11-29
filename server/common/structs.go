package common

import (
	"encoding/json"
	"log"
)

func Struct2JSON(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		log.Print("[Warning] Couldn't convert to json", err.Error())
		return ""
	}
	return string(b)
}
