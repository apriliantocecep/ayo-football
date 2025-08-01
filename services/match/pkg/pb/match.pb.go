// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.28.0
// source: match.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Match struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Date          string                 `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	Venue         string                 `protobuf:"bytes,3,opt,name=venue,proto3" json:"venue,omitempty"`
	HomeTeamId    string                 `protobuf:"bytes,4,opt,name=home_team_id,json=homeTeamId,proto3" json:"home_team_id,omitempty"`
	AwayTeamId    string                 `protobuf:"bytes,5,opt,name=away_team_id,json=awayTeamId,proto3" json:"away_team_id,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     string                 `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Match) Reset() {
	*x = Match{}
	mi := &file_match_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Match) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Match) ProtoMessage() {}

func (x *Match) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Match.ProtoReflect.Descriptor instead.
func (*Match) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{0}
}

func (x *Match) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Match) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *Match) GetVenue() string {
	if x != nil {
		return x.Venue
	}
	return ""
}

func (x *Match) GetHomeTeamId() string {
	if x != nil {
		return x.HomeTeamId
	}
	return ""
}

func (x *Match) GetAwayTeamId() string {
	if x != nil {
		return x.AwayTeamId
	}
	return ""
}

func (x *Match) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Match) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

type Goal struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MatchId       string                 `protobuf:"bytes,2,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	PlayerId      string                 `protobuf:"bytes,3,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	ScoredAt      string                 `protobuf:"bytes,4,opt,name=scored_at,json=scoredAt,proto3" json:"scored_at,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     string                 `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Goal) Reset() {
	*x = Goal{}
	mi := &file_match_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Goal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Goal) ProtoMessage() {}

func (x *Goal) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Goal.ProtoReflect.Descriptor instead.
func (*Goal) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{1}
}

func (x *Goal) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Goal) GetMatchId() string {
	if x != nil {
		return x.MatchId
	}
	return ""
}

func (x *Goal) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *Goal) GetScoredAt() string {
	if x != nil {
		return x.ScoredAt
	}
	return ""
}

func (x *Goal) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Goal) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

type CreateMatchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Date          string                 `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Venue         string                 `protobuf:"bytes,2,opt,name=venue,proto3" json:"venue,omitempty"`
	HomeTeamId    string                 `protobuf:"bytes,3,opt,name=home_team_id,json=homeTeamId,proto3" json:"home_team_id,omitempty"`
	AwayTeamId    string                 `protobuf:"bytes,4,opt,name=away_team_id,json=awayTeamId,proto3" json:"away_team_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateMatchRequest) Reset() {
	*x = CreateMatchRequest{}
	mi := &file_match_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMatchRequest) ProtoMessage() {}

func (x *CreateMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMatchRequest.ProtoReflect.Descriptor instead.
func (*CreateMatchRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{2}
}

func (x *CreateMatchRequest) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *CreateMatchRequest) GetVenue() string {
	if x != nil {
		return x.Venue
	}
	return ""
}

func (x *CreateMatchRequest) GetHomeTeamId() string {
	if x != nil {
		return x.HomeTeamId
	}
	return ""
}

func (x *CreateMatchRequest) GetAwayTeamId() string {
	if x != nil {
		return x.AwayTeamId
	}
	return ""
}

type UpdateMatchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Date          string                 `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	Venue         string                 `protobuf:"bytes,3,opt,name=venue,proto3" json:"venue,omitempty"`
	HomeTeamId    string                 `protobuf:"bytes,4,opt,name=home_team_id,json=homeTeamId,proto3" json:"home_team_id,omitempty"`
	AwayTeamId    string                 `protobuf:"bytes,5,opt,name=away_team_id,json=awayTeamId,proto3" json:"away_team_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateMatchRequest) Reset() {
	*x = UpdateMatchRequest{}
	mi := &file_match_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMatchRequest) ProtoMessage() {}

func (x *UpdateMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMatchRequest.ProtoReflect.Descriptor instead.
func (*UpdateMatchRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateMatchRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateMatchRequest) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *UpdateMatchRequest) GetVenue() string {
	if x != nil {
		return x.Venue
	}
	return ""
}

func (x *UpdateMatchRequest) GetHomeTeamId() string {
	if x != nil {
		return x.HomeTeamId
	}
	return ""
}

func (x *UpdateMatchRequest) GetAwayTeamId() string {
	if x != nil {
		return x.AwayTeamId
	}
	return ""
}

type DeleteMatchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteMatchRequest) Reset() {
	*x = DeleteMatchRequest{}
	mi := &file_match_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMatchRequest) ProtoMessage() {}

func (x *DeleteMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMatchRequest.ProtoReflect.Descriptor instead.
func (*DeleteMatchRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteMatchRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetMatchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMatchRequest) Reset() {
	*x = GetMatchRequest{}
	mi := &file_match_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMatchRequest) ProtoMessage() {}

func (x *GetMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMatchRequest.ProtoReflect.Descriptor instead.
func (*GetMatchRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{5}
}

func (x *GetMatchRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ListMatchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListMatchRequest) Reset() {
	*x = ListMatchRequest{}
	mi := &file_match_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMatchRequest) ProtoMessage() {}

func (x *ListMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMatchRequest.ProtoReflect.Descriptor instead.
func (*ListMatchRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{6}
}

func (x *ListMatchRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListMatchRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ListMatchResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Matches       []*Match               `protobuf:"bytes,1,rep,name=matches,proto3" json:"matches,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListMatchResponse) Reset() {
	*x = ListMatchResponse{}
	mi := &file_match_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMatchResponse) ProtoMessage() {}

