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

// NewRegistry creates a new Registry instance with the given service name and client configuration.
func NewRegistry(serviceName string, etcdEndpoints []string) (*Registry, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: defaultDialTimeout,
	})
	if err != nil {
		return nil, err
	}

	manager, err := endpoints.NewManager(client, serviceName)
	if err != nil {
		return nil, err
	}

	return &Registry{
		client:      client,
		manager:     manager,
		serviceName: serviceName,
	}, nil
}

// Registry implements the Registry interface using etcd.
type Registry struct {
	serviceID   string
	serviceName string
	client      *clientv3.Client
	manager     endpoints.Manager
}

// Register registers the service address with etcd.
func (r *Registry) Register(ctx context.Context, addr string) error {
	r.serviceID = r.serviceName + "/" + addr
	if err := r.manager.AddEndpoint(ctx, r.serviceID, endpoints.Endpoint{Addr: addr}); err != nil {
		log.Fatalf("%s register error: %s", logPrefix, err.Error())
	}

	return nil
}

// Deregister removes the service address from etcd.
func (r *Registry) Deregister(ctx context.Context) error {
	if err := r.manager.DeleteEndpoint(ctx, r.serviceID); err != nil {
		log.Fatalf("%s deregister error: %s", logPrefix, err.Error())
	}

	return r.client.Close()
}
