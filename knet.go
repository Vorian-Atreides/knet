// Package knet implement networking functions when running in Kuberntes.
package knet

import (
	"context"
	"fmt"
	"net"
)

const (
	defaultNamespace = "default"
)

var (
	defaultResolver netResolver = net.DefaultResolver
)

type netResolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
}

// DefaultResolver is the resolver used by the package-level Lookup functions.
var DefaultResolver = &Resolver{
	resolver: net.DefaultResolver,
}

// Resolver looks up names and numbers.
type Resolver struct {
	resolver netResolver
}

// LookupHost looks up the given host using the local resolver.
// It returns a slice of that host's addresses.
//
// If no namespace is provided, the function will use the default kubernetes
// namespace.
func LookupHost(serviceName string, namespace string) ([]string, error) {
	return DefaultResolver.LookupHost(context.Background(), serviceName, namespace)
}

// LookupHost looks up the given host using the local resolver.
// It returns a slice of that host's addresses.
//
// If no namespace is provided, the function will use the default kubernetes
// namespace.
func (r *Resolver) LookupHost(ctx context.Context, serviceName string, namespace string) ([]string, error) {
	_, srvs, err := r.LookupSRV(ctx, serviceName, namespace)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, srv := range srvs {
		hosts, err := defaultResolver.LookupHost(ctx, srv.Target)
		if err != nil {
			l.Warningf("failed dns lookup due to: %v.\n", err)
			continue
		}
		addrs = append(addrs, hosts...)
	}
	return addrs, nil
}

// LookupSRV constructs the DNS name with the format used in Kubernetes.
// The returned records are sorted by priority and randomized by weight within
// a priority.
//
// If no namespace is provided, the function will use the default kubernetes
// namespace.
func LookupSRV(serviceName string, namespace string) (string, []*net.SRV, error) {
	return DefaultResolver.LookupSRV(context.Background(), serviceName, namespace)
}

// LookupSRV constructs the DNS name with the format used in Kubernetes.
// The returned records are sorted by priority and randomized by weight within
// a priority.
//
// If no namespace is provided, the function will use the default kubernetes
// namespace.
func (r *Resolver) LookupSRV(ctx context.Context, serviceName string, namespace string) (string, []*net.SRV, error) {
	nm := defaultNamespace
	if len(namespace) > 0 {
		nm = namespace
	}

	name := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, nm)
	return defaultResolver.LookupSRV(ctx, "", "", name)
}
