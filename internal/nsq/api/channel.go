package api

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func CreateChannel(endpoint, topic, channel string) error {
	code, _, err := fasthttp.Post(nil, fmt.Sprintf("http://%v/channel/create?topic=%v&channel=%v", endpoint, topic, channel), nil)
	if err != nil {
		return err
	}
	if code > 300 {
		return fmt.Errorf("createing nsq channel: %v", err)
	}
	return nil
}
