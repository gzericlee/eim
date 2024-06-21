// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v4.25.3
// source: internal/model/biz.proto

package model

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Biz_BizType int32

const (
	Biz_UNKNOWN Biz_BizType = 0
	Biz_USER    Biz_BizType = 1
	Biz_GROUP   Biz_BizType = 2
	Biz_SERVICE Biz_BizType = 3
)

// Enum value maps for Biz_BizType.
var (
	Biz_BizType_name = map[int32]string{
		0: "UNKNOWN",
		1: "USER",
		2: "GROUP",
		3: "SERVICE",
	}
	Biz_BizType_value = map[string]int32{
		"UNKNOWN": 0,
		"USER":    1,
		"GROUP":   2,
		"SERVICE": 3,
	}
)

func (x Biz_BizType) Enum() *Biz_BizType {
	p := new(Biz_BizType)
	*p = x
	return p
}

func (x Biz_BizType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Biz_BizType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_model_biz_proto_enumTypes[0].Descriptor()
}

func (Biz_BizType) Type() protoreflect.EnumType {
	return &file_internal_model_biz_proto_enumTypes[0]
}

func (x Biz_BizType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Biz_BizType.Descriptor instead.
func (Biz_BizType) EnumDescriptor() ([]byte, []int) {
	return file_internal_model_biz_proto_rawDescGZIP(), []int{0, 0}
}

type Biz struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BizId      string                `protobuf:"bytes,1,opt,name=biz_id,json=bizId,proto3" json:"biz_id,omitempty" bson:"biz_id"`
	BizType    Biz_BizType           `protobuf:"varint,2,opt,name=biz_type,json=bizType,proto3,enum=model.Biz_BizType" json:"biz_type,omitempty" bson:"biz_type"`
	BizName    string                `protobuf:"bytes,3,opt,name=biz_name,json=bizName,proto3" json:"biz_name,omitempty" bson:"biz_name"`
	TenantId   string                `protobuf:"bytes,4,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty" bson:"tenant_id"`
	TenantName string                `protobuf:"bytes,5,opt,name=tenant_name,json=tenantName,proto3" json:"tenant_name,omitempty" bson:"tenant_name"`
	Attributes map[string]*anypb.Any `protobuf:"bytes,6,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" bson:"attributes"`
}

func (x *Biz) Reset() {
	*x = Biz{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_model_biz_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Biz) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Biz) ProtoMessage() {}

func (x *Biz) ProtoReflect() protoreflect.Message {
	mi := &file_internal_model_biz_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Biz.ProtoReflect.Descriptor instead.
func (*Biz) Descriptor() ([]byte, []int) {
	return file_internal_model_biz_proto_rawDescGZIP(), []int{0}
}

func (x *Biz) GetBizId() string {
	if x != nil {
		return x.BizId
	}
	return ""
}

func (x *Biz) GetBizType() Biz_BizType {
	if x != nil {
		return x.BizType
	}
	return Biz_UNKNOWN
}

func (x *Biz) GetBizName() string {
	if x != nil {
		return x.BizName
	}
	return ""
}

func (x *Biz) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *Biz) GetTenantName() string {
	if x != nil {
		return x.TenantName
	}
	return ""
}

func (x *Biz) GetAttributes() map[string]*anypb.Any {
	if x != nil {
		return x.Attributes
	}
	return nil
}

var File_internal_model_biz_proto protoreflect.FileDescriptor

var file_internal_model_biz_proto_rawDesc = []byte{
	0x0a, 0x18, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2f, 0x62, 0x69, 0x7a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xef, 0x02, 0x0a,
	0x03, 0x42, 0x69, 0x7a, 0x12, 0x15, 0x0a, 0x06, 0x62, 0x69, 0x7a, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62, 0x69, 0x7a, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x08, 0x62,
	0x69, 0x7a, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x42, 0x69, 0x7a, 0x2e, 0x42, 0x69, 0x7a, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x07, 0x62, 0x69, 0x7a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x69,
	0x7a, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x69,
	0x7a, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e,
	0x42, 0x69, 0x7a, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x1a,
	0x53, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x38, 0x0a, 0x07, 0x42, 0x69, 0x7a, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x55, 0x53, 0x45, 0x52, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x52, 0x4f, 0x55, 0x50, 0x10,
	0x02, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x10, 0x03, 0x42, 0x0a,
	0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_internal_model_biz_proto_rawDescOnce sync.Once
	file_internal_model_biz_proto_rawDescData = file_internal_model_biz_proto_rawDesc
)

func file_internal_model_biz_proto_rawDescGZIP() []byte {
	file_internal_model_biz_proto_rawDescOnce.Do(func() {
		file_internal_model_biz_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_model_biz_proto_rawDescData)
	})
	return file_internal_model_biz_proto_rawDescData
}

var file_internal_model_biz_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_model_biz_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_model_biz_proto_goTypes = []interface{}{
	(Biz_BizType)(0),  // 0: model.Biz.BizType
	(*Biz)(nil),       // 1: model.Biz
	nil,               // 2: model.Biz.AttributesEntry
	(*anypb.Any)(nil), // 3: google.protobuf.Any
}
var file_internal_model_biz_proto_depIdxs = []int32{
	0, // 0: model.Biz.biz_type:type_name -> model.Biz.BizType
	2, // 1: model.Biz.attributes:type_name -> model.Biz.AttributesEntry
	3, // 2: model.Biz.AttributesEntry.value:type_name -> google.protobuf.Any
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_model_biz_proto_init() }
func file_internal_model_biz_proto_init() {
	if File_internal_model_biz_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_model_biz_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Biz); i {
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
			RawDescriptor: file_internal_model_biz_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_model_biz_proto_goTypes,
		DependencyIndexes: file_internal_model_biz_proto_depIdxs,
		EnumInfos:         file_internal_model_biz_proto_enumTypes,
		MessageInfos:      file_internal_model_biz_proto_msgTypes,
	}.Build()
	File_internal_model_biz_proto = out.File
	file_internal_model_biz_proto_rawDesc = nil
	file_internal_model_biz_proto_goTypes = nil
	file_internal_model_biz_proto_depIdxs = nil
}
