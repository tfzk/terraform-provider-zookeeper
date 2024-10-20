package client

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Pool wraps a collection of Client.
// Each client is associated to a unique set of construction parameters.
type Pool struct {
	mu   sync.Mutex
	pool map[string]*Client
}

func newPool() *Pool {
	return &Pool{
		pool: make(map[string]*Client),
	}
}

// GetClient retrieves (or creates) a Client.
func (p *Pool) GetClient(servers string, sessionTimeoutSec int, username string, password string) (*Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	clientKey := fmt.Sprintf("%s::%d::%s::%s", servers, sessionTimeoutSec, username, password)
	if client, found := p.pool[clientKey]; found {
		return client, nil
	}

	client, err := NewClient(servers, sessionTimeoutSec, username, password)
	p.pool[clientKey] = client

	return client, err
}

// GetClientFromEnv retrieves (or creates) a Client instance from environment variables.
//
// The only mandatory environment variable is EnvZooKeeperServer.
func (p *Pool) GetClientFromEnv() (*Client, error) {
	zkServers, ok := os.LookupEnv(EnvZooKeeperServer)
	if !ok {
		return nil, fmt.Errorf("missing environment variable: %s", EnvZooKeeperServer)
	}

	zkSession, ok := os.LookupEnv(EnvZooKeeperSessionSec)
	if !ok {
		zkSession = strconv.FormatInt(DefaultZooKeeperSessionSec, 10)
	}

	zkSessionInt, err := strconv.Atoi(zkSession)
	if err != nil {
		return nil, fmt.Errorf("failed to convert '%s' to integer: %w", zkSession, err)
	}

	zkUsername, _ := os.LookupEnv(EnvZooKeeperUsername)
	zkPassword, _ := os.LookupEnv(EnvZooKeeperPassword)

	return p.GetClient(zkServers, zkSessionInt, zkUsername, zkPassword)
}

//nolint:gochecknoglobals
var defaultPool *Pool

func init() {
	defaultPool = newPool()
}

// DefaultPool returns a pointer to the default (and unique) instance of Pool.
// This is essentially a singleton.
func DefaultPool() *Pool {
	return defaultPool
}
