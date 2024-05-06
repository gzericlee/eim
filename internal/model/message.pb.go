// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v4.25.3
// source: internal/model/message.proto

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgId      int64  `protobuf:"varint,1,opt,name=msgId,proto3" json:"msgId,omitempty"`
	SeqId      int64  `protobuf:"varint,2,opt,name=seqId,proto3" json:"seqId,omitempty"`
	MsgType    int64  `protobuf:"varint,3,opt,name=msgType,proto3" json:"msgType,omitempty"`
	Content    string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	FromType   int64  `protobuf:"varint,5,opt,name=fromType,proto3" json:"fromType,omitempty"`
	FromId     string `protobuf:"bytes,6,opt,name=fromId,proto3" json:"fromId,omitempty"`
	FromDevice string `protobuf:"bytes,7,opt,name=fromDevice,proto3" json:"fromDevice,omitempty"`
	ToType     int64  `protobuf:"varint,8,opt,name=toType,proto3" json:"toType,omitempty"`
	ToId       string `protobuf:"bytes,9,opt,name=toId,proto3" json:"toId,omitempty"`
	ToDevice   string `protobuf:"bytes,10,opt,name=toDevice,proto3" json:"toDevice,omitempty"`
	SendTime   int64  `protobuf:"varint,11,opt,name=sendTime,proto3" json:"sendTime,omitempty"`
	UserId     string `protobuf:"bytes,12,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_model_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_internal_model_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_internal_model_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetMsgId() int64 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

func (x *Message) GetSeqId() int64 {
	if x != nil {
		return x.SeqId
	}
	return 0
}

func (x *Message) GetMsgType() int64 {
	if x != nil {
		return x.MsgType
	}
	return 0
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetFromType() int64 {
	if x != nil {
		return x.FromType
	}
	return 0
}

func (x *Message) GetFromId() string {
	if x != nil {
		return x.FromId
	}
	return ""
}

func (x *Message) GetFromDevice() string {
	if x != nil {
		return x.FromDevice
	}
	return ""
}

func (x *Message) GetToType() int64 {
	if x != nil {
		return x.ToType
	}
	return 0
}

func (x *Message) GetToId() string {
	if x != nil {
		return x.ToId
	}
	return ""
}

func (x *Message) GetToDevice() string {
	if x != nil {
		return x.ToDevice
	}
	return ""
}

func (x *Message) GetSendTime() int64 {
	if x != nil {
		return x.SendTime
	}
	return 0
}

func (x *Message) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

var File_internal_model_message_proto protoreflect.FileDescriptor

var file_internal_model_message_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0xb9, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x65, 0x71, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x65, 0x71, 0x49, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x66, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66,
	0x72, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x6f, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x6f, 0x49, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x6f, 0x49,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x6f, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x6f, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_model_message_proto_rawDescOnce sync.Once
	file_internal_model_message_proto_rawDescData = file_internal_model_message_proto_rawDesc
)

func file_internal_model_message_proto_rawDescGZIP() []byte {
	file_internal_model_message_proto_rawDescOnce.Do(func() {
		file_internal_model_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_model_message_proto_rawDescData)
	})
	return file_internal_model_message_proto_rawDescData
}

var file_internal_model_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_model_message_proto_goTypes = []interface{}{
	(*Message)(nil), // 0: model.Message
}
var file_internal_model_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_model_message_proto_init() }
func file_internal_model_message_proto_init() {
	if File_internal_model_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_model_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_internal_model_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_model_message_proto_goTypes,
		DependencyIndexes: file_internal_model_message_proto_depIdxs,
		MessageInfos:      file_internal_model_message_proto_msgTypes,
	}.Build()
	File_internal_model_message_proto = out.File
	file_internal_model_message_proto_rawDesc = nil
	file_internal_model_message_proto_goTypes = nil
	file_internal_model_message_proto_depIdxs = nil
}