// Package registry provides service registration and discovery interfaces.
package registry

import (
	"context"
)

// Prefix is the key prefix used for service registration in the registry.
const Prefix = "aokiji/services"

// Registry is the interface for a service registry.
type Registry interface {
	Register(context.Context, string, string) error
	Deregister(context.Context, string, string) error
}

// NewEmptyRegistry creates a new Registry instance that does nothing.
func NewEmptyRegistry() Registry {
	return &emptyRegistry{}
}

type emptyRegistry struct {
}

func (r *emptyRegistry) Register(_ context.Context, _, _ string) error {
	return nil
}

func (r *emptyRegistry) Deregister(_ context.Context, _, _ string) error {
	return nil
}
