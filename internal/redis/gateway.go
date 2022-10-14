package redis

import (
	"fmt"
	"time"

	"eim/internal/types"
)

func RegisterGateway(gateway *types.Gateway) error {
	body, _ := gateway.Serialize()
	return rdsClient.Set(fmt.Sprintf("gateway:%v", gateway.Ip), body, time.Second*6)
}

func GetGateways() ([]*types.Gateway, error) {
	values, err := rdsClient.GetAll(fmt.Sprintf("gateway:*"))
	if err != nil {
		return nil, err
	}
	var gateways []*types.Gateway
	for _, value := range values {
		gateway := &types.Gateway{}
		err = gateway.Deserialize([]byte(value))
		if err != nil {
			continue
		}
		gateways = append(gateways, gateway)
	}
	return gateways, err
}
