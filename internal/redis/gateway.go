package redis

import (
	"fmt"
	"time"

	"eim/internal/model"
)

func (its *Manager) RegisterGateway(gateway *model.Gateway) error {
	body, _ := gateway.Serialize()
	return its.rdsClient.Set(fmt.Sprintf("gateway:%v", gateway.Ip), body, time.Second*6)
}

func (its *Manager) GetGateways() ([]*model.Gateway, error) {
	values, err := its.rdsClient.GetAll(fmt.Sprintf("gateway:*"))
	if err != nil {
		return nil, err
	}
	var gateways []*model.Gateway
	for _, value := range values {
		gateway := &model.Gateway{}
		err = gateway.Deserialize([]byte(value))
		if err != nil {
			continue
		}
		gateways = append(gateways, gateway)
	}
	return gateways, err
}
