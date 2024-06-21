package jsonutil

import "encoding/json"

func ToAny[T any](src any) (dst T, err error) {
	var body []byte

	switch src.(type) {
	case string:
		{
			body = []byte(src.(string))
			break
		}
	case []byte:
		{
			body = src.([]byte)
		}
	case interface{}:
		{
			body, err = json.Marshal(src)
			if err != nil {
				return
			}
			break
		}
	}

	switch any(dst).(type) {
	case string:
		{
			dst = any(string(body)).(T)
		}
	case []byte:
		{
			dst = any(body).(T)
		}
	case interface{}:
		{
			err = json.Unmarshal(body, &dst)
			if err != nil {
				return
			}
			break
		}
	}

	return
}

func ToIndentAny[T any](src any) (dst T, err error) {
	var body []byte

	switch src.(type) {
	case string:
		{
			body = []byte(src.(string))
			break
		}
	case []byte:
		{
			body = src.([]byte)
		}
	case interface{}:
		{
			body, err = json.MarshalIndent(src, "", "    ")
			if err != nil {
				return
			}
			break
		}
	}

	switch any(dst).(type) {
	case string:
		{
			dst = any(string(body)).(T)
		}
	case []byte:
		{
			dst = any(body).(T)
		}
	case interface{}:
		{
			err = json.Unmarshal(body, &dst)
			if err != nil {
				return
			}
			break
		}
	}

	return
}

func ToString(src any) (string, error) {
	switch src.(type) {
	case string:
		{
			return src.(string), nil
		}
	}
	b, err := json.Marshal(src)
	return string(b), err
}
