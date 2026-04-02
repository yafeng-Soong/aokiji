// Package resolver provides a gresolver.Builder and gresolver.Resolver for etcd.
package resolver

import (
	"fmt"
	"path"

	"github.com/yafeng-Soong/aokiji/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	gresolver "google.golang.org/grpc/resolver"
)

// NewEtcdBuilder creates a new gresolver.Builder for etcd.
func NewEtcdBuilder(cfg clientv3.Config) (gresolver.Builder, error) {
	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %v", err)
	}

	builder, err := resolver.NewBuilder(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create resolver builder for etcd: %v", err)
	}

	return &etcdBuilder{c: c, b: builder}, nil
}

type etcdBuilder struct {
	c *clientv3.Client
	b gresolver.Builder
}

// Build implements gresolver.Builder.Build.
func (b *etcdBuilder) Build(
	target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions,
) (gresolver.Resolver, error) {
	target.URL.Path = path.Join(registry.Prefix, target.URL.Path)
	r, err := b.b.Build(target, cc, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to build resolver for etcd: %v", err)
	}

	return &etcdResolver{r: r}, nil
}

// Scheme implements gresolver.Builder.Scheme.
func (b *etcdBuilder) Scheme() string {
	return "aokiji-etcd"
}

type etcdResolver struct {
	r gresolver.Resolver
}

// ResolveNow implements gresolver.Resolver.ResolveNow.
func (r *etcdResolver) ResolveNow(opts gresolver.ResolveNowOptions) {
	r.r.ResolveNow(opts)
}

// Close implements gresolver.Resolver.Close.
func (r *etcdResolver) Close() {
	r.r.Close()
}
