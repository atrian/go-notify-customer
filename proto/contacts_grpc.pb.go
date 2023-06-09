// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: proto/contacts.proto

package go_notify_client

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// VaultClient is the client API for Vault service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VaultClient interface {
	GetContacts(ctx context.Context, in *GetContactsRequest, opts ...grpc.CallOption) (*GetContactsResponse, error)
}

type vaultClient struct {
	cc grpc.ClientConnInterface
}

func NewVaultClient(cc grpc.ClientConnInterface) VaultClient {
	return &vaultClient{cc}
}

func (c *vaultClient) GetContacts(ctx context.Context, in *GetContactsRequest, opts ...grpc.CallOption) (*GetContactsResponse, error) {
	out := new(GetContactsResponse)
	err := c.cc.Invoke(ctx, "/contacts.Vault/GetContacts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VaultServer is the server API for Vault service.
// All implementations must embed UnimplementedVaultServer
// for forward compatibility
type VaultServer interface {
	GetContacts(context.Context, *GetContactsRequest) (*GetContactsResponse, error)
	mustEmbedUnimplementedVaultServer()
}

// UnimplementedVaultServer must be embedded to have forward compatible implementations.
type UnimplementedVaultServer struct {
}

func (UnimplementedVaultServer) GetContacts(context.Context, *GetContactsRequest) (*GetContactsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContacts not implemented")
}
func (UnimplementedVaultServer) mustEmbedUnimplementedVaultServer() {}

// UnsafeVaultServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VaultServer will
// result in compilation errors.
type UnsafeVaultServer interface {
	mustEmbedUnimplementedVaultServer()
}

func RegisterVaultServer(s grpc.ServiceRegistrar, srv VaultServer) {
	s.RegisterService(&Vault_ServiceDesc, srv)
}

func _Vault_GetContacts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetContactsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServer).GetContacts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contacts.Vault/GetContacts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServer).GetContacts(ctx, req.(*GetContactsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Vault_ServiceDesc is the grpc.ServiceDesc for Vault service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Vault_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "contacts.Vault",
	HandlerType: (*VaultServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetContacts",
			Handler:    _Vault_GetContacts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/contacts.proto",
}
