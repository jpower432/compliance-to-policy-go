// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.0--rc2
// source: api/proto/v1alpha1/models.proto

package v1alpha1

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

type Result int32

const (
	Result_RESULT_UNSPECIFIED Result = 0
	Result_RESULT_ERROR       Result = 1
	Result_RESULT_WARNING     Result = 2
	Result_RESULT_PASS        Result = 3
	Result_RESULT_FAILURE     Result = 4
)

// Enum value maps for Result.
var (
	Result_name = map[int32]string{
		0: "RESULT_UNSPECIFIED",
		1: "RESULT_ERROR",
		2: "RESULT_WARNING",
		3: "RESULT_PASS",
		4: "RESULT_FAILURE",
	}
	Result_value = map[string]int32{
		"RESULT_UNSPECIFIED": 0,
		"RESULT_ERROR":       1,
		"RESULT_WARNING":     2,
		"RESULT_PASS":        3,
		"RESULT_FAILURE":     4,
	}
)

func (x Result) Enum() *Result {
	p := new(Result)
	*p = x
	return p
}

func (x Result) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Result) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_v1alpha1_models_proto_enumTypes[0].Descriptor()
}

func (Result) Type() protoreflect.EnumType {
	return &file_api_proto_v1alpha1_models_proto_enumTypes[0]
}

func (x Result) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Result.Descriptor instead.
func (Result) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{0}
}

// Define a single rule parameter
type Parameter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name is the human-readable parameter identifier
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// description is the human-readable documentation for the parameter
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// selected values for the parameter
	SelectedValue string `protobuf:"bytes,3,opt,name=selected_value,json=selectedValue,proto3" json:"selected_value,omitempty"`
}

func (x *Parameter) Reset() {
	*x = Parameter{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Parameter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Parameter) ProtoMessage() {}

func (x *Parameter) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Parameter.ProtoReflect.Descriptor instead.
func (*Parameter) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{0}
}

func (x *Parameter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Parameter) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Parameter) GetSelectedValue() string {
	if x != nil {
		return x.SelectedValue
	}
	return ""
}

// Define a single check
type Check struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name is the human-readable check identifier
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// description is the human-readable documentation for the check
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Check) Reset() {
	*x = Check{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Check) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Check) ProtoMessage() {}

func (x *Check) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Check.ProtoReflect.Descriptor instead.
func (*Check) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{1}
}

func (x *Check) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Check) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// Define a single rule
type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name is the human-readable technical rule identifier
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// description is the human-readable documentation for the technical rule
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Check Mapped to rule
	Check *Check `protobuf:"bytes,4,opt,name=check,proto3" json:"check,omitempty"`
	// Parameter associated with Rule
	Parameter *Parameter `protobuf:"bytes,5,opt,name=parameter,proto3,oneof" json:"parameter,omitempty"`
}

func (x *Rule) Reset() {
	*x = Rule{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{2}
}

func (x *Rule) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Rule) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Rule) GetCheck() *Check {
	if x != nil {
		return x.Check
	}
	return nil
}

func (x *Rule) GetParameter() *Parameter {
	if x != nil {
		return x.Parameter
	}
	return nil
}

type Policy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rules      []*Rule      `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
	Parameters []*Parameter `protobuf:"bytes,2,rep,name=parameters,proto3" json:"parameters,omitempty"`
}

func (x *Policy) Reset() {
	*x = Policy{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Policy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Policy) ProtoMessage() {}

func (x *Policy) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Policy.ProtoReflect.Descriptor instead.
func (*Policy) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{3}
}

func (x *Policy) GetRules() []*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

func (x *Policy) GetParameters() []*Parameter {
	if x != nil {
		return x.Parameters
	}
	return nil
}

type Subject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title       string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	ResourceId  string                 `protobuf:"bytes,2,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Result      Result                 `protobuf:"varint,3,opt,name=result,proto3,enum=protocols.Result" json:"result,omitempty"`
	EvaluatedOn *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=evaluated_on,json=evaluatedOn,proto3" json:"evaluated_on,omitempty"`
	Reason      string                 `protobuf:"bytes,5,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *Subject) Reset() {
	*x = Subject{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Subject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subject) ProtoMessage() {}

func (x *Subject) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subject.ProtoReflect.Descriptor instead.
func (*Subject) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{4}
}

func (x *Subject) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Subject) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *Subject) GetResult() Result {
	if x != nil {
		return x.Result
	}
	return Result_RESULT_UNSPECIFIED
}

func (x *Subject) GetEvaluatedOn() *timestamppb.Timestamp {
	if x != nil {
		return x.EvaluatedOn
	}
	return nil
}

func (x *Subject) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type ObservationByCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description  string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	CheckId      string                 `protobuf:"bytes,3,opt,name=check_id,json=checkId,proto3" json:"check_id,omitempty"`
	Methods      []string               `protobuf:"bytes,4,rep,name=methods,proto3" json:"methods,omitempty"`
	CollectedAt  *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=collected_at,json=collectedAt,proto3" json:"collected_at,omitempty"`
	Subjects     []*Subject             `protobuf:"bytes,6,rep,name=subjects,proto3" json:"subjects,omitempty"`
	EvidenceRefs []string               `protobuf:"bytes,7,rep,name=evidence_refs,json=evidenceRefs,proto3" json:"evidence_refs,omitempty"`
}

