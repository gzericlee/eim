package idgenerator

import (
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"

	"eim/util/log"
)

var node *snowflake.Node

func Init(redisEndpoints []string, password string) {
	id, err := registry(config{
		redisEndpoints: redisEndpoints,
		redisPassword:  password,
		maxWorkerId:    1023,
		minWorkerId:    1,
	})
	if err != nil {
		log.Error("IdGenerator registry error", zap.Error(err))
		os.Exit(1)
	}

	log.Info("IdGenerator registry node", zap.Int64("workerId", id))

	node, err = snowflake.NewNode(workerId)
	if err != nil {
		log.Error("IdGenerator new node error", zap.Error(err))
		os.Exit(1)
	}
}

func NextId() int64 {
	var id int64
	var previousTimestamp int64

	for {
		id = node.Generate().Int64()
		if id == 0 {
			log.Warn("IdGenerator generate id is 0")
			continue
		}

		currentTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
		if currentTimestamp < previousTimestamp {
			//TODO 更好的处理方式
			panic("IdGenerator generate id is not monotonic")
		}
		previousTimestamp = currentTimestamp

		break
	}

	return id
}
