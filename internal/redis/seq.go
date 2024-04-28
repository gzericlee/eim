package redis

import (
	"time"

	"eim/internal/model"
)

func (its *Manager) Incr(key string) (int64, error) {
	return its.rdsClient.Incr(key)
}

func (its *Manager) Decr(key string) (int64, error) {
	return its.rdsClient.Decr(key)
}

func (its *Manager) GetSegmentSeq(id string) (*model.Seq, error) {
	key := id + ":seq"
	value, _ := its.rdsClient.Get(key)
	seq := &model.Seq{}

	if value == "" {
		seq.Id = id
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
	err := its.rdsClient.Set(key, body, 0)
	return seq, err
}
