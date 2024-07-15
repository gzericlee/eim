package model

import "github.com/gzericlee/eim/internal/model"

type BizMemberArgs struct {
	BizMember *model.BizMember
}

type BizMembersReply struct {
	Members []string
}
