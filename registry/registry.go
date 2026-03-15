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
func NewEmptyRegistryBuilder(config RegistryConfig) (Registry, error) {
	return &emptyRegistry{}, nil
}

type emptyRegistry struct {
}

func (r *emptyRegistry) Register(ctx context.Context, addr string) error {
	return nil
}

func (r *emptyRegistry) Deregister(ctx context.Context) error {
	return nil
}

// RegistryConfig holds the configuration for a Registry.
type RegistryConfig struct {
	ServiceName string
}

// NewEtcdRegistryBuilder creates an EtcdRegistry builder with the given service name and etcd endpoints.
func NewEtcdRegistryBuilder(endpoints []string) RegistryBuilder {
	// serviceName := path.Join(prefix, config.ServiceName)
	// return etcd.NewRegistry(serviceName, config.RegistryEndpoints)
	return func(config RegistryConfig) (Registry, error) {
		serviceName := path.Join(prefix, config.ServiceName)
		return etcd.NewRegistry(serviceName, endpoints)
	}
}

// RegistryBuilder is a function type that builds a Registry based on the given configuration.
type RegistryBuilder func(config RegistryConfig) (Registry, error)
