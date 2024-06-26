// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v5.27.1
// source: internal/model/device.proto

package model

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId      string                 `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty" bson:"device_id"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" bson:"user_id"`
	DeviceType    string                 `protobuf:"bytes,3,opt,name=device_type,json=deviceType,proto3" json:"device_type,omitempty" bson:"device_type"`
	DeviceVersion string                 `protobuf:"bytes,4,opt,name=device_version,json=deviceVersion,proto3" json:"device_version,omitempty" bson:"device_version"`
	GatewayIp     string                 `protobuf:"bytes,5,opt,name=gateway_ip,json=gatewayIp,proto3" json:"gateway_ip,omitempty" bson:"gateway_ip"`
	OnlineAt      *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=online_at,json=onlineAt,proto3" json:"online_at,omitempty" bson:"online_at"`
	OfflineAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=offline_at,json=offlineAt,proto3" json:"offline_at,omitempty" bson:"offline_at"`
	State         int32                  `protobuf:"varint,8,opt,name=state,proto3" json:"state,omitempty" bson:"state"`
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_model_device_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_internal_model_device_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_internal_model_device_proto_rawDescGZIP(), []int{0}
}

func (x *Device) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *Device) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Device) GetDeviceType() string {
	if x != nil {
		return x.DeviceType
	}
	return ""
}

func (x *Device) GetDeviceVersion() string {
	if x != nil {
		return x.DeviceVersion
	}
	return ""
}

func (x *Device) GetGatewayIp() string {
	if x != nil {
		return x.GatewayIp
	}
	return ""
}

func (x *Device) GetOnlineAt() *timestamppb.Timestamp {
	if x != nil {
		return x.OnlineAt
	}
	return nil
}

func (x *Device) GetOfflineAt() *timestamppb.Timestamp {
	if x != nil {
		return x.OfflineAt
	}
	return nil
}

func (x *Device) GetState() int32 {
	if x != nil {
		return x.State
	}
	return 0
}

var File_internal_model_device_proto protoreflect.FileDescriptor

var file_internal_model_device_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaf,
	0x02, 0x0a, 0x06, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x25, 0x0a, 0x0e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x5f, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x49, 0x70, 0x12, 0x37, 0x0a, 0x09, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65,
	0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_model_device_proto_rawDescOnce sync.Once
	file_internal_model_device_proto_rawDescData = file_internal_model_device_proto_rawDesc
)

func file_internal_model_device_proto_rawDescGZIP() []byte {
	file_internal_model_device_proto_rawDescOnce.Do(func() {
		file_internal_model_device_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_model_device_proto_rawDescData)
	})
	return file_internal_model_device_proto_rawDescData
}

var file_internal_model_device_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_model_device_proto_goTypes = []interface{}{
	(*Device)(nil),                // 0: Device
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_internal_model_device_proto_depIdxs = []int32{
	1, // 0: Device.online_at:type_name -> google.protobuf.Timestamp
	1, // 1: Device.offline_at:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_model_device_proto_init() }
func file_internal_model_device_proto_init() {
	if File_internal_model_device_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_model_device_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
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
			RawDescriptor: file_internal_model_device_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_model_device_proto_goTypes,
		DependencyIndexes: file_internal_model_device_proto_depIdxs,
		MessageInfos:      file_internal_model_device_proto_msgTypes,
	}.Build()
	File_internal_model_device_proto = out.File
	file_internal_model_device_proto_rawDesc = nil
	file_internal_model_device_proto_goTypes = nil
	file_internal_model_device_proto_depIdxs = nil
}
