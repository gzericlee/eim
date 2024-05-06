package idgenerator

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"eim/pkg/log"
)

var ctx = context.Background()
var workerIdLock sync.Mutex
var redisClient redis.UniversalClient

var workerId int32
var loopCount int32 = 0
var lifeIndex int32 = -1

var maxLoopCount int32 = 20
var workerIdLifeTimeSeconds int32 = 15
var sleepMillisecondEveryLoop int32 = 200
var maxWorkerId int32 = 0
var minWorkerId int32 = 0
var workerIdFlag = "Y"

var workerIdIndexKey string = "id_gen:worker_id:index"
var workerIdValueKeyPrefix string = "id_gen:worker_id:value"

type config struct {
	redisEndpoints  []string
	redisPassword   string
	maxWorkerId     int32
	minWorkerId     int32
	lifeTimeSeconds int32
}

func unRegistry() {
	if redisClient == nil {
		return
	}

	workerIdLock.Lock()
	defer workerIdLock.Unlock()

	lifeIndex = -1
	if workerId > -1 {
		redisClient.Del(ctx, fmt.Sprintf("%s:%d", workerIdValueKeyPrefix, workerId))
	}

	workerId = -1
}

func registry(cfg config) int32 {
	if cfg.maxWorkerId < 0 || cfg.minWorkerId > cfg.maxWorkerId {
		return -2
	}

	redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.redisEndpoints,
		Password:     cfg.redisPassword,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		MaxRedirects: 8,
		PoolTimeout:  30 * time.Second,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Error("redis ping failed", zap.Error(err))
		return -1
	}

	maxWorkerId = cfg.maxWorkerId
	minWorkerId = cfg.minWorkerId
	workerIdLifeTimeSeconds = cfg.lifeTimeSeconds
	loopCount = 0

	unRegistry()

	lifeIndex++
	id := getNextWorkerId(lifeIndex)
	if id > -1 {
		workerId = id
		go extendLifeTime(lifeIndex)
	}

	return id
}

func getNextWorkerId(lifeTime int32) int32 {
	// 获取当前 WorkerIdIndex
	r, err := redisClient.Incr(ctx, workerIdIndexKey).Result()
	if err != nil {
		return -1
	}

	candidateId := int32(r)

	// 设置最小值
	if candidateId < minWorkerId {
		candidateId = minWorkerId
		setWorkerIdIndex(minWorkerId)
	}

	// 如果 candidateId 大于最大值，则重置
	if candidateId > maxWorkerId {
		if canReset() {
			// 当前应用获得重置 WorkerIdIndex 的权限
			//setWorkerIdIndex(-1)
			setWorkerIdIndex(minWorkerId - 1)
			endReset() // 此步有可能不被执行？
			loopCount++

			// 超过一定次数，直接终止操作
			if loopCount > maxLoopCount {
				loopCount = 0

				// 返回错误
				return -1
			}

			// 每次一个大循环后，暂停一些时间
			time.Sleep(time.Duration(sleepMillisecondEveryLoop*loopCount) * time.Millisecond)

			return getNextWorkerId(lifeTime)
		} else {
			// 如果有其它应用正在编辑，则本应用暂停200ms后，再继续
			time.Sleep(time.Duration(200) * time.Millisecond)

			return getNextWorkerId(lifeTime)
		}
	}

	if isAvailable(candidateId) {
		// 最新获得的 WorkerIdIndex，在 redis 中是可用状态
		setWorkerIdFlag(candidateId)
		loopCount = 0

		// 获取到可用 WorkerId 后，启用新线程，每隔 1/3个 _WorkerIdLifeTimeSeconds 时间，向服务器续期（延长一次 LifeTime）
		// go extendWorkerIdLifeTime(lifeTime, candidateId)

		return candidateId
	} else {
		// 最新获得的 WorkerIdIndex，在 redis 中是不可用状态，则继续下一个 WorkerIdIndex
		return getNextWorkerId(lifeTime)
	}
}

func extendLifeTime(myLifeIndex int32) {
	// 获取到可用 WorkerId 后，启用新线程，每隔 1/3个 _WorkerIdLifeTimeSeconds 时间，向服务器续期（延长一次 LifeTime）

	// 循环操作：间隔一定时间，刷新 WorkerId 在 redis 中的有效时间。
	for {
		time.Sleep(time.Duration(workerIdLifeTimeSeconds/3) * time.Second)

		// 上锁操作，防止跟 UnRegister 操作重叠
		workerIdLock.Lock()

		// 如果临时变量 myLifeIndex 不等于 全局变量 _lifeIndex，表明全局状态被修改，当前线程可终止，不应继续操作 redis
		// 还应主动释放 redis 键值缓存
		if myLifeIndex != lifeIndex {
			break
		}

		// 已经被注销，则终止（此步是上一步的二次验证）
		if workerId == -1 {
			break
		}

		// 延长 redis 数据有效期
		if workerId > -1 {
			extendWorkerIdFlag(workerId)
		}

		workerIdLock.Unlock()
	}
}

func extendWorkerIdLifeTime(myLifeIndex int32, workerId int32) {
	// 循环操作：间隔一定时间，刷新 WorkerId 在 redis 中的有效时间。
	for {
		time.Sleep(time.Duration(workerIdLifeTimeSeconds/3) * time.Second)

		workerIdLock.Lock()

		if myLifeIndex != lifeIndex {
			break
		}

		extendWorkerIdFlag(workerId)

		workerIdLock.Unlock()
	}
}

func setWorkerIdIndex(val int32) {
	redisClient.Set(ctx, workerIdIndexKey, val, 0)
}

func setWorkerIdFlag(workerId int32) {
	redisClient.Set(ctx, fmt.Sprintf("%s:%d", workerIdValueKeyPrefix, workerId), workerIdFlag, time.Duration(workerIdLifeTimeSeconds)*time.Second)
}

func extendWorkerIdFlag(workerId int32) {
	redisClient.Expire(ctx, workerIdValueKeyPrefix+strconv.Itoa(int(workerId)), time.Duration(workerIdLifeTimeSeconds)*time.Second)
}

func canReset() bool {
	r, err := redisClient.Incr(ctx, fmt.Sprintf("%s:%s", workerIdValueKeyPrefix, "edit")).Result()
	if err != nil {
		return false
	}
	return r != 1
}

func endReset() {
	redisClient.Set(ctx, fmt.Sprintf("%s:%s", workerIdValueKeyPrefix, "edit"), 0, 0)
}

func getWorkerIdFlag(workerId int32) (string, bool) {
	r, err := redisClient.Get(ctx, fmt.Sprintf("%s:%d", workerIdValueKeyPrefix, workerId)).Result()
	if err != nil {
		return "", false
	}
	return r, true
}

func isAvailable(workerId int32) bool {
	r, err := redisClient.Get(ctx, fmt.Sprintf("%s:%d", workerIdValueKeyPrefix, workerId)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return true
		}
		return false
	}
	return r != workerIdFlag
}
