package util

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// FakeGRPCServerStream implements google.golang.org/grpc.ServerStream
// interface for unit tests.
type FakeGRPCServerStream struct {
	OnSetHeader  func(metadata.MD) error
	OnSendHeader func(metadata.MD) error
	OnSetTrailer func(m metadata.MD)
	OnContext    func() context.Context
	OnSendMsg    func(m interface{}) error
	OnRecvMsg    func(m interface{}) error
}

// SetHeader implements grpc.ServerStream.SetHeader.
func (s *FakeGRPCServerStream) SetHeader(m metadata.MD) error {
	if s.OnSetHeader != nil {
		return s.OnSetHeader(m)
	}
	panic("OnSetHeader not set")
}

// SendHeader implements grpc.ServerStream.SendHeader.
func (s *FakeGRPCServerStream) SendHeader(m metadata.MD) error {
	if s.OnSendHeader != nil {
		return s.OnSendHeader(m)
	}
	panic("OnSendHeader not set")
}

// SetTrailer implements grpc.ServerStream.SetTrailer.
func (s *FakeGRPCServerStream) SetTrailer(m metadata.MD) {
	if s.OnSetTrailer != nil {
		s.OnSetTrailer(m)
	}
	panic("OnSetTrailer not set")
}

// Context implements grpc.ServerStream.Context.
func (s *FakeGRPCServerStream) Context() context.Context {
	if s.OnContext != nil {
		return s.OnContext()
	}
	panic("OnContext not set")
}

// SendMsg implements grpc.ServerStream.SendMsg.
func (s *FakeGRPCServerStream) SendMsg(m interface{}) error {
	if s.OnSendMsg != nil {
		return s.OnSendMsg(m)
	}
	panic("OnSendMsg not set")
}

// RecvMsg implements grpc.ServerStream.RecvMsg.
func (s *FakeGRPCServerStream) RecvMsg(m interface{}) error {
	if s.OnRecvMsg != nil {
		return s.OnRecvMsg(m)
	}
	panic("OnRecvMsg not set")
}
