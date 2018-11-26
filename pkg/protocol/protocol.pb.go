// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol.proto

/*
Package protocol is a generated protocol buffer package.

It is generated from these files:
	protocol.proto

It has these top-level messages:
	SignupRequest
	SignupReply
	LoginRequest
	LoginReply
	FollowUnFollowRequest
	FollowUnFollowReply
	AddTweetRequest
	AddTweetReply
	GetFollowingTweetsRequest
	GetFollowingTweetsReply
	Tweet
	Timestamp
	GetUserProfileRequest
	GetUserProfileReply
*/
package protocol

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SignupRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *SignupRequest) Reset()                    { *m = SignupRequest{} }
func (m *SignupRequest) String() string            { return proto.CompactTextString(m) }
func (*SignupRequest) ProtoMessage()               {}
func (*SignupRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SignupRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *SignupRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type SignupReply struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Success  bool   `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
}

func (m *SignupReply) Reset()                    { *m = SignupReply{} }
func (m *SignupReply) String() string            { return proto.CompactTextString(m) }
func (*SignupReply) ProtoMessage()               {}
func (*SignupReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SignupReply) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *SignupReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type LoginRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *LoginRequest) Reset()                    { *m = LoginRequest{} }
func (m *LoginRequest) String() string            { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()               {}
func (*LoginRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *LoginRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LoginReply struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Success  bool   `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
}

func (m *LoginReply) Reset()                    { *m = LoginReply{} }
func (m *LoginReply) String() string            { return proto.CompactTextString(m) }
func (*LoginReply) ProtoMessage()               {}
func (*LoginReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LoginReply) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type FollowUnFollowRequest struct {
	Username   string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Targetname string `protobuf:"bytes,2,opt,name=targetname" json:"targetname,omitempty"`
}

func (m *FollowUnFollowRequest) Reset()                    { *m = FollowUnFollowRequest{} }
func (m *FollowUnFollowRequest) String() string            { return proto.CompactTextString(m) }
func (*FollowUnFollowRequest) ProtoMessage()               {}
func (*FollowUnFollowRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *FollowUnFollowRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *FollowUnFollowRequest) GetTargetname() string {
	if m != nil {
		return m.Targetname
	}
	return ""
}

type FollowUnFollowReply struct {
	Username   string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Targetname string `protobuf:"bytes,2,opt,name=targetname" json:"targetname,omitempty"`
	Success    bool   `protobuf:"varint,3,opt,name=success" json:"success,omitempty"`
}

func (m *FollowUnFollowReply) Reset()                    { *m = FollowUnFollowReply{} }
func (m *FollowUnFollowReply) String() string            { return proto.CompactTextString(m) }
func (*FollowUnFollowReply) ProtoMessage()               {}
func (*FollowUnFollowReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *FollowUnFollowReply) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *FollowUnFollowReply) GetTargetname() string {
	if m != nil {
		return m.Targetname
	}
	return ""
}

func (m *FollowUnFollowReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type AddTweetRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Post     string `protobuf:"bytes,2,opt,name=post" json:"post,omitempty"`
}

func (m *AddTweetRequest) Reset()                    { *m = AddTweetRequest{} }
func (m *AddTweetRequest) String() string            { return proto.CompactTextString(m) }
func (*AddTweetRequest) ProtoMessage()               {}
func (*AddTweetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *AddTweetRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AddTweetRequest) GetPost() string {
	if m != nil {
		return m.Post
	}
	return ""
}

type AddTweetReply struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Success  bool   `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
}

func (m *AddTweetReply) Reset()                    { *m = AddTweetReply{} }
func (m *AddTweetReply) String() string            { return proto.CompactTextString(m) }
func (*AddTweetReply) ProtoMessage()               {}
func (*AddTweetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *AddTweetReply) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AddTweetReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type GetFollowingTweetsRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
}

func (m *GetFollowingTweetsRequest) Reset()                    { *m = GetFollowingTweetsRequest{} }
func (m *GetFollowingTweetsRequest) String() string            { return proto.CompactTextString(m) }
func (*GetFollowingTweetsRequest) ProtoMessage()               {}
func (*GetFollowingTweetsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *GetFollowingTweetsRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

type GetFollowingTweetsReply struct {
	Username string   `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Tweet    []*Tweet `protobuf:"bytes,2,rep,name=tweet" json:"tweet,omitempty"`
	Success  bool     `protobuf:"varint,3,opt,name=success" json:"success,omitempty"`
}

