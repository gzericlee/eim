package api

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func CreateTopic(endpoint, topic string) error {
	code, _, err := fasthttp.Post(nil, fmt.Sprintf("http://%v/topic/create?topic=%v", endpoint, topic), nil)
	if err != nil {
		return err
	}
	if code > 300 {
		return fmt.Errorf("createing nsq topic: %v", err)
	}
	return nil
}
