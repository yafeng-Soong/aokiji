// Package registry provides service registration and discovery interfaces.
package registry

import (
	"context"
	"path"

	"github.com/yafeng-Soong/aokiji/registry/internal/etcd"
)

const prefix = "aokiji/services"

// Registry is the interface for a service registry.
type Registry interface {
	Register(ctx context.Context, addr string) error
	Deregister(ctx context.Context) error
}

// NewEmptyRegistryBuilder creates an empty Registry builder.
func NewEmptyRegistryBuilder(_ Config) (Registry, error) {
	return &emptyRegistry{}, nil
}

type emptyRegistry struct {
}

func (r *emptyRegistry) Register(_ context.Context, _ string) error {
	return nil
}

func (r *emptyRegistry) Deregister(_ context.Context) error {
	return nil
}

// Config holds the configuration for a Registry.
type Config struct {
	ServiceName string
}

// NewEtcdRegistryBuilder creates an EtcdRegistry builder with the given service name and etcd endpoints.
func NewEtcdRegistryBuilder(endpoints []string) Builder {
	// serviceName := path.Join(prefix, config.ServiceName)
	// return etcd.NewRegistry(serviceName, config.RegistryEndpoints)
	return func(config Config) (Registry, error) {
		serviceName := path.Join(prefix, config.ServiceName)
		return etcd.NewRegistry(serviceName, endpoints)
	}
}

// Builder is a function type that builds a Registry based on the given configuration.
type Builder func(config Config) (Registry, error)