func (m *GetFollowingTweetsReply) Reset()                    { *m = GetFollowingTweetsReply{} }
func (m *GetFollowingTweetsReply) String() string            { return proto.CompactTextString(m) }
func (*GetFollowingTweetsReply) ProtoMessage()               {}
func (*GetFollowingTweetsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *GetFollowingTweetsReply) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *GetFollowingTweetsReply) GetTweet() []*Tweet {
	if m != nil {
		return m.Tweet
	}
	return nil
}

func (m *GetFollowingTweetsReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type Tweet struct {
	UserName  string     `protobuf:"bytes,1,opt,name=UserName,json=userName" json:"UserName,omitempty"`
	Timestamp *Timestamp `protobuf:"bytes,2,opt,name=Timestamp,json=timestamp" json:"Timestamp,omitempty"`
	Body      string     `protobuf:"bytes,3,opt,name=Body,json=body" json:"Body,omitempty"`
}

func (m *Tweet) Reset()                    { *m = Tweet{} }
func (m *Tweet) String() string            { return proto.CompactTextString(m) }
func (*Tweet) ProtoMessage()               {}
func (*Tweet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Tweet) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *Tweet) GetTimestamp() *Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Tweet) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

type Timestamp struct {
	// Represents seconds of UTC time since Unix epoch
	// 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
	// 9999-12-31T23:59:59Z inclusive.
	Seconds int64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
	// Non-negative fractions of a second at nanosecond resolution. Negative
	// second values with fractions must still have non-negative nanos values
	// that count forward in time. Must be from 0 to 999,999,999
	// inclusive.
	Nanos int32 `protobuf:"varint,2,opt,name=nanos" json:"nanos,omitempty"`
}

func (m *Timestamp) Reset()                    { *m = Timestamp{} }
func (m *Timestamp) String() string            { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()               {}
func (*Timestamp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Timestamp) GetSeconds() int64 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func (m *Timestamp) GetNanos() int32 {
	if m != nil {
		return m.Nanos
	}
	return 0
}

type GetUserProfileRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
}

