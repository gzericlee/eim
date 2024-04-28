package model

import "eim/pkg/json"

type Gateway struct {
	Ip             string  `json:"ip"`
	ClientTotal    int64   `json:"clientTotal"`
	SendTotal      int64   `json:"sentTotal"`
	ReceivedTotal  int64   `json:"receivedTotal"`
	InvalidTotal   int64   `json:"invalidTotal"`
	GoroutineTotal int64   `json:"goroutineTotal"`
	MemUsed        float64 `json:"memUsed"`
	CpuUsed        float64 `json:"cpuUsed"`
}

func (its *Gateway) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *Gateway) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