func (x *ListMatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMatchResponse.ProtoReflect.Descriptor instead.
func (*ListMatchResponse) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{7}
}

func (x *ListMatchResponse) GetMatches() []*Match {
	if x != nil {
		return x.Matches
	}
	return nil
}

type CreateGoalRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MatchId       string                 `protobuf:"bytes,1,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	PlayerId      string                 `protobuf:"bytes,2,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	ScoredAt      string                 `protobuf:"bytes,3,opt,name=scored_at,json=scoredAt,proto3" json:"scored_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateGoalRequest) Reset() {
	*x = CreateGoalRequest{}
	mi := &file_match_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateGoalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateGoalRequest) ProtoMessage() {}

func (x *CreateGoalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateGoalRequest.ProtoReflect.Descriptor instead.
func (*CreateGoalRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{8}
}

func (x *CreateGoalRequest) GetMatchId() string {
	if x != nil {
		return x.MatchId
	}
	return ""
}

func (x *CreateGoalRequest) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *CreateGoalRequest) GetScoredAt() string {
	if x != nil {
		return x.ScoredAt
	}
	return ""
}

type UpdateGoalRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MatchId       string                 `protobuf:"bytes,2,opt,name=match_id,json=matchId,proto3" json:"match_id,omitempty"`
	PlayerId      string                 `protobuf:"bytes,3,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	ScoredAt      string                 `protobuf:"bytes,4,opt,name=scored_at,json=scoredAt,proto3" json:"scored_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateGoalRequest) Reset() {
	*x = UpdateGoalRequest{}
	mi := &file_match_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateGoalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGoalRequest) ProtoMessage() {}

func (x *UpdateGoalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGoalRequest.ProtoReflect.Descriptor instead.
func (*UpdateGoalRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{9}
}

func (x *UpdateGoalRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateGoalRequest) GetMatchId() string {
	if x != nil {
		return x.MatchId
	}
	return ""
}

func (x *UpdateGoalRequest) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *UpdateGoalRequest) GetScoredAt() string {
	if x != nil {
		return x.ScoredAt
	}
	return ""
}

type GetGoalRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetGoalRequest) Reset() {
	*x = GetGoalRequest{}
	mi := &file_match_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGoalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGoalRequest) ProtoMessage() {}

