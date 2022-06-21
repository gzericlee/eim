package redis

import (
	"time"

	"eim/model"
)

func GetSegmentSeq(id string) (*model.Seq, error) {
	key := id + ":seq"
	value, _ := rdsClient.Get(key)
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
	err := rdsClient.Set(key, body, 0)
	return seq, err
}
