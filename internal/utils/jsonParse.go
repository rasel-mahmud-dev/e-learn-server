package utils

import "encoding/json"

func ParseJson(str string, dest *any) {
	err := json.Unmarshal([]byte(str), dest)
	if err != nil {

	}
}
