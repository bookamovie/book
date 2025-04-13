package utils

import "encoding/json"

func MarshalJSON(v any) []byte {
	vByte, _ := json.Marshal(v)

	return vByte
}
