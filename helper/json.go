package helper

import "encoding/json"

func JsonEncode(value interface{}) []byte {
	m, _ := json.Marshal(value)

	return m
}
