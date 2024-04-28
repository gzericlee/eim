package model

import "eim/pkg/json"

type Message struct {
	MsgId      string `json:"msgId" bson:"msgId"`
	SeqId      int64  `json:"seqId" bson:"seqId"`
	MsgType    int64  `json:"msgType" bson:"msgType"`
	Content    string `json:"content" bson:"content"`
	FromType   int64  `json:"fromType" bson:"fromType"`
	FromId     string `json:"fromId" bson:"fromId"`
	FromDevice string `json:"fromDevice" bson:"fromDevice"`
	ToType     int64  `json:"toType" bson:"toType"`
	ToId       string `json:"toId" bson:"toId"`
	ToDevice   string `json:"toDevice" bson:"toDevice"`
	SendTime   int64  `json:"sendTime" bson:"sendTime"`
	UserId     string `json:"-" bson:"-"`
}

func (its *Message) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *Message) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
