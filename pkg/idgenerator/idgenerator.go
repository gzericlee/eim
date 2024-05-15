package idgenerator

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"

	"eim/util/log"
)

func Init(redisEndpoints []string, password string) {
	id := registry(config{
		redisEndpoints: redisEndpoints,
		redisPassword:  password,
		maxWorkerId:    1023,
		minWorkerId:    1,
	})
	log.Info("IdGenerator worker id", zap.Any("id", id))

	var options = idgen.NewIdGeneratorOptions(uint16(id))
	options.WorkerIdBitLength = 10
	options.SeqBitLength = 10
	idgen.SetIdGenerator(options)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		unRegistry()
		os.Exit(0)
	}()
}

func NextId() int64 {
	return idgen.NextId()
}
