// Package etcd provides an etcd-backed registry implementation.
package etcd

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

const (
	logPrefix          = "[etcd-registry]"
	defaultDialTimeout = 5 * time.Second
)

var defaultEndpoints = []string{"localhost:2379"}

// Option is a function that configures a registry option.
type Option func(*options)

type options struct {
	dialTimeout time.Duration
	endpoints   []string
}

// WithEndpoints sets the etcd endpoints for the registry.
func WithEndpoints(endpoints []string) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

// WithDialTimeout sets the dial timeout for etcd client.
func WithDialTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.dialTimeout = timeout
	}
}

// NewRegistry creates a new Registry instance with the given etcd endpoints and options.
func NewRegistry(opts ...Option) (*Registry, error) {
	op := &options{
		dialTimeout: defaultDialTimeout,
		endpoints:   defaultEndpoints,
	}
	for _, apply := range opts {
		apply(op)
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   op.endpoints,
		DialTimeout: op.dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &Registry{
		client: client,
	}, nil
}

// Registry implements the Registry interface using etcd.
type Registry struct {
	client  *clientv3.Client
	manager endpoints.Manager
}

// Register registers the service address with etcd.
func (r *Registry) Register(ctx context.Context, serviceName string, addr string) error {
	manager, err := endpoints.NewManager(r.client, serviceName)
	if err != nil {
		return err
	}
	r.manager = manager

	serviceID := serviceName + "/" + addr
	if err := r.manager.AddEndpoint(ctx, serviceID, endpoints.Endpoint{Addr: addr}); err != nil {
		log.Printf("%s register error: %s", logPrefix, err.Error())
		return err
	}

	return nil
}

// Deregister removes the service address from etcd.
func (r *Registry) Deregister(ctx context.Context, serviceName string, addr string) error {
	if r.manager == nil {
		return nil
	}

	serviceID := serviceName + "/" + addr
	if err := r.manager.DeleteEndpoint(ctx, serviceID); err != nil {
		log.Printf("%s deregister error: %s", logPrefix, err.Error())
		return err
	}

	return nil
}