func (m *GetUserProfileRequest) Reset()                    { *m = GetUserProfileRequest{} }
func (m *GetUserProfileRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserProfileRequest) ProtoMessage()               {}
func (*GetUserProfileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *GetUserProfileRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

type GetUserProfileReply struct {
	Username      string   `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	TweetList     []*Tweet `protobuf:"bytes,2,rep,name=tweetList" json:"tweetList,omitempty"`
	FollowingList []string `protobuf:"bytes,3,rep,name=followingList" json:"followingList,omitempty"`
	FollowerList  []string `protobuf:"bytes,4,rep,name=followerList" json:"followerList,omitempty"`
	Success       bool     `protobuf:"varint,5,opt,name=success" json:"success,omitempty"`
}

func (m *GetUserProfileReply) Reset()                    { *m = GetUserProfileReply{} }
func (m *GetUserProfileReply) String() string            { return proto.CompactTextString(m) }
func (*GetUserProfileReply) ProtoMessage()               {}
func (*GetUserProfileReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *GetUserProfileReply) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *GetUserProfileReply) GetTweetList() []*Tweet {
	if m != nil {
		return m.TweetList
	}
	return nil
}

func (m *GetUserProfileReply) GetFollowingList() []string {
	if m != nil {
		return m.FollowingList
	}
	return nil
}

func (m *GetUserProfileReply) GetFollowerList() []string {
	if m != nil {
		return m.FollowerList
	}
	return nil
}

func (m *GetUserProfileReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterType((*SignupRequest)(nil), "protocol.SignupRequest")
	proto.RegisterType((*SignupReply)(nil), "protocol.SignupReply")
	proto.RegisterType((*LoginRequest)(nil), "protocol.LoginRequest")
	proto.RegisterType((*LoginReply)(nil), "protocol.LoginReply")
	proto.RegisterType((*FollowUnFollowRequest)(nil), "protocol.FollowUnFollowRequest")
	proto.RegisterType((*FollowUnFollowReply)(nil), "protocol.FollowUnFollowReply")
	proto.RegisterType((*AddTweetRequest)(nil), "protocol.AddTweetRequest")
	proto.RegisterType((*AddTweetReply)(nil), "protocol.AddTweetReply")
	proto.RegisterType((*GetFollowingTweetsRequest)(nil), "protocol.GetFollowingTweetsRequest")
	proto.RegisterType((*GetFollowingTweetsReply)(nil), "protocol.GetFollowingTweetsReply")
	proto.RegisterType((*Tweet)(nil), "protocol.Tweet")
	proto.RegisterType((*Timestamp)(nil), "protocol.Timestamp")
	proto.RegisterType((*GetUserProfileRequest)(nil), "protocol.GetUserProfileRequest")
	proto.RegisterType((*GetUserProfileReply)(nil), "protocol.GetUserProfileReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for TwitterRPC service

type TwitterRPCClient interface {
	SignupRPC(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupReply, error)
	LoginRPC(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
	FollowUnFollowRPC(ctx context.Context, in *FollowUnFollowRequest, opts ...grpc.CallOption) (*FollowUnFollowReply, error)
	AddTweetRPC(ctx context.Context, in *AddTweetRequest, opts ...grpc.CallOption) (*AddTweetReply, error)
	GetFollowingTweetsRPC(ctx context.Context, in *GetFollowingTweetsRequest, opts ...grpc.CallOption) (*GetFollowingTweetsReply, error)
	GetUserProfileRPC(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*GetUserProfileReply, error)
}

type twitterRPCClient struct {
	cc *grpc.ClientConn
}

func NewTwitterRPCClient(cc *grpc.ClientConn) TwitterRPCClient {
	return &twitterRPCClient{cc}
}

func (c *twitterRPCClient) SignupRPC(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupReply, error) {
	out := new(SignupReply)
	err := grpc.Invoke(ctx, "/protocol.TwitterRPC/SignupRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *twitterRPCClient) LoginRPC(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := grpc.Invoke(ctx, "/protocol.TwitterRPC/LoginRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *twitterRPCClient) FollowUnFollowRPC(ctx context.Context, in *FollowUnFollowRequest, opts ...grpc.CallOption) (*FollowUnFollowReply, error) {
	out := new(FollowUnFollowReply)
	err := grpc.Invoke(ctx, "/protocol.TwitterRPC/FollowUnFollowRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *twitterRPCClient) AddTweetRPC(ctx context.Context, in *AddTweetRequest, opts ...grpc.CallOption) (*AddTweetReply, error) {
	out := new(AddTweetReply)
	err := grpc.Invoke(ctx, "/protocol.TwitterRPC/AddTweetRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *twitterRPCClient) GetFollowingTweetsRPC(ctx context.Context, in *GetFollowingTweetsRequest, opts ...grpc.CallOption) (*GetFollowingTweetsReply, error) {
	out := new(GetFollowingTweetsReply)
	err := grpc.Invoke(ctx, "/protocol.TwitterRPC/GetFollowingTweetsRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *twitterRPCClient) GetUserProfileRPC(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*GetUserProfileReply, error) {
	out := new(GetUserProfileReply)
	err := grpc.Invoke(ctx, "/protocol.TwitterRPC/GetUserProfileRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TwitterRPC service

type TwitterRPCServer interface {
	SignupRPC(context.Context, *SignupRequest) (*SignupReply, error)
	LoginRPC(context.Context, *LoginRequest) (*LoginReply, error)
	FollowUnFollowRPC(context.Context, *FollowUnFollowRequest) (*FollowUnFollowReply, error)
	AddTweetRPC(context.Context, *AddTweetRequest) (*AddTweetReply, error)
	GetFollowingTweetsRPC(context.Context, *GetFollowingTweetsRequest) (*GetFollowingTweetsReply, error)
	GetUserProfileRPC(context.Context, *GetUserProfileRequest) (*GetUserProfileReply, error)
}

func RegisterTwitterRPCServer(s *grpc.Server, srv TwitterRPCServer) {
	s.RegisterService(&_TwitterRPC_serviceDesc, srv)
}

func _TwitterRPC_SignupRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TwitterRPCServer).SignupRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.TwitterRPC/SignupRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TwitterRPCServer).SignupRPC(ctx, req.(*SignupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TwitterRPC_LoginRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TwitterRPCServer).LoginRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.TwitterRPC/LoginRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TwitterRPCServer).LoginRPC(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TwitterRPC_FollowUnFollowRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowUnFollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TwitterRPCServer).FollowUnFollowRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.TwitterRPC/FollowUnFollowRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TwitterRPCServer).FollowUnFollowRPC(ctx, req.(*FollowUnFollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TwitterRPC_AddTweetRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTweetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TwitterRPCServer).AddTweetRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.TwitterRPC/AddTweetRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TwitterRPCServer).AddTweetRPC(ctx, req.(*AddTweetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TwitterRPC_GetFollowingTweetsRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFollowingTweetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TwitterRPCServer).GetFollowingTweetsRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.TwitterRPC/GetFollowingTweetsRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TwitterRPCServer).GetFollowingTweetsRPC(ctx, req.(*GetFollowingTweetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TwitterRPC_GetUserProfileRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TwitterRPCServer).GetUserProfileRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.TwitterRPC/GetUserProfileRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TwitterRPCServer).GetUserProfileRPC(ctx, req.(*GetUserProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TwitterRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.TwitterRPC",
	HandlerType: (*TwitterRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignupRPC",
			Handler:    _TwitterRPC_SignupRPC_Handler,
		},
		{
			MethodName: "LoginRPC",
			Handler:    _TwitterRPC_LoginRPC_Handler,
		},
		{
			MethodName: "FollowUnFollowRPC",
			Handler:    _TwitterRPC_FollowUnFollowRPC_Handler,
		},
		{
			MethodName: "AddTweetRPC",
			Handler:    _TwitterRPC_AddTweetRPC_Handler,
		},
		{
			MethodName: "GetFollowingTweetsRPC",
			Handler:    _TwitterRPC_GetFollowingTweetsRPC_Handler,
		},
		{
			MethodName: "GetUserProfileRPC",
			Handler:    _TwitterRPC_GetUserProfileRPC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protocol.proto",
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 550 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x6e, 0xea, 0x18, 0xe2, 0x49, 0x43, 0xc5, 0xa6, 0x21, 0xa9, 0x25, 0x20, 0x2c, 0x20, 0xe5,
	0x42, 0x25, 0xd2, 0x03, 0x12, 0x88, 0x43, 0x6b, 0xd1, 0x5c, 0x2a, 0x64, 0x6d, 0x93, 0x07, 0x70,
	0xe3, 0x6d, 0x64, 0x70, 0xbc, 0xc6, 0xbb, 0x21, 0xca, 0xfb, 0xf1, 0x1a, 0xbc, 0x0b, 0xf2, 0x6e,
	0x1c, 0xff, 0xd4, 0x49, 0xac, 0x70, 0xb2, 0x67, 0xe7, 0xdb, 0x6f, 0xbe, 0x99, 0x9d, 0x19, 0x78,
	0x16, 0x46, 0x4c, 0xb0, 0x29, 0xf3, 0x2f, 0xe4, 0x0f, 0x6a, 0x24, 0x36, 0x1e, 0x41, 0xeb, 0xce,
	0x9b, 0x05, 0x8b, 0x90, 0xd0, 0x5f, 0x0b, 0xca, 0x05, 0x32, 0xa1, 0xb1, 0xe0, 0x34, 0x0a, 0x9c,
	0x39, 0xed, 0xd5, 0xfa, 0xb5, 0x81, 0x41, 0x36, 0x76, 0xec, 0x0b, 0x1d, 0xce, 0x97, 0x2c, 0x72,
	0x7b, 0xc7, 0xca, 0x97, 0xd8, 0xd8, 0x82, 0x66, 0x42, 0x14, 0xfa, 0xab, 0x9d, 0x34, 0x3d, 0x78,
	0xca, 0x17, 0xd3, 0x29, 0xe5, 0x5c, 0xb2, 0x34, 0x48, 0x62, 0xe2, 0x1b, 0x38, 0xb9, 0x65, 0x33,
	0x2f, 0xf8, 0x5f, 0x31, 0xd7, 0x00, 0x6b, 0x9e, 0xc3, 0xb5, 0xdc, 0x41, 0xe7, 0x86, 0xf9, 0x3e,
	0x5b, 0x4e, 0x02, 0xf5, 0xad, 0x22, 0xea, 0x15, 0x80, 0x70, 0xa2, 0x19, 0x15, 0xd2, 0xab, 0x64,
	0x65, 0x4e, 0xf0, 0x4f, 0x68, 0x17, 0x49, 0xf7, 0x29, 0xdc, 0x43, 0x99, 0xcd, 0x40, 0xcb, 0x67,
	0x70, 0x05, 0xa7, 0x57, 0xae, 0x3b, 0x5e, 0x52, 0x2a, 0xaa, 0x68, 0x47, 0x50, 0x0f, 0x19, 0x17,
	0xeb, 0x10, 0xf2, 0x1f, 0x7f, 0x83, 0x56, 0x4a, 0x71, 0x78, 0x2d, 0x3f, 0xc1, 0xf9, 0x88, 0x0a,
	0x95, 0xb1, 0x17, 0xcc, 0x24, 0x1f, 0xaf, 0xa0, 0x09, 0xff, 0x86, 0x6e, 0xd9, 0xc5, 0x7d, 0x4a,
	0xde, 0x83, 0x2e, 0x62, 0x68, 0xef, 0xb8, 0xaf, 0x0d, 0x9a, 0xc3, 0xd3, 0x8b, 0x4d, 0xff, 0xab,
	0x54, 0x94, 0x77, 0x47, 0xe9, 0x7e, 0x80, 0x2e, 0x91, 0x71, 0x94, 0x09, 0xa7, 0xd1, 0xf7, 0x42,
	0x94, 0xd8, 0x46, 0x1f, 0xc1, 0x18, 0x7b, 0x73, 0xca, 0x85, 0x33, 0x0f, 0x65, 0xc6, 0xcd, 0x61,
	0x3b, 0x13, 0x29, 0x71, 0x11, 0x43, 0x24, 0xbf, 0x71, 0x8d, 0xaf, 0x99, 0xbb, 0x92, 0xe1, 0x0c,
	0x52, 0xbf, 0x67, 0xee, 0x0a, 0x7f, 0xc9, 0xd0, 0x48, 0x49, 0x74, 0xca, 0x02, 0x97, 0xcb, 0x70,
	0x1a, 0x49, 0x4c, 0x74, 0x06, 0x7a, 0xe0, 0x04, 0x4c, 0xd5, 0x56, 0x27, 0xca, 0xc0, 0x97, 0xd0,
	0x19, 0x51, 0x11, 0x4b, 0xb4, 0x23, 0xf6, 0xe0, 0xf9, 0xb4, 0x4a, 0x55, 0xff, 0xd4, 0xa0, 0x5d,
	0xbc, 0xb5, 0xaf, 0xa4, 0x1f, 0xc0, 0x90, 0x45, 0xbb, 0xf5, 0xf8, 0xd6, 0xb2, 0xa6, 0x08, 0xf4,
	0x0e, 0x5a, 0x0f, 0xc9, 0xab, 0xc9, 0x2b, 0x5a, 0x5f, 0x1b, 0x18, 0x24, 0x7f, 0x88, 0x30, 0x9c,
	0xa8, 0x03, 0x1a, 0x49, 0x50, 0x5d, 0x82, 0x72, 0x67, 0xd9, 0x47, 0xd2, 0x73, 0x8f, 0x34, 0xfc,
	0xab, 0x01, 0x8c, 0x97, 0x9e, 0x10, 0x34, 0x22, 0xb6, 0x85, 0xbe, 0x82, 0xb1, 0xde, 0x40, 0xb6,
	0x85, 0xba, 0xa9, 0xb6, 0xdc, 0x7e, 0x33, 0x3b, 0x8f, 0x1d, 0xa1, 0xbf, 0xc2, 0x47, 0xe8, 0x33,
	0x34, 0xd4, 0xce, 0xb0, 0x2d, 0xf4, 0x22, 0x05, 0x65, 0xf7, 0x91, 0x79, 0xf6, 0xe8, 0x5c, 0xdd,
	0x9d, 0xc0, 0xf3, 0xc2, 0x58, 0xdb, 0x16, 0x7a, 0x9d, 0x82, 0x4b, 0x17, 0x89, 0xf9, 0x72, 0x3b,
	0x40, 0xd1, 0x5a, 0xd0, 0xdc, 0x4c, 0x9f, 0x6d, 0xa1, 0xf3, 0x14, 0x5f, 0x98, 0x6b, 0xb3, 0x5b,
	0xe6, 0x52, 0x24, 0x8e, 0xec, 0x90, 0xe2, 0x08, 0xd9, 0x16, 0x7a, 0x9b, 0xde, 0xd9, 0x3a, 0x9c,
	0xe6, 0x9b, 0xdd, 0xa0, 0x4d, 0xfa, 0x85, 0x76, 0xca, 0xa7, 0x5f, 0xda, 0xa1, 0xd9, 0xf4, 0x4b,
	0x9a, 0x11, 0x1f, 0xdd, 0x3f, 0x91, 0xfe, 0xcb, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xfd, 0xf5,
	0x8d, 0x6f, 0xbe, 0x06, 0x00, 0x00,
}