func (x *GetGoalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGoalRequest.ProtoReflect.Descriptor instead.
func (*GetGoalRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{10}
}

func (x *GetGoalRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteGoalRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteGoalRequest) Reset() {
	*x = DeleteGoalRequest{}
	mi := &file_match_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteGoalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteGoalRequest) ProtoMessage() {}

func (x *DeleteGoalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteGoalRequest.ProtoReflect.Descriptor instead.
func (*DeleteGoalRequest) Descriptor() ([]byte, []int) {
	return file_match_proto_rawDescGZIP(), []int{11}
}

func (x *DeleteGoalRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_match_proto protoreflect.FileDescriptor

const file_match_proto_rawDesc = "" +
	"\n" +
	"\vmatch.proto\x1a\x1bgoogle/protobuf/empty.proto\"\xc3\x01\n" +
	"\x05Match\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04date\x18\x02 \x01(\tR\x04date\x12\x14\n" +
	"\x05venue\x18\x03 \x01(\tR\x05venue\x12 \n" +
	"\fhome_team_id\x18\x04 \x01(\tR\n" +
	"homeTeamId\x12 \n" +
	"\faway_team_id\x18\x05 \x01(\tR\n" +
	"awayTeamId\x12\x1d\n" +
	"\n" +
	"created_at\x18\x06 \x01(\tR\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\a \x01(\tR\tupdatedAt\"\xa9\x01\n" +
	"\x04Goal\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x19\n" +
	"\bmatch_id\x18\x02 \x01(\tR\amatchId\x12\x1b\n" +
	"\tplayer_id\x18\x03 \x01(\tR\bplayerId\x12\x1b\n" +
	"\tscored_at\x18\x04 \x01(\tR\bscoredAt\x12\x1d\n" +
	"\n" +
	"created_at\x18\x05 \x01(\tR\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x06 \x01(\tR\tupdatedAt\"\x82\x01\n" +
	"\x12CreateMatchRequest\x12\x12\n" +
	"\x04date\x18\x01 \x01(\tR\x04date\x12\x14\n" +
	"\x05venue\x18\x02 \x01(\tR\x05venue\x12 \n" +
	"\fhome_team_id\x18\x03 \x01(\tR\n" +
	"homeTeamId\x12 \n" +
	"\faway_team_id\x18\x04 \x01(\tR\n" +
	"awayTeamId\"\x92\x01\n" +
	"\x12UpdateMatchRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04date\x18\x02 \x01(\tR\x04date\x12\x14\n" +
	"\x05venue\x18\x03 \x01(\tR\x05venue\x12 \n" +
	"\fhome_team_id\x18\x04 \x01(\tR\n" +
	"homeTeamId\x12 \n" +
	"\faway_team_id\x18\x05 \x01(\tR\n" +
	"awayTeamId\"$\n" +
	"\x12DeleteMatchRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"!\n" +
	"\x0fGetMatchRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"C\n" +
	"\x10ListMatchRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\x05R\bpageSize\"5\n" +
	"\x11ListMatchResponse\x12 \n" +
	"\amatches\x18\x01 \x03(\v2\x06.MatchR\amatches\"h\n" +
	"\x11CreateGoalRequest\x12\x19\n" +
	"\bmatch_id\x18\x01 \x01(\tR\amatchId\x12\x1b\n" +
	"\tplayer_id\x18\x02 \x01(\tR\bplayerId\x12\x1b\n" +
	"\tscored_at\x18\x03 \x01(\tR\bscoredAt\"x\n" +
	"\x11UpdateGoalRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x19\n" +
	"\bmatch_id\x18\x02 \x01(\tR\amatchId\x12\x1b\n" +
	"\tplayer_id\x18\x03 \x01(\tR\bplayerId\x12\x1b\n" +
	"\tscored_at\x18\x04 \x01(\tR\bscoredAt\" \n" +
	"\x0eGetGoalRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"#\n" +
	"\x11DeleteGoalRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id2\xab\x03\n" +
	"\fMatchService\x12*\n" +
	"\vCreateMatch\x12\x13.CreateMatchRequest\x1a\x06.Match\x12$\n" +
	"\bGetMatch\x12\x10.GetMatchRequest\x1a\x06.Match\x12*\n" +
	"\vUpdateMatch\x12\x13.UpdateMatchRequest\x1a\x06.Match\x12:\n" +
	"\vDeleteMatch\x12\x13.DeleteMatchRequest\x1a\x16.google.protobuf.Empty\x122\n" +
	"\tListMatch\x12\x11.ListMatchRequest\x1a\x12.ListMatchResponse\x12'\n" +
	"\n" +
	"CreateGoal\x12\x12.CreateGoalRequest\x1a\x05.Goal\x128\n" +
	"\n" +
	"DeleteGoal\x12\x12.DeleteGoalRequest\x1a\x16.google.protobuf.Empty\x12'\n" +
	"\n" +
	"UpdateGoal\x12\x12.UpdateGoalRequest\x1a\x05.Goal\x12!\n" +
	"\aGetGoal\x12\x0f.GetGoalRequest\x1a\x05.GoalB?Z=github.com/apriliantocecep/ayo-football/services/match/pkg/pbb\x06proto3"

var (
	file_match_proto_rawDescOnce sync.Once
	file_match_proto_rawDescData []byte
)

func file_match_proto_rawDescGZIP() []byte {
	file_match_proto_rawDescOnce.Do(func() {
		file_match_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_match_proto_rawDesc), len(file_match_proto_rawDesc)))
	})
	return file_match_proto_rawDescData
}

