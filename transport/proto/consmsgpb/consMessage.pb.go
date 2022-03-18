// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.3
// source: consMessage.proto

package consmsgpb

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

type MessageType int32

const (
	MessageType_VAL                  MessageType = 0
	MessageType_ECHO                 MessageType = 1
	MessageType_READY                MessageType = 2
	MessageType_BVAL                 MessageType = 3
	MessageType_AUX                  MessageType = 4
	MessageType_COIN                 MessageType = 5
	MessageType_CON                  MessageType = 6
	MessageType_SKIP                 MessageType = 7
	MessageType_ECHO_COLLECTION      MessageType = 8
	MessageType_READY_COLLECTION     MessageType = 9
	MessageType_BVAL_ZERO_COLLECTION MessageType = 10
	MessageType_BVAL_ONE_COLLECTION  MessageType = 11
	MessageType_AUX_ZERO_COLLECTION  MessageType = 12
	MessageType_AUX_ONE_COLLECTION   MessageType = 13
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0:  "VAL",
		1:  "ECHO",
		2:  "READY",
		3:  "BVAL",
		4:  "AUX",
		5:  "COIN",
		6:  "CON",
		7:  "SKIP",
		8:  "ECHO_COLLECTION",
		9:  "READY_COLLECTION",
		10: "BVAL_ZERO_COLLECTION",
		11: "BVAL_ONE_COLLECTION",
		12: "AUX_ZERO_COLLECTION",
		13: "AUX_ONE_COLLECTION",
	}
	MessageType_value = map[string]int32{
		"VAL":                  0,
		"ECHO":                 1,
		"READY":                2,
		"BVAL":                 3,
		"AUX":                  4,
		"COIN":                 5,
		"CON":                  6,
		"SKIP":                 7,
		"ECHO_COLLECTION":      8,
		"READY_COLLECTION":     9,
		"BVAL_ZERO_COLLECTION": 10,
		"BVAL_ONE_COLLECTION":  11,
		"AUX_ZERO_COLLECTION":  12,
		"AUX_ONE_COLLECTION":   13,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_consMessage_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_consMessage_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_consMessage_proto_rawDescGZIP(), []int{0}
}

type Content struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (x *Content) Reset() {
	*x = Content{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consMessage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content) ProtoMessage() {}

func (x *Content) ProtoReflect() protoreflect.Message {
	mi := &file_consMessage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content.ProtoReflect.Descriptor instead.
func (*Content) Descriptor() ([]byte, []int) {
	return file_consMessage_proto_rawDescGZIP(), []int{0}
}

func (x *Content) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type Collections struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Collections [][]byte `protobuf:"bytes,1,rep,name=Collections,proto3" json:"Collections,omitempty"`
}

func (x *Collections) Reset() {
	*x = Collections{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consMessage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Collections) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Collections) ProtoMessage() {}

func (x *Collections) ProtoReflect() protoreflect.Message {
	mi := &file_consMessage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Collections.ProtoReflect.Descriptor instead.
func (*Collections) Descriptor() ([]byte, []int) {
	return file_consMessage_proto_rawDescGZIP(), []int{1}
}

func (x *Collections) GetCollections() [][]byte {
	if x != nil {
		return x.Collections
	}
	return nil
}

type WholeMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConsMsg    *ConsMessage `protobuf:"bytes,1,opt,name=ConsMsg,proto3" json:"ConsMsg,omitempty"`
	From       uint32       `protobuf:"varint,2,opt,name=From,proto3" json:"From,omitempty"`
	Signature  []byte       `protobuf:"bytes,3,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Collection []byte       `protobuf:"bytes,4,opt,name=Collection,proto3" json:"Collection,omitempty"`
}

func (x *WholeMessage) Reset() {
	*x = WholeMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consMessage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WholeMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WholeMessage) ProtoMessage() {}

func (x *WholeMessage) ProtoReflect() protoreflect.Message {
	mi := &file_consMessage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WholeMessage.ProtoReflect.Descriptor instead.
func (*WholeMessage) Descriptor() ([]byte, []int) {
	return file_consMessage_proto_rawDescGZIP(), []int{2}
}

func (x *WholeMessage) GetConsMsg() *ConsMessage {
	if x != nil {
		return x.ConsMsg
	}
	return nil
}

func (x *WholeMessage) GetFrom() uint32 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *WholeMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *WholeMessage) GetCollection() []byte {
	if x != nil {
		return x.Collection
	}
	return nil
}

type ConsMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     MessageType `protobuf:"varint,1,opt,name=type,proto3,enum=consmsg.MessageType" json:"type,omitempty"`
	Proposer uint32      `protobuf:"varint,2,opt,name=proposer,proto3" json:"proposer,omitempty"`
	Round    uint32      `protobuf:"varint,3,opt,name=round,proto3" json:"round,omitempty"`
	Sequence uint64      `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Content  []byte      `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	Single   bool        `protobuf:"varint,6,opt,name=single,proto3" json:"single,omitempty"`
	BinVals  []byte      `protobuf:"bytes,7,opt,name=binVals,proto3" json:"binVals,omitempty"`
}

func (x *ConsMessage) Reset() {
	*x = ConsMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consMessage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsMessage) ProtoMessage() {}

func (x *ConsMessage) ProtoReflect() protoreflect.Message {
	mi := &file_consMessage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsMessage.ProtoReflect.Descriptor instead.
func (*ConsMessage) Descriptor() ([]byte, []int) {
	return file_consMessage_proto_rawDescGZIP(), []int{3}
}

func (x *ConsMessage) GetType() MessageType {
	if x != nil {
		return x.Type
	}
	return MessageType_VAL
}

func (x *ConsMessage) GetProposer() uint32 {
	if x != nil {
		return x.Proposer
	}
	return 0
}

func (x *ConsMessage) GetRound() uint32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *ConsMessage) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *ConsMessage) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *ConsMessage) GetSingle() bool {
	if x != nil {
		return x.Single
	}
	return false
}

func (x *ConsMessage) GetBinVals() []byte {
	if x != nil {
		return x.BinVals
	}
	return nil
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From      uint32 `protobuf:"varint,1,opt,name=From,proto3" json:"From,omitempty"`
	Sequence  uint64 `protobuf:"varint,2,opt,name=Sequence,proto3" json:"Sequence,omitempty"`
	Signature []byte `protobuf:"bytes,3,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Content   []byte `protobuf:"bytes,4,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_consMessage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_consMessage_proto_msgTypes[4]
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
	return file_consMessage_proto_rawDescGZIP(), []int{4}
}

func (x *Request) GetFrom() uint32 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *Request) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *Request) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Request) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_consMessage_proto protoreflect.FileDescriptor

