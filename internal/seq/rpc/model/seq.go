package model

type Request struct {
	BizId    string
	TenantId string
}

type Reply struct {
	Number int64
}