var file_match_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_match_proto_goTypes = []any{
	(*Match)(nil),              // 0: Match
	(*Goal)(nil),               // 1: Goal
	(*CreateMatchRequest)(nil), // 2: CreateMatchRequest
	(*UpdateMatchRequest)(nil), // 3: UpdateMatchRequest
	(*DeleteMatchRequest)(nil), // 4: DeleteMatchRequest
	(*GetMatchRequest)(nil),    // 5: GetMatchRequest
	(*ListMatchRequest)(nil),   // 6: ListMatchRequest
	(*ListMatchResponse)(nil),  // 7: ListMatchResponse
	(*CreateGoalRequest)(nil),  // 8: CreateGoalRequest
	(*UpdateGoalRequest)(nil),  // 9: UpdateGoalRequest
	(*GetGoalRequest)(nil),     // 10: GetGoalRequest
	(*DeleteGoalRequest)(nil),  // 11: DeleteGoalRequest
	(*emptypb.Empty)(nil),      // 12: google.protobuf.Empty
}
var file_match_proto_depIdxs = []int32{
	0,  // 0: ListMatchResponse.matches:type_name -> Match
	2,  // 1: MatchService.CreateMatch:input_type -> CreateMatchRequest
	5,  // 2: MatchService.GetMatch:input_type -> GetMatchRequest
	3,  // 3: MatchService.UpdateMatch:input_type -> UpdateMatchRequest
	4,  // 4: MatchService.DeleteMatch:input_type -> DeleteMatchRequest
	6,  // 5: MatchService.ListMatch:input_type -> ListMatchRequest
	8,  // 6: MatchService.CreateGoal:input_type -> CreateGoalRequest
	11, // 7: MatchService.DeleteGoal:input_type -> DeleteGoalRequest
	9,  // 8: MatchService.UpdateGoal:input_type -> UpdateGoalRequest
	10, // 9: MatchService.GetGoal:input_type -> GetGoalRequest
	0,  // 10: MatchService.CreateMatch:output_type -> Match
	0,  // 11: MatchService.GetMatch:output_type -> Match
	0,  // 12: MatchService.UpdateMatch:output_type -> Match
	12, // 13: MatchService.DeleteMatch:output_type -> google.protobuf.Empty
	7,  // 14: MatchService.ListMatch:output_type -> ListMatchResponse
	1,  // 15: MatchService.CreateGoal:output_type -> Goal
	12, // 16: MatchService.DeleteGoal:output_type -> google.protobuf.Empty
	1,  // 17: MatchService.UpdateGoal:output_type -> Goal
	1,  // 18: MatchService.GetGoal:output_type -> Goal
	10, // [10:19] is the sub-list for method output_type
	1,  // [1:10] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_match_proto_init() }
func file_match_proto_init() {
	if File_match_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_match_proto_rawDesc), len(file_match_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_match_proto_goTypes,
		DependencyIndexes: file_match_proto_depIdxs,
		MessageInfos:      file_match_proto_msgTypes,
	}.Build()
	File_match_proto = out.File
	file_match_proto_goTypes = nil
	file_match_proto_depIdxs = nil
}
