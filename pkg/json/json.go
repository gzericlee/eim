package json

import (
	jsoniter "github.com/json-iterator/go"
)

func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
