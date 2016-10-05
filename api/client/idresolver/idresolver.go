package idresolver

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
)

// IDResolver provides ID to Name resolution.
type IDResolver struct {
	client    client.APIClient
	noResolve bool
	cache     map[string]string
}

// New creates a new IDResolver.
func New(client client.APIClient, noResolve bool) *IDResolver {
	return &IDResolver{
		client:    client,
		noResolve: noResolve,
		cache:     make(map[string]string),
	}
}

// Resolve will attempt to resolve an ID to a Name by querying the manager.
// Results are stored into a cache.
// If the `-n` flag is used in the command-line, resolution is disabled.
func (r *IDResolver) Resolve(ctx context.Context, t interface{}, id string) (string, error) {
	if r.noResolve {
		return id, nil
	}
	if name, ok := r.cache[id]; ok {
		return name, nil
	}

	r.cache[id] = name
	return name, nil
}