func (x *ObservationByCheck) Reset() {
	*x = ObservationByCheck{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ObservationByCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObservationByCheck) ProtoMessage() {}

func (x *ObservationByCheck) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObservationByCheck.ProtoReflect.Descriptor instead.
func (*ObservationByCheck) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{5}
}

func (x *ObservationByCheck) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ObservationByCheck) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ObservationByCheck) GetCheckId() string {
	if x != nil {
		return x.CheckId
	}
	return ""
}

func (x *ObservationByCheck) GetMethods() []string {
	if x != nil {
		return x.Methods
	}
	return nil
}

func (x *ObservationByCheck) GetCollectedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CollectedAt
	}
	return nil
}

func (x *ObservationByCheck) GetSubjects() []*Subject {
	if x != nil {
		return x.Subjects
	}
	return nil
}

func (x *ObservationByCheck) GetEvidenceRefs() []string {
	if x != nil {
		return x.EvidenceRefs
	}
	return nil
}

// OSCAL Finding
type Finding struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name                string                `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description         string                `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	RelatedObservations []*ObservationByCheck `protobuf:"bytes,3,rep,name=related_observations,json=relatedObservations,proto3" json:"related_observations,omitempty"`
}

func (x *Finding) Reset() {
	*x = Finding{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Finding) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Finding) ProtoMessage() {}

func (x *Finding) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Finding.ProtoReflect.Descriptor instead.
func (*Finding) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{6}
}

func (x *Finding) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Finding) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Finding) GetRelatedObservations() []*ObservationByCheck {
	if x != nil {
		return x.RelatedObservations
	}
	return nil
}

type PVPResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Observations []*ObservationByCheck `protobuf:"bytes,1,rep,name=observations,proto3" json:"observations,omitempty"`
	Links        map[string]string     `protobuf:"bytes,2,rep,name=links,proto3" json:"links,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *PVPResult) Reset() {
	*x = PVPResult{}
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PVPResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PVPResult) ProtoMessage() {}

func (x *PVPResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1alpha1_models_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PVPResult.ProtoReflect.Descriptor instead.
func (*PVPResult) Descriptor() ([]byte, []int) {
	return file_api_proto_v1alpha1_models_proto_rawDescGZIP(), []int{7}
}

func (x *PVPResult) GetObservations() []*ObservationByCheck {
	if x != nil {
		return x.Observations
	}
	return nil
}

func (x *PVPResult) GetLinks() map[string]string {
	if x != nil {
		return x.Links
	}
	return nil
}

var File_api_proto_v1alpha1_models_proto protoreflect.FileDescriptor

var file_api_proto_v1alpha1_models_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x68, 0x0a,
	0x09, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3d, 0x0a, 0x05, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xab, 0x01, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x05, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73,
	0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x05, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x37, 0x0a,
	0x09, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x2e, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x48, 0x00, 0x52, 0x09, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x22, 0x65, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x25,
	0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x34, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22, 0xc2, 0x01, 0x0a, 0x07,
	0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x29,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x3d, 0x0a, 0x0c, 0x65, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x65, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x65, 0x64, 0x4f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x22, 0x93, 0x02, 0x0a, 0x12, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x42, 0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a,
	0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x73, 0x12, 0x3d, 0x0a, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x2e, 0x0a, 0x08, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x2e,
	0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x08, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x73, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x72, 0x65,
	0x66, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e,
	0x63, 0x65, 0x52, 0x65, 0x66, 0x73, 0x22, 0x91, 0x01, 0x0a, 0x07, 0x46, 0x69, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x50, 0x0a, 0x14, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x73, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x13, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x4f, 0x62,
	0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xbf, 0x01, 0x0a, 0x09, 0x50,
	0x56, 0x50, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x41, 0x0a, 0x0c, 0x6f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x0c, 0x6f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x35, 0x0a, 0x05, 0x6c,
	0x69, 0x6e, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x73, 0x2e, 0x50, 0x56, 0x50, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x6c, 0x69, 0x6e,
	0x6b, 0x73, 0x1a, 0x38, 0x0a, 0x0a, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x6b, 0x0a, 0x06,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x0c, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01,
	0x12, 0x12, 0x0a, 0x0e, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x57, 0x41, 0x52, 0x4e, 0x49,
	0x4e, 0x47, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x50,
	0x41, 0x53, 0x53, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x04, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x2d, 0x63, 0x6f,
	0x6d, 0x70, 0x61, 0x73, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65,
	0x2d, 0x74, 0x6f, 0x2d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2d, 0x67, 0x6f, 0x2f, 0x76, 0x32,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_v1alpha1_models_proto_rawDescOnce sync.Once
	file_api_proto_v1alpha1_models_proto_rawDescData = file_api_proto_v1alpha1_models_proto_rawDesc
)

func file_api_proto_v1alpha1_models_proto_rawDescGZIP() []byte {
	file_api_proto_v1alpha1_models_proto_rawDescOnce.Do(func() {
		file_api_proto_v1alpha1_models_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_v1alpha1_models_proto_rawDescData)
	})
	return file_api_proto_v1alpha1_models_proto_rawDescData
}

var file_api_proto_v1alpha1_models_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_v1alpha1_models_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_proto_v1alpha1_models_proto_goTypes = []any{
	(Result)(0),                   // 0: protocols.Result
	(*Parameter)(nil),             // 1: protocols.Parameter
	(*Check)(nil),                 // 2: protocols.Check
	(*Rule)(nil),                  // 3: protocols.Rule
	(*Policy)(nil),                // 4: protocols.Policy
	(*Subject)(nil),               // 5: protocols.Subject
	(*ObservationByCheck)(nil),    // 6: protocols.ObservationByCheck
	(*Finding)(nil),               // 7: protocols.Finding
	(*PVPResult)(nil),             // 8: protocols.PVPResult
	nil,                           // 9: protocols.PVPResult.LinksEntry
	(*timestamppb.Timestamp)(nil), // 10: google.protobuf.Timestamp
}
var file_api_proto_v1alpha1_models_proto_depIdxs = []int32{
	2,  // 0: protocols.Rule.check:type_name -> protocols.Check
	1,  // 1: protocols.Rule.parameter:type_name -> protocols.Parameter
	3,  // 2: protocols.Policy.rules:type_name -> protocols.Rule
	1,  // 3: protocols.Policy.parameters:type_name -> protocols.Parameter
	0,  // 4: protocols.Subject.result:type_name -> protocols.Result
	10, // 5: protocols.Subject.evaluated_on:type_name -> google.protobuf.Timestamp
	10, // 6: protocols.ObservationByCheck.collected_at:type_name -> google.protobuf.Timestamp
	5,  // 7: protocols.ObservationByCheck.subjects:type_name -> protocols.Subject
	6,  // 8: protocols.Finding.related_observations:type_name -> protocols.ObservationByCheck
	6,  // 9: protocols.PVPResult.observations:type_name -> protocols.ObservationByCheck
	9,  // 10: protocols.PVPResult.links:type_name -> protocols.PVPResult.LinksEntry
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_api_proto_v1alpha1_models_proto_init() }
func file_api_proto_v1alpha1_models_proto_init() {
	if File_api_proto_v1alpha1_models_proto != nil {
		return
	}
	file_api_proto_v1alpha1_models_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_v1alpha1_models_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_v1alpha1_models_proto_goTypes,
		DependencyIndexes: file_api_proto_v1alpha1_models_proto_depIdxs,
		EnumInfos:         file_api_proto_v1alpha1_models_proto_enumTypes,
		MessageInfos:      file_api_proto_v1alpha1_models_proto_msgTypes,
	}.Build()
	File_api_proto_v1alpha1_models_proto = out.File
	file_api_proto_v1alpha1_models_proto_rawDesc = nil
	file_api_proto_v1alpha1_models_proto_goTypes = nil
	file_api_proto_v1alpha1_models_proto_depIdxs = nil
}