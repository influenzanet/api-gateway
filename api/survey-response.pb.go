// Code generated by protoc-gen-go. DO NOT EDIT.
// source: survey-response.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SurveyResponse struct {
	Key                  string                `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	ParticipantId        string                `protobuf:"bytes,2,opt,name=participant_id,json=participantId,proto3" json:"participant_id,omitempty"`
	SubmittedAt          int64                 `protobuf:"varint,3,opt,name=submitted_at,json=submittedAt,proto3" json:"submitted_at,omitempty"`
	Responses            []*SurveyItemResponse `protobuf:"bytes,4,rep,name=responses,proto3" json:"responses,omitempty"`
	Context              map[string]string     `protobuf:"bytes,5,rep,name=context,proto3" json:"context,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *SurveyResponse) Reset()         { *m = SurveyResponse{} }
func (m *SurveyResponse) String() string { return proto.CompactTextString(m) }
func (*SurveyResponse) ProtoMessage()    {}
func (*SurveyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f7d4a107e829578, []int{0}
}

func (m *SurveyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SurveyResponse.Unmarshal(m, b)
}
func (m *SurveyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SurveyResponse.Marshal(b, m, deterministic)
}
func (m *SurveyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SurveyResponse.Merge(m, src)
}
func (m *SurveyResponse) XXX_Size() int {
	return xxx_messageInfo_SurveyResponse.Size(m)
}
func (m *SurveyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SurveyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SurveyResponse proto.InternalMessageInfo

func (m *SurveyResponse) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SurveyResponse) GetParticipantId() string {
	if m != nil {
		return m.ParticipantId
	}
	return ""
}

func (m *SurveyResponse) GetSubmittedAt() int64 {
	if m != nil {
		return m.SubmittedAt
	}
	return 0
}

func (m *SurveyResponse) GetResponses() []*SurveyItemResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

func (m *SurveyResponse) GetContext() map[string]string {
	if m != nil {
		return m.Context
	}
	return nil
}

type SurveyItemResponse struct {
	Key  string        `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Meta *ResponseMeta `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	// for item groups:
	Items []*SurveyItemResponse `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	// for single items:
	Response             *ResponseItem `protobuf:"bytes,4,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SurveyItemResponse) Reset()         { *m = SurveyItemResponse{} }
func (m *SurveyItemResponse) String() string { return proto.CompactTextString(m) }
func (*SurveyItemResponse) ProtoMessage()    {}
func (*SurveyItemResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f7d4a107e829578, []int{1}
}

func (m *SurveyItemResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SurveyItemResponse.Unmarshal(m, b)
}
func (m *SurveyItemResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SurveyItemResponse.Marshal(b, m, deterministic)
}
func (m *SurveyItemResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SurveyItemResponse.Merge(m, src)
}
func (m *SurveyItemResponse) XXX_Size() int {
	return xxx_messageInfo_SurveyItemResponse.Size(m)
}
func (m *SurveyItemResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SurveyItemResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SurveyItemResponse proto.InternalMessageInfo

func (m *SurveyItemResponse) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SurveyItemResponse) GetMeta() *ResponseMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *SurveyItemResponse) GetItems() []*SurveyItemResponse {
	if m != nil {
		return m.Items
	}
	return nil
}

func (m *SurveyItemResponse) GetResponse() *ResponseItem {
	if m != nil {
		return m.Response
	}
	return nil
}

