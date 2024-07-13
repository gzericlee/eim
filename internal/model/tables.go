package model

func (x *Biz) TableName() string {
	return "biz"
}

func (x *BizMember) TableName() string {
	return "biz_member"
}

func (x *Device) TableName() string {
	return "device"
}

func (x *Message) TableName() string {
	return "message"
}

func (x *Segment) TableName() string {
	return "segment"
}

func (x *Tenant) TableName() string {
	return "tenant"
}
