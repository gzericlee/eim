package httputil

import (
	"fmt"
	"net/url"
)

func ToQueryString(input map[string]interface{}) string {
	values := url.Values{}
	for key, value := range input {
		values.Set(key, fmt.Sprintf("%v", value))
	}
	return values.Encode()
}