type ResponseItem struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Dtype string `protobuf:"bytes,3,opt,name=dtype,proto3" json:"dtype,omitempty"`
	// For response option groups:
	Items                []*ResponseItem `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ResponseItem) Reset()         { *m = ResponseItem{} }
func (m *ResponseItem) String() string { return proto.CompactTextString(m) }
func (*ResponseItem) ProtoMessage()    {}
func (*ResponseItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f7d4a107e829578, []int{2}
}

func (m *ResponseItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseItem.Unmarshal(m, b)
}
func (m *ResponseItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseItem.Marshal(b, m, deterministic)
}
func (m *ResponseItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseItem.Merge(m, src)
}
func (m *ResponseItem) XXX_Size() int {
	return xxx_messageInfo_ResponseItem.Size(m)
}
func (m *ResponseItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseItem.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseItem proto.InternalMessageInfo

func (m *ResponseItem) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ResponseItem) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *ResponseItem) GetDtype() string {
	if m != nil {
		return m.Dtype
	}
	return ""
}

func (m *ResponseItem) GetItems() []*ResponseItem {
	if m != nil {
		return m.Items
	}
	return nil
}

type ResponseMeta struct {
	Position   int32  `protobuf:"varint,1,opt,name=position,proto3" json:"position,omitempty"`
	LocaleCode string `protobuf:"bytes,2,opt,name=locale_code,json=localeCode,proto3" json:"locale_code,omitempty"`
	Version    int32  `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	// timestamps:
	Rendered             []int64  `protobuf:"varint,4,rep,packed,name=rendered,proto3" json:"rendered,omitempty"`
	Displayed            []int64  `protobuf:"varint,5,rep,packed,name=displayed,proto3" json:"displayed,omitempty"`
	Responded            []int64  `protobuf:"varint,6,rep,packed,name=responded,proto3" json:"responded,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseMeta) Reset()         { *m = ResponseMeta{} }
func (m *ResponseMeta) String() string { return proto.CompactTextString(m) }
func (*ResponseMeta) ProtoMessage()    {}
func (*ResponseMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f7d4a107e829578, []int{3}
}

func (m *ResponseMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseMeta.Unmarshal(m, b)
}
func (m *ResponseMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseMeta.Marshal(b, m, deterministic)
}
func (m *ResponseMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseMeta.Merge(m, src)
}
func (m *ResponseMeta) XXX_Size() int {
	return xxx_messageInfo_ResponseMeta.Size(m)
}
func (m *ResponseMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseMeta.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseMeta proto.InternalMessageInfo

func (m *ResponseMeta) GetPosition() int32 {
	if m != nil {
		return m.Position
	}
	return 0
}

func (m *ResponseMeta) GetLocaleCode() string {
	if m != nil {
		return m.LocaleCode
	}
	return ""
}

func (m *ResponseMeta) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ResponseMeta) GetRendered() []int64 {
	if m != nil {
		return m.Rendered
	}
	return nil
}

func (m *ResponseMeta) GetDisplayed() []int64 {
	if m != nil {
		return m.Displayed
	}
	return nil
}

func (m *ResponseMeta) GetResponded() []int64 {
	if m != nil {
		return m.Responded
	}
	return nil
}

func init() {
	proto.RegisterType((*SurveyResponse)(nil), "inf.survey_response.SurveyResponse")
	proto.RegisterMapType((map[string]string)(nil), "inf.survey_response.SurveyResponse.ContextEntry")
	proto.RegisterType((*SurveyItemResponse)(nil), "inf.survey_response.SurveyItemResponse")
	proto.RegisterType((*ResponseItem)(nil), "inf.survey_response.ResponseItem")
	proto.RegisterType((*ResponseMeta)(nil), "inf.survey_response.ResponseMeta")
}

func init() {
	proto.RegisterFile("survey-response.proto", fileDescriptor_9f7d4a107e829578)
}

var fileDescriptor_9f7d4a107e829578 = []byte{
	// 420 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xd1, 0x8a, 0x13, 0x31,
	0x14, 0x65, 0x9a, 0xce, 0xee, 0xf6, 0xb6, 0x2e, 0x12, 0x15, 0xc2, 0x22, 0xd8, 0x2d, 0x88, 0x7d,
	0x71, 0x90, 0x15, 0x51, 0x16, 0xf6, 0x41, 0x97, 0x7d, 0x58, 0xc1, 0x97, 0xf8, 0xe6, 0x4b, 0xc9,
	0x36, 0x57, 0x08, 0xb6, 0x93, 0x90, 0xdc, 0x16, 0xe7, 0x03, 0xfc, 0x2a, 0xf1, 0x4b, 0xfc, 0x19,
	0x99, 0x64, 0x33, 0x2d, 0x38, 0x48, 0x7d, 0x9b, 0x7b, 0xce, 0x3d, 0x27, 0xe7, 0xde, 0x49, 0xe0,
	0x49, 0xd8, 0xf8, 0x2d, 0x36, 0x2f, 0x3d, 0x06, 0x67, 0xeb, 0x80, 0x95, 0xf3, 0x96, 0x2c, 0x7f,
	0x64, 0xea, 0xaf, 0x55, 0xa2, 0x16, 0x99, 0x9a, 0xfd, 0x1c, 0xc0, 0xe9, 0xe7, 0x88, 0xc9, 0x7b,
	0x88, 0x3f, 0x04, 0xf6, 0x0d, 0x1b, 0x51, 0x4c, 0x8b, 0xf9, 0x48, 0xb6, 0x9f, 0xfc, 0x39, 0x9c,
	0x3a, 0xe5, 0xc9, 0x2c, 0x8d, 0x53, 0x35, 0x2d, 0x8c, 0x16, 0x83, 0x48, 0x3e, 0xd8, 0x43, 0x6f,
	0x35, 0x3f, 0x87, 0x49, 0xd8, 0xdc, 0xad, 0x0d, 0x11, 0xea, 0x85, 0x22, 0xc1, 0xa6, 0xc5, 0x9c,
	0xc9, 0x71, 0x87, 0xbd, 0x27, 0x7e, 0x03, 0xa3, 0x7c, 0x74, 0x10, 0xc3, 0x29, 0x9b, 0x8f, 0x2f,
	0x5e, 0x54, 0x3d, 0xb9, 0xaa, 0x94, 0xe9, 0x96, 0x70, 0x9d, 0x73, 0xc9, 0x9d, 0x92, 0x7f, 0x84,
	0xe3, 0xa5, 0xad, 0x09, 0xbf, 0x93, 0x28, 0xa3, 0xc9, 0xab, 0x7f, 0x98, 0x64, 0x83, 0xea, 0x3a,
	0x49, 0x6e, 0x6a, 0xf2, 0x8d, 0xcc, 0x06, 0x67, 0x97, 0x30, 0xd9, 0x27, 0x7a, 0xc6, 0x7f, 0x0c,
	0xe5, 0x56, 0xad, 0x36, 0x78, 0x3f, 0x75, 0x2a, 0x2e, 0x07, 0xef, 0x8a, 0xd9, 0xef, 0x02, 0xf8,
	0xdf, 0x49, 0x7b, 0x2c, 0xde, 0xc0, 0x70, 0x8d, 0xa4, 0xa2, 0xc3, 0xf8, 0xe2, 0xbc, 0x37, 0x6d,
	0x96, 0x7f, 0x42, 0x52, 0x32, 0xb6, 0xf3, 0x2b, 0x28, 0x0d, 0xe1, 0x3a, 0x08, 0xf6, 0x7f, 0xab,
	0x4a, 0x2a, 0x7e, 0x05, 0x27, 0xb9, 0x4b, 0x0c, 0x0f, 0x38, 0x39, 0x7a, 0x74, 0x92, 0xd9, 0x8f,
	0x02, 0x26, 0xfb, 0xd4, 0xa1, 0xab, 0x69, 0x51, 0x4d, 0x8d, 0xc3, 0x78, 0x03, 0x46, 0x32, 0x15,
	0xfc, 0x6d, 0x1e, 0x26, 0xfd, 0xf7, 0x03, 0xa2, 0xa4, 0xfe, 0xd9, 0xaf, 0xbd, 0x1c, 0xed, 0x72,
	0xf8, 0x19, 0x9c, 0x38, 0x1b, 0x0c, 0x19, 0x5b, 0xc7, 0x30, 0xa5, 0xec, 0x6a, 0xfe, 0x0c, 0xc6,
	0x2b, 0xbb, 0x54, 0x2b, 0x5c, 0x2c, 0xad, 0xce, 0xb9, 0x20, 0x41, 0xd7, 0x56, 0x23, 0x17, 0x70,
	0xbc, 0x45, 0x1f, 0x5a, 0x2d, 0x8b, 0xda, 0x5c, 0xb6, 0xb6, 0x1e, 0x6b, 0x8d, 0x1e, 0x75, 0xcc,
	0xc8, 0x64, 0x57, 0xf3, 0xa7, 0x30, 0xd2, 0x26, 0xb8, 0x95, 0x6a, 0x50, 0xc7, 0x3b, 0xc7, 0xe4,
	0x0e, 0x68, 0xd9, 0x34, 0x81, 0x46, 0x2d, 0x8e, 0x12, 0xdb, 0x01, 0x1f, 0xca, 0x2f, 0x4c, 0x39,
	0x73, 0x77, 0x14, 0x9f, 0xe1, 0xeb, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x42, 0xb1, 0x67, 0x58,
	0x9f, 0x03, 0x00, 0x00,
}