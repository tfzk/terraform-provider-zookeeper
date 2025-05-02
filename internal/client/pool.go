package client

import (
	"fmt"
	"sync"
)

// Pool contains a pool of Client.
// Each client is associated to a unique set of construction parameters.
type Pool struct {
	mu   sync.Mutex
	pool map[string]*Client
}

// NewPool creates a new Pool.
func NewPool() *Pool {
	return &Pool{
		pool: make(map[string]*Client),
	}
}

// GetOrCreateClient retrieves (or creates) a Client.
// A new client is created for each unique set of construction parameters.
func (p *Pool) GetOrCreateClient(
	servers string,
	sessionTimeoutSec int,
	username string,
	password string,
) (*Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	clientKey := fmt.Sprintf("%s::%d::%s::%s", servers, sessionTimeoutSec, username, password)

	// Return client if already present for the same key
	if client, found := p.pool[clientKey]; found {
		return client, nil
	}

	// Create new client, and cache it for the given key
	client, err := NewClient(servers, sessionTimeoutSec, username, password)
	p.pool[clientKey] = client

	return client, err
}
