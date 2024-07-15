package model

import "github.com/gzericlee/eim/internal/model"

type AuthArgs struct {
	Token string
}

type AuthReply struct {
	Biz *model.Biz
}
