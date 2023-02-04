// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: strsvc.proto

package strsvcv1

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*Request_LowerCase
	//	*Request_UpperCase
	Message isRequest_Message `protobuf_oneof:"message"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_strsvc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_strsvc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_strsvc_proto_rawDescGZIP(), []int{0}
}

func (m *Request) GetMessage() isRequest_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Request) GetLowerCase() *LowerCase {
	if x, ok := x.GetMessage().(*Request_LowerCase); ok {
		return x.LowerCase
	}
	return nil
}

func (x *Request) GetUpperCase() *UpperCase {
	if x, ok := x.GetMessage().(*Request_UpperCase); ok {
		return x.UpperCase
	}
	return nil
}

type isRequest_Message interface {
	isRequest_Message()
}

type Request_LowerCase struct {
	LowerCase *LowerCase `protobuf:"bytes,1,opt,name=lower_case,json=lowerCase,proto3,oneof"`
}

type Request_UpperCase struct {
	UpperCase *UpperCase `protobuf:"bytes,2,opt,name=upper_case,json=upperCase,proto3,oneof"`
}

func (*Request_LowerCase) isRequest_Message() {}

func (*Request_UpperCase) isRequest_Message() {}

type LowerCase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Input string `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
}

func (x *LowerCase) Reset() {
	*x = LowerCase{}
	if protoimpl.UnsafeEnabled {
		mi := &file_strsvc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LowerCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LowerCase) ProtoMessage() {}

func (x *LowerCase) ProtoReflect() protoreflect.Message {
	mi := &file_strsvc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LowerCase.ProtoReflect.Descriptor instead.
func (*LowerCase) Descriptor() ([]byte, []int) {
	return file_strsvc_proto_rawDescGZIP(), []int{1}
}

func (x *LowerCase) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

type UpperCase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Input string `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
}

func (x *UpperCase) Reset() {
	*x = UpperCase{}
	if protoimpl.UnsafeEnabled {
		mi := &file_strsvc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpperCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpperCase) ProtoMessage() {}

func (x *UpperCase) ProtoReflect() protoreflect.Message {
	mi := &file_strsvc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpperCase.ProtoReflect.Descriptor instead.
func (*UpperCase) Descriptor() ([]byte, []int) {
	return file_strsvc_proto_rawDescGZIP(), []int{2}
}

func (x *UpperCase) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_strsvc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_strsvc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_strsvc_proto_rawDescGZIP(), []int{3}
}

func (x *Response) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_strsvc_proto protoreflect.FileDescriptor

var file_strsvc_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x74, 0x72, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x73, 0x74, 0x72, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x22, 0x82, 0x01, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0a, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x63,
	0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x74, 0x72, 0x73,
	0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x77, 0x65, 0x72, 0x43, 0x61, 0x73, 0x65, 0x48,
	0x00, 0x52, 0x09, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x43, 0x61, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0a,
	0x75, 0x70, 0x70, 0x65, 0x72, 0x5f, 0x63, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x73, 0x74, 0x72, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x70,
	0x65, 0x72, 0x43, 0x61, 0x73, 0x65, 0x48, 0x00, 0x52, 0x09, 0x75, 0x70, 0x70, 0x65, 0x72, 0x43,
	0x61, 0x73, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x21,
	0x0a, 0x09, 0x4c, 0x6f, 0x77, 0x65, 0x72, 0x43, 0x61, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x22, 0x21, 0x0a, 0x09, 0x55, 0x70, 0x70, 0x65, 0x72, 0x43, 0x61, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x22, 0x22, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x73, 0x74,
	0x72, 0x73, 0x76, 0x63, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_strsvc_proto_rawDescOnce sync.Once
	file_strsvc_proto_rawDescData = file_strsvc_proto_rawDesc
)

func file_strsvc_proto_rawDescGZIP() []byte {
	file_strsvc_proto_rawDescOnce.Do(func() {
		file_strsvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_strsvc_proto_rawDescData)
	})
	return file_strsvc_proto_rawDescData
}

var file_strsvc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_strsvc_proto_goTypes = []interface{}{
	(*Request)(nil),   // 0: strsvc.v1.Request
	(*LowerCase)(nil), // 1: strsvc.v1.LowerCase
	(*UpperCase)(nil), // 2: strsvc.v1.UpperCase
	(*Response)(nil),  // 3: strsvc.v1.Response
}
var file_strsvc_proto_depIdxs = []int32{
	1, // 0: strsvc.v1.Request.lower_case:type_name -> strsvc.v1.LowerCase
	2, // 1: strsvc.v1.Request.upper_case:type_name -> strsvc.v1.UpperCase
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_strsvc_proto_init() }
func file_strsvc_proto_init() {
	if File_strsvc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_strsvc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_strsvc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LowerCase); i {
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
		file_strsvc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpperCase); i {
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
		file_strsvc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
	file_strsvc_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Request_LowerCase)(nil),
		(*Request_UpperCase)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_strsvc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_strsvc_proto_goTypes,
		DependencyIndexes: file_strsvc_proto_depIdxs,
		MessageInfos:      file_strsvc_proto_msgTypes,
	}.Build()
	File_strsvc_proto = out.File
	file_strsvc_proto_rawDesc = nil
	file_strsvc_proto_goTypes = nil
	file_strsvc_proto_depIdxs = nil
}