var file_consMessage_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6d, 0x73, 0x67, 0x22, 0x23, 0x0a, 0x07,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x22, 0x2f, 0x0a, 0x0b, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0b, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x22, 0x90, 0x01, 0x0a, 0x0c, 0x57, 0x68, 0x6f, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x73, 0x4d, 0x73, 0x67, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x43,
	0x6f, 0x6e, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x73,
	0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xd1, 0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x73, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x69, 0x6e, 0x67, 0x6c,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x62, 0x69, 0x6e, 0x56, 0x61, 0x6c, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x62, 0x69, 0x6e, 0x56, 0x61, 0x6c, 0x73, 0x22, 0x71, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2a, 0xea, 0x01, 0x0a,
	0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03,
	0x56, 0x41, 0x4c, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x45, 0x43, 0x48, 0x4f, 0x10, 0x01, 0x12,
	0x09, 0x0a, 0x05, 0x52, 0x45, 0x41, 0x44, 0x59, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x56,
	0x41, 0x4c, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x55, 0x58, 0x10, 0x04, 0x12, 0x08, 0x0a,
	0x04, 0x43, 0x4f, 0x49, 0x4e, 0x10, 0x05, 0x12, 0x07, 0x0a, 0x03, 0x43, 0x4f, 0x4e, 0x10, 0x06,
	0x12, 0x08, 0x0a, 0x04, 0x53, 0x4b, 0x49, 0x50, 0x10, 0x07, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x43,
	0x48, 0x4f, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x08, 0x12,
	0x14, 0x0a, 0x10, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54,
	0x49, 0x4f, 0x4e, 0x10, 0x09, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x56, 0x41, 0x4c, 0x5f, 0x5a, 0x45,
	0x52, 0x4f, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x0a, 0x12,
	0x17, 0x0a, 0x13, 0x42, 0x56, 0x41, 0x4c, 0x5f, 0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x4c,
	0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x0b, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x55, 0x58, 0x5f,
	0x5a, 0x45, 0x52, 0x4f, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10,
	0x0c, 0x12, 0x16, 0x0a, 0x12, 0x41, 0x55, 0x58, 0x5f, 0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x4f, 0x4c,
	0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x0d, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x63,
	0x6f, 0x6e, 0x73, 0x6d, 0x73, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_consMessage_proto_rawDescOnce sync.Once
	file_consMessage_proto_rawDescData = file_consMessage_proto_rawDesc
)

func file_consMessage_proto_rawDescGZIP() []byte {
	file_consMessage_proto_rawDescOnce.Do(func() {
		file_consMessage_proto_rawDescData = protoimpl.X.CompressGZIP(file_consMessage_proto_rawDescData)
	})
	return file_consMessage_proto_rawDescData
}

var file_consMessage_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_consMessage_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_consMessage_proto_goTypes = []interface{}{
	(MessageType)(0),     // 0: consmsg.MessageType
	(*Content)(nil),      // 1: consmsg.Content
	(*Collections)(nil),  // 2: consmsg.Collections
	(*WholeMessage)(nil), // 3: consmsg.WholeMessage
	(*ConsMessage)(nil),  // 4: consmsg.ConsMessage
	(*Request)(nil),      // 5: consmsg.Request
}
var file_consMessage_proto_depIdxs = []int32{
	4, // 0: consmsg.WholeMessage.ConsMsg:type_name -> consmsg.ConsMessage
	0, // 1: consmsg.ConsMessage.type:type_name -> consmsg.MessageType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_consMessage_proto_init() }
func file_consMessage_proto_init() {
	if File_consMessage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_consMessage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content); i {
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
		file_consMessage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Collections); i {
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
		file_consMessage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WholeMessage); i {
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
		file_consMessage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsMessage); i {
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
		file_consMessage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_consMessage_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_consMessage_proto_goTypes,
		DependencyIndexes: file_consMessage_proto_depIdxs,
		EnumInfos:         file_consMessage_proto_enumTypes,
		MessageInfos:      file_consMessage_proto_msgTypes,
	}.Build()
	File_consMessage_proto = out.File
	file_consMessage_proto_rawDesc = nil
	file_consMessage_proto_goTypes = nil
	file_consMessage_proto_depIdxs = nil
}
