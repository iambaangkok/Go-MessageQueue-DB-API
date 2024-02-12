package utils

import (
	"encoding/json"
	"log"
	"models"
)

func CompressToJsonBytes(obj any) []byte {
	raw, _ := json.Marshal(obj)
 	return raw
}

func MsgToJson(msg models.Message) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}