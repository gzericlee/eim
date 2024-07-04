package snowflake

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"

	"eim/pkg/log"
)

type GeneratorConfig struct {
	RedisEndpoints []string
	RedisPassword  string
	MaxWorkerId    int64
	MinWorkerId    int64
	// NodeCount is the number of nodes that generate ids
	// NodeCount = 1 monotonic increasing, > 1 trend increasing
	NodeCount int
}

type Generator struct {
	idChan chan int64
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	generator := &Generator{idChan: make(chan int64, cfg.NodeCount)}

	workIdManager, err := wewWorkerIdManager(workerIdConfig{
		redisEndpoints: cfg.RedisEndpoints,
		redisPassword:  cfg.RedisPassword,
		maxWorkerId:    cfg.MaxWorkerId,
		minWorkerId:    cfg.MinWorkerId,
	})
	if err != nil {
		return nil, fmt.Errorf("new worker id manager -> %w", err)
	}

	for i := 0; i < cfg.NodeCount; i++ {
		workerId, err := workIdManager.nextId()
		if err != nil {
			return nil, fmt.Errorf("generate worker id -> %w", err)
		}

		node, err := snowflake.NewNode(workerId)
		if err != nil {
			return nil, fmt.Errorf("new snowflake node -> %w", err)
		}

		log.Info("New snowflake node", zap.Int64("workerId", workerId))

		go func(node *snowflake.Node) {
			var previousTimestamp int64
			for {
				id := node.Generate().Int64()
				if id == 0 {
					log.Warn("IdGenerator generate id is 0")
					time.Sleep(time.Millisecond * 200)
					continue
				}

				currentTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
				if currentTimestamp < previousTimestamp {
					log.Error("IdGenerator generate id is not monotonic")
					time.Sleep(time.Millisecond * 200)
					continue
				}
				previousTimestamp = currentTimestamp

				if id != 0 {
					generator.idChan <- id
				}
			}
		}(node)
	}

	return generator, nil
}

func (its *Generator) NextId() int64 {
	return <-its.idChan
}
