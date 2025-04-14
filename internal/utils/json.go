package utils

import "encoding/json"

// MarshalJSON() marshals a given value (v) into a JSON byte slice.
//
// It uses the standard library's json.Marshal function to convert the value into JSON format. The function does not handle errors explicitly, it will silently ignore any errors.
func MarshalJSON(v any) []byte {
	vByte, _ := json.Marshal(v)

	return vByte
}
