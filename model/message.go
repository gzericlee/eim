package model

import "eim/pkg/json"

type Message struct {
	MsgId      string `gorm:"primary_key;column:msg_id" json:"msgId"`
	SeqId      int64  `gorm:"column:seq_id" json:"seqId"`
	MsgType    int64  `gorm:"column:msg_type" json:"msgType"`
	Content    string `gorm:"column:content" json:"content"`
	FromType   int64  `gorm:"column:from_type" json:"fromType"`
	FromId     string `gorm:"column:from_id" json:"fromId"`
	FromDevice string `gorm:"column:from_device" json:"fromDevice"`
	ToType     int64  `gorm:"column:to_type" json:"toType"`
	ToId       string `gorm:"column:to_id" json:"toId"`
	ToDevice   string `gorm:"column:to_device" json:"toDevice"`
	SendTime   int64  `gorm:"column:send_time" json:"sendTime"`
	UserId     string `gorm:"-" json:"userId"`
}

func (its *Message) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *Message) Serialize() ([]byte, error) {
	return json.Marshal(its)
}

func (its *Message) TableName() string {
	return "eim_message"
}
