package jsonutil

import "encoding/json"

func DeepCopy[T any](orig T) (result T, err error) {
	origJSON, err := json.Marshal(orig)
	if err != nil {
		return result, err
	}

	return result, json.Unmarshal(origJSON, &result)
}
