// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.26.1
// source: proto/employer.proto

package employer

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

// The Employer model.
type Employer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email       string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Phone       string `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	CompanyName string `protobuf:"bytes,5,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
}

func (x *Employer) Reset() {
	*x = Employer{}
	mi := &file_proto_employer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Employer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Employer) ProtoMessage() {}

func (x *Employer) ProtoReflect() protoreflect.Message {
	mi := &file_proto_employer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Employer.ProtoReflect.Descriptor instead.
func (*Employer) Descriptor() ([]byte, []int) {
	return file_proto_employer_proto_rawDescGZIP(), []int{0}
}

func (x *Employer) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Employer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Employer) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Employer) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *Employer) GetCompanyName() string {
	if x != nil {
		return x.CompanyName
	}
	return ""
}

// Request message for listing employers.
type ListEmployersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListEmployersRequest) Reset() {
	*x = ListEmployersRequest{}
	mi := &file_proto_employer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEmployersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEmployersRequest) ProtoMessage() {}

func (x *ListEmployersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_employer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEmployersRequest.ProtoReflect.Descriptor instead.
func (*ListEmployersRequest) Descriptor() ([]byte, []int) {
	return file_proto_employer_proto_rawDescGZIP(), []int{1}
}

// Response message containing a list of employers.
type ListEmployersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Employers []*Employer `protobuf:"bytes,1,rep,name=employers,proto3" json:"employers,omitempty"`
}

func (x *ListEmployersResponse) Reset() {
	*x = ListEmployersResponse{}
	mi := &file_proto_employer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListEmployersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListEmployersResponse) ProtoMessage() {}

func (x *ListEmployersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_employer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListEmployersResponse.ProtoReflect.Descriptor instead.
func (*ListEmployersResponse) Descriptor() ([]byte, []int) {
	return file_proto_employer_proto_rawDescGZIP(), []int{2}
}

func (x *ListEmployersResponse) GetEmployers() []*Employer {
	if x != nil {
		return x.Employers
	}
	return nil
}

var File_proto_employer_proto protoreflect.FileDescriptor

var file_proto_employer_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72,
	0x22, 0x7d, 0x0a, 0x08, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x16, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x49, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x45,
	0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x30, 0x0a, 0x09, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x2e, 0x45,
	0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x52, 0x09, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x72, 0x73, 0x32, 0x63, 0x0a, 0x0f, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x50, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x6d, 0x70,
	0x6c, 0x6f, 0x79, 0x65, 0x72, 0x73, 0x12, 0x1e, 0x2e, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x28, 0x5a, 0x26, 0x66, 0x69, 0x6e, 0x65, 0x44,
	0x65, 0x65, 0x64, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x3b, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_employer_proto_rawDescOnce sync.Once
	file_proto_employer_proto_rawDescData = file_proto_employer_proto_rawDesc
)

func file_proto_employer_proto_rawDescGZIP() []byte {
	file_proto_employer_proto_rawDescOnce.Do(func() {
		file_proto_employer_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_employer_proto_rawDescData)
	})
	return file_proto_employer_proto_rawDescData
}

var file_proto_employer_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_employer_proto_goTypes = []any{
	(*Employer)(nil),              // 0: employer.Employer
	(*ListEmployersRequest)(nil),  // 1: employer.ListEmployersRequest
	(*ListEmployersResponse)(nil), // 2: employer.ListEmployersResponse
}
var file_proto_employer_proto_depIdxs = []int32{
	0, // 0: employer.ListEmployersResponse.employers:type_name -> employer.Employer
	1, // 1: employer.EmployerService.ListEmployers:input_type -> employer.ListEmployersRequest
	2, // 2: employer.EmployerService.ListEmployers:output_type -> employer.ListEmployersResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_employer_proto_init() }
func file_proto_employer_proto_init() {
	if File_proto_employer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_employer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_employer_proto_goTypes,
		DependencyIndexes: file_proto_employer_proto_depIdxs,
		MessageInfos:      file_proto_employer_proto_msgTypes,
	}.Build()
	File_proto_employer_proto = out.File
	file_proto_employer_proto_rawDesc = nil
	file_proto_employer_proto_goTypes = nil
	file_proto_employer_proto_depIdxs = nil
}
