package maputil

import (
	"errors"

	"github.com/gzericlee/eim/pkg/conv/typeconv"
)

func GetString(vars map[string]interface{}, key, defaultVal string) string {
	if val, ok := vars[key]; ok {
		return val.(string)
	} else {
		return defaultVal
	}
}

func GetAny[T any](data map[string]interface{}, keys ...string) (T, error) {
	var tmp interface{}
	var result T
	var tmpMap map[string]interface{}
	for i, key := range keys {
		if i == 0 {
			tmp = data[key]
		} else {
			if tmpMap == nil {
				return result, errors.New("not a valid map")
			}
			tmp = tmpMap[key]
		}
		if tmp == nil {
			return result, errors.New("key not found")
		}
		tmpMap, _ = tmp.(map[string]interface{})
	}
	return typeconv.ToAnyE[T](tmp)
}

func GetInterface(data map[string]interface{}, keys ...string) (interface{}, error) {
	var result interface{}
	var tmpMap map[string]interface{}
	for i, key := range keys {
		if i == 0 {
			result = data[key]
		} else {
			if tmpMap == nil {
				return nil, errors.New("not a valid map")
			}
			result = tmpMap[key]
		}
		if result == nil {
			return nil, errors.New("key not found")
		}
		tmpMap, _ = result.(map[string]interface{})
	}
	return result, nil
}

func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
