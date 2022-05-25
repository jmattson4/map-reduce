// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: grpc/coordinator.proto

package grpc

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

type TaskArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	HandlerType int32  `protobuf:"varint,2,opt,name=HandlerType,proto3" json:"HandlerType,omitempty"`
}

func (x *TaskArgs) Reset() {
	*x = TaskArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_coordinator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskArgs) ProtoMessage() {}

func (x *TaskArgs) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_coordinator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskArgs.ProtoReflect.Descriptor instead.
func (*TaskArgs) Descriptor() ([]byte, []int) {
	return file_grpc_coordinator_proto_rawDescGZIP(), []int{0}
}

func (x *TaskArgs) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskArgs) GetHandlerType() int32 {
	if x != nil {
		return x.HandlerType
	}
	return 0
}

type TaskReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Type    int32  `protobuf:"varint,2,opt,name=Type,proto3" json:"Type,omitempty"`
	NReduce int32  `protobuf:"varint,3,opt,name=NReduce,proto3" json:"NReduce,omitempty"`
}

func (x *TaskReply) Reset() {
	*x = TaskReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_coordinator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskReply) ProtoMessage() {}

func (x *TaskReply) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_coordinator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskReply.ProtoReflect.Descriptor instead.
func (*TaskReply) Descriptor() ([]byte, []int) {
	return file_grpc_coordinator_proto_rawDescGZIP(), []int{1}
}

func (x *TaskReply) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskReply) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *TaskReply) GetNReduce() int32 {
	if x != nil {
		return x.NReduce
	}
	return 0
}

var File_grpc_coordinator_proto protoreflect.FileDescriptor

var file_grpc_coordinator_proto_rawDesc = []byte{
	0x0a, 0x16, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b,
	0x41, 0x72, 0x67, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x48, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x22, 0x49, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4e, 0x52, 0x65, 0x64, 0x75,
	0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x4e, 0x52, 0x65, 0x64, 0x75, 0x63,
	0x65, 0x32, 0x31, 0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x22, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x09, 0x2e, 0x54, 0x61,
	0x73, 0x6b, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x0a, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x00, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6a, 0x6d, 0x61, 0x74, 0x74, 0x73, 0x6f, 0x6e, 0x34, 0x2f, 0x6d, 0x61, 0x70,
	0x2d, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_coordinator_proto_rawDescOnce sync.Once
	file_grpc_coordinator_proto_rawDescData = file_grpc_coordinator_proto_rawDesc
)

func file_grpc_coordinator_proto_rawDescGZIP() []byte {
	file_grpc_coordinator_proto_rawDescOnce.Do(func() {
		file_grpc_coordinator_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_coordinator_proto_rawDescData)
	})
	return file_grpc_coordinator_proto_rawDescData
}

var file_grpc_coordinator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_grpc_coordinator_proto_goTypes = []interface{}{
	(*TaskArgs)(nil),  // 0: TaskArgs
	(*TaskReply)(nil), // 1: TaskReply
}
var file_grpc_coordinator_proto_depIdxs = []int32{
	0, // 0: TaskService.GetTask:input_type -> TaskArgs
	1, // 1: TaskService.GetTask:output_type -> TaskReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpc_coordinator_proto_init() }
func file_grpc_coordinator_proto_init() {
	if File_grpc_coordinator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_coordinator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskArgs); i {
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
		file_grpc_coordinator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskReply); i {
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
			RawDescriptor: file_grpc_coordinator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_coordinator_proto_goTypes,
		DependencyIndexes: file_grpc_coordinator_proto_depIdxs,
		MessageInfos:      file_grpc_coordinator_proto_msgTypes,
	}.Build()
	File_grpc_coordinator_proto = out.File
	file_grpc_coordinator_proto_rawDesc = nil
	file_grpc_coordinator_proto_goTypes = nil
	file_grpc_coordinator_proto_depIdxs = nil
}