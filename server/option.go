// Package server provides gRPC server and related configurations.
package server

import (
	"github.com/yafeng-Soong/aokiji/registry"
	"google.golang.org/grpc"
)

type serverOption struct {
	serviceName string
	grpcServer  *grpc.Server
	registry    registry.Registry
}

// Option is a function that configures a serverOption.
type Option func(*serverOption)

// WithServiceName sets the service name for the server.
func WithServiceName(serviceName string) Option {
	return func(s *serverOption) {
		s.serviceName = serviceName
	}
}

// WithGRPCServer sets the gRPC server for the server.
func WithGRPCServer(grpcServer *grpc.Server) Option {
	return func(s *serverOption) {
		s.grpcServer = grpcServer
	}
}

// WithRegistry sets the registry for the server.
func WithRegistry(registry registry.Registry) Option {
	return func(s *serverOption) {
		s.registry = registry
	}
}
