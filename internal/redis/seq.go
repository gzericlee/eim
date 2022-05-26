package redis

import (
	"time"

	"go.uber.org/zap"

	"eim/global"
	"eim/model"
)

func GetIncrSeq(userId string) int64 {
	key := userId + ":seq"
	id, err := rdsClient.Incr(key)
	if err != nil {
		global.Logger.Error("Error incr seq id", zap.String("key", key), zap.Error(err))
	}
	return id
}

func GetSegmentSeq(userId string) (*model.Seq, error) {
	key := userId + ":seq"
	value, _ := rdsClient.Get(key)
	seq := &model.Seq{}

	if value == "" {
		seq.UserId = userId
		seq.MaxId = 0
		seq.Step = 1000
		seq.CreateAt = time.Now().Local()
		seq.UpdateAt = time.Now().Local()
	} else {
		err := seq.Deserialize([]byte(value))
		if err != nil {
			return seq, err
		}
		seq.MaxId = seq.MaxId + int64(seq.Step)
		seq.UpdateAt = time.Now().Local()
	}

	body, _ := seq.Serialize()
	err := rdsClient.Set(key, body, 0)
	return seq, err
}
