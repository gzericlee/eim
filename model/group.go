package model

type Group struct {
	GroupId   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type GroupMember struct {
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
	SeqId   int64  `json:"seqId"`
}
