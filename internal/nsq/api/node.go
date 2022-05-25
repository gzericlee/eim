package api

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"eim/pkg/json"
)

type Node struct {
	BroadcastAddress string `json:"broadcast_address"`
	HttpPort         int    `json:"http_port"`
	TcpPort          int    `json:"tcp_port"`
}

func GetNodes(endpoint string) ([]*Node, error) {
	code, body, err := fasthttp.Get(nil, fmt.Sprintf("http://%v/nodes", endpoint))
	if err != nil {
		return nil, err
	}
	var result struct {
		Nodes []*Node `json:"producers"`
	}
	if code > 300 {
		return nil, fmt.Errorf("geting nsq nodes: %v", err)
	}
	err = json.Unmarshal(body, &result)
	return result.Nodes, nil
}
