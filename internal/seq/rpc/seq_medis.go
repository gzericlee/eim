package rpc

import (
	"context"
	"sync"

	"eim/internal/seq/medis"
)

type medisSeq struct {
	mutex    sync.RWMutex
	lock     sync.Mutex
	instance *medis.Instance
}

func (its *medisSeq) Number(ctx context.Context, req *Request, reply *Reply) error {
	its.mutex.Lock()
	var magazine = medis.GetMagazine()
	// 取值时从信道取值，每次判断信道余量，余量小于等于阈值时从 kv 拉取一批值
	if int64(len(magazine.Channel)) <= magazine.Threshold && medis.Freedom == 0 {
		need := magazine.Capacity - int64(len(magazine.Channel))
		// 用全局变量锁定 Goroutine 用完再释放
		medis.Freedom = 1
		go func() {
			err := its.instance.KvToChannel(magazine.Channel, need, magazine.KvThreshold)
			if err != nil {
				return
			}
			medis.Freedom = 0
		}()
	}
	number := <-magazine.Channel
	its.mutex.Unlock()
	reply.Number = number
	return nil
}
