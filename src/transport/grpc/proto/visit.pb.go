// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: src/transport/grpc/proto/visit.proto

package pb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type VisitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        uint32 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	FirstName string `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName,omitempty"`
	LastName  string `protobuf:"bytes,3,opt,name=LastName,proto3" json:"LastName,omitempty"`
}

func (x *VisitRequest) Reset() {
	*x = VisitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_transport_grpc_proto_visit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VisitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VisitRequest) ProtoMessage() {}

func (x *VisitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_src_transport_grpc_proto_visit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VisitRequest.ProtoReflect.Descriptor instead.
func (*VisitRequest) Descriptor() ([]byte, []int) {
	return file_src_transport_grpc_proto_visit_proto_rawDescGZIP(), []int{0}
}

func (x *VisitRequest) GetID() uint32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *VisitRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *VisitRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

type VisitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        uint32                 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	FirstName string                 `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName,omitempty"`
	LastName  string                 `protobuf:"bytes,3,opt,name=LastName,proto3" json:"LastName,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
}

func (x *VisitResponse) Reset() {
	*x = VisitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_transport_grpc_proto_visit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VisitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VisitResponse) ProtoMessage() {}

func (x *VisitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_src_transport_grpc_proto_visit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VisitResponse.ProtoReflect.Descriptor instead.
func (*VisitResponse) Descriptor() ([]byte, []int) {
	return file_src_transport_grpc_proto_visit_proto_rawDescGZIP(), []int{1}
}

func (x *VisitResponse) GetID() uint32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *VisitResponse) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *VisitResponse) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *VisitResponse) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *VisitResponse) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_src_transport_grpc_proto_visit_proto protoreflect.FileDescriptor

var file_src_transport_grpc_proto_visit_proto_rawDesc = []byte{
	0x0a, 0x24, 0x73, 0x72, 0x63, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x73, 0x72,
	0x63, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x0c, 0x56, 0x69, 0x73, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x28, 0x00, 0x52, 0x02, 0x49, 0x44, 0x12, 0x25,
	0x0a, 0x09, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x02, 0x52, 0x09, 0x46, 0x69, 0x72, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x08, 0x4c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x02,
	0x52, 0x08, 0x4c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xcd, 0x01, 0x0a, 0x0d, 0x56,
	0x69, 0x73, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09,
	0x46, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x61,
	0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4c, 0x61,
	0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x38, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x32, 0x93, 0x01, 0x0a, 0x05, 0x56,
	0x69, 0x73, 0x69, 0x74, 0x12, 0x38, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x06, 0x2e, 0x70, 0x62,
	0x2e, 0x49, 0x44, 0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x69, 0x73, 0x69, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e,
	0x2f, 0x76, 0x31, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x74, 0x2f, 0x7b, 0x49, 0x44, 0x7d, 0x12, 0x50,
	0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x69, 0x73, 0x69, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x69, 0x73,
	0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x1e, 0x22, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x74, 0x3a, 0x01, 0x2a,
	0x5a, 0x0e, 0x1a, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x74, 0x3a, 0x01, 0x2a,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_src_transport_grpc_proto_visit_proto_rawDescOnce sync.Once
	file_src_transport_grpc_proto_visit_proto_rawDescData = file_src_transport_grpc_proto_visit_proto_rawDesc
)

func file_src_transport_grpc_proto_visit_proto_rawDescGZIP() []byte {
	file_src_transport_grpc_proto_visit_proto_rawDescOnce.Do(func() {
		file_src_transport_grpc_proto_visit_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_transport_grpc_proto_visit_proto_rawDescData)
	})
	return file_src_transport_grpc_proto_visit_proto_rawDescData
}

var file_src_transport_grpc_proto_visit_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_src_transport_grpc_proto_visit_proto_goTypes = []interface{}{
	(*VisitRequest)(nil),          // 0: pb.VisitRequest
	(*VisitResponse)(nil),         // 1: pb.VisitResponse
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
	(*ID)(nil),                    // 3: pb.ID
}
var file_src_transport_grpc_proto_visit_proto_depIdxs = []int32{
	2, // 0: pb.VisitResponse.CreatedAt:type_name -> google.protobuf.Timestamp
	2, // 1: pb.VisitResponse.UpdatedAt:type_name -> google.protobuf.Timestamp
	3, // 2: pb.Visit.Get:input_type -> pb.ID
	0, // 3: pb.Visit.Set:input_type -> pb.VisitRequest
	1, // 4: pb.Visit.Get:output_type -> pb.VisitResponse
	1, // 5: pb.Visit.Set:output_type -> pb.VisitResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_src_transport_grpc_proto_visit_proto_init() }
func file_src_transport_grpc_proto_visit_proto_init() {
	if File_src_transport_grpc_proto_visit_proto != nil {
		return
	}
	file_src_transport_grpc_proto_generics_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_src_transport_grpc_proto_visit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VisitRequest); i {
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
		file_src_transport_grpc_proto_visit_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VisitResponse); i {
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
			RawDescriptor: file_src_transport_grpc_proto_visit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_src_transport_grpc_proto_visit_proto_goTypes,
		DependencyIndexes: file_src_transport_grpc_proto_visit_proto_depIdxs,
		MessageInfos:      file_src_transport_grpc_proto_visit_proto_msgTypes,
	}.Build()
	File_src_transport_grpc_proto_visit_proto = out.File
	file_src_transport_grpc_proto_visit_proto_rawDesc = nil
	file_src_transport_grpc_proto_visit_proto_goTypes = nil
	file_src_transport_grpc_proto_visit_proto_depIdxs = nil
}