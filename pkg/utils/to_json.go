package utils

import "encoding/json"

func ToJson(encoding interface{}) ([]byte, error) {
	return json.Marshal(encoding)
}
