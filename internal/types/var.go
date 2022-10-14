package types

import "eim/internal/config"

var (
	MessageSendTopic = config.SystemConfig.LocalIp + "_send"
)
