// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v5.27.2
// source: internal/model/biz_member.proto

package model

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BizMember struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BizId          string `protobuf:"bytes,1,opt,name=biz_id,json=bizId,proto3" json:"biz_id,omitempty" bson:"biz_id"`
	MemberId       string `protobuf:"bytes,2,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty" bson:"member_id"`
	BizTenantId    string `protobuf:"bytes,3,opt,name=biz_tenant_id,json=bizTenantId,proto3" json:"biz_tenant_id,omitempty" bson:"biz_tenant_id"`
	MemberTenantId string `protobuf:"bytes,4,opt,name=member_tenant_id,json=memberTenantId,proto3" json:"member_tenant_id,omitempty" bson:"member_tenant_id"`
}

func (x *BizMember) Reset() {
	*x = BizMember{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_model_biz_member_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BizMember) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BizMember) ProtoMessage() {}

func (x *BizMember) ProtoReflect() protoreflect.Message {
	mi := &file_internal_model_biz_member_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BizMember.ProtoReflect.Descriptor instead.
func (*BizMember) Descriptor() ([]byte, []int) {
	return file_internal_model_biz_member_proto_rawDescGZIP(), []int{0}
}

func (x *BizMember) GetBizId() string {
	if x != nil {
		return x.BizId
	}
	return ""
}

func (x *BizMember) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

func (x *BizMember) GetBizTenantId() string {
	if x != nil {
		return x.BizTenantId
	}
	return ""
}

func (x *BizMember) GetMemberTenantId() string {
	if x != nil {
		return x.MemberTenantId
	}
	return ""
}

var File_internal_model_biz_member_proto protoreflect.FileDescriptor

var file_internal_model_biz_member_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2f, 0x62, 0x69, 0x7a, 0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0x8d, 0x01, 0x0a, 0x09, 0x42, 0x69, 0x7a,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x15, 0x0a, 0x06, 0x62, 0x69, 0x7a, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62, 0x69, 0x7a, 0x49, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x62, 0x69,
	0x7a, 0x5f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x62, 0x69, 0x7a, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x28,
	0x0a, 0x10, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_model_biz_member_proto_rawDescOnce sync.Once
	file_internal_model_biz_member_proto_rawDescData = file_internal_model_biz_member_proto_rawDesc
)

func file_internal_model_biz_member_proto_rawDescGZIP() []byte {
	file_internal_model_biz_member_proto_rawDescOnce.Do(func() {
		file_internal_model_biz_member_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_model_biz_member_proto_rawDescData)
	})
	return file_internal_model_biz_member_proto_rawDescData
}

var file_internal_model_biz_member_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_model_biz_member_proto_goTypes = []interface{}{
	(*BizMember)(nil), // 0: model.BizMember
}
var file_internal_model_biz_member_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_model_biz_member_proto_init() }
func file_internal_model_biz_member_proto_init() {
	if File_internal_model_biz_member_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_model_biz_member_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BizMember); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_model_biz_member_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_model_biz_member_proto_goTypes,
		DependencyIndexes: file_internal_model_biz_member_proto_depIdxs,
		MessageInfos:      file_internal_model_biz_member_proto_msgTypes,
	}.Build()
	File_internal_model_biz_member_proto = out.File
	file_internal_model_biz_member_proto_rawDesc = nil
	file_internal_model_biz_member_proto_goTypes = nil
	file_internal_model_biz_member_proto_depIdxs = nil
}
