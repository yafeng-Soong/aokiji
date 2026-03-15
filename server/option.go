package server

import (
	"github.com/yafeng-Soong/aokiji/registry"

	"google.golang.org/grpc"
)

var emptyRegistryBuilder = registryBuilder{
	build: registry.NewEmptyRegistryBuilder,
}

type serverOption struct {
	serviceName     string
	grpcServer      *grpc.Server
	registryBuilder registryBuilder
}

type registryBuilder struct {
	build registry.RegistryBuilder
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

// WithEtcdRegistry sets the etcd registry for the server.
func WithEtcdRegistry(endpoints []string) Option {
	return func(s *serverOption) {
		s.registryBuilder = registryBuilder{
			build: registry.NewEtcdRegistryBuilder(endpoints),
		}
	}
}
