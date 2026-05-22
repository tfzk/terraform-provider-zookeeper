// Package client provides a wrapper around the go-zookeeper library.
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

// Client wraps a go-zookeeper `zk.Conn` object.
//
// It's designed to offer the functionalities that we will expose via the
// actual Terraform Provider.
type Client struct {
	zkConn *zk.Conn
}

// ZNode represents, obviously, a ZooKeeper Node.
//
// While `Path` and `Data` fields are pretty self-explanatory,
// the `Stat` contains multiple ZooKeeper related metadata.
// See `zk.Stat` for details.
type ZNode struct {
	Path string
	Stat *zk.Stat
	Data []byte
	ACL  []zk.ACL
}

// Re-exporting errors from the ZK library for better encapsulation.
var (
	ErrZNodeAlreadyExists = zk.ErrNodeExists
	ErrZNodeDoesNotExist  = zk.ErrNoNode
	ErrZNodeHasChildren   = zk.ErrNotEmpty
	ErrConnectionClosed   = zk.ErrConnectionClosed
	ErrInvalidArguments   = zk.ErrBadArguments
)

var (
	// ErrUserPassBothOrNone returned when only one of username and password is specified: either both or none is allowed.
	ErrUserPassBothOrNone = errors.New("both username and password must be specified together")

	// ErrTLSParseCACert returned when parsing the root CA certificate failed.
	ErrTLSParseCACert = errors.New("unable to parse TLS root CA cert")

	// ErrTLSCertKeyBothOrNone returned when one of either client certificate or client key are specified, but the other is not.
	ErrTLSCertKeyBothOrNone = errors.New("TLS cert and key file paths are mutually inclusive " +
		"(if one is specified, the other must be too)")
)

const (
	serversStringSeparator = ","
	zNodeRootPath          = "/"
	zNodePathSeparator     = '/'

	// matchAnyVersion is used when submitting an update/delete request.
	// Providing `version = -1` means that the operation will match any
	// version of the ZNode found.
	matchAnyVersion = -1

	// EnvZooKeeperServer environment variable containing a comma separated
	// list of 'host:port' pairs, pointing at ZooKeeper Server(s).
	// This is used by NewClientFromEnv.
	EnvZooKeeperServer = "ZOOKEEPER_SERVERS"

	// EnvZooKeeperSessionSec environment variable defining how many seconds
	// a session is considered valid after losing connectivity.
	// This is used by NewClientFromEnv.
	EnvZooKeeperSessionSec = "ZOOKEEPER_SESSION"

	// DefaultZooKeeperSessionSec is the default amount of seconds configured for the
	// Client timeout session, in case EnvZooKeeperSessionSec is not set.
	DefaultZooKeeperSessionSec = 30

	// EnvZooKeeperUsername environment variable providing the username part of a digest auth credentials.
	// This is used by NewClientFromEnv.
	EnvZooKeeperUsername = "ZOOKEEPER_USERNAME"

	// EnvZooKeeperPassword environment variable providing the password part of a digest auth credentials.
	// This is used by NewClientFromEnv.
	EnvZooKeeperPassword = "ZOOKEEPER_PASSWORD"

	// EnvZooKeeperTLSEnable environment variable enabling a TLS connection to the server(s).
	// This is used by NewClientFromEnv.
	EnvZooKeeperTLSEnable = "ZOOKEEPER_TLS_ENABLE"

	// EnvZooKeeperTLSSkipVerify environment variable disabling verification of server's certificate chain and host name.
	// This is used by NewClientFromEnv.
	EnvZooKeeperTLSSkipVerify = "ZOOKEEPER_TLS_SKIP_VERIFY"

	// EnvZooKeeperTLSRootCertPath environment variable providing file path to the TLS root certificate.
	// This is used by NewClientFromEnv.
	EnvZooKeeperTLSRootCertPath = "ZOOKEEPER_TLS_ROOT_CA_CERT_PATH"

	// EnvZooKeeperTLSCertPath environment variable providing file path to the TLS certificate.
	// This is used by NewClientFromEnv.
	EnvZooKeeperTLSCertPath = "ZOOKEEPER_TLS_CERT_PATH"

	// EnvZooKeeperTLSKeyPath environment variable providing file path to the TLS key.
	// This is used by NewClientFromEnv.
	EnvZooKeeperTLSKeyPath = "ZOOKEEPER_TLS_KEY_PATH"
)

// NewClient constructs a new Client instance.
func NewClient(
	servers string,
	sessionTimeoutSec int,
	username string,
	password string,
	tlsConfig *TLSConfig,
) (*Client, error) {
	serversSplit := strings.Split(servers, serversStringSeparator)

	dialer, err := newDialer(tlsConfig)
	if err != nil {
		return nil, err
	}

	conn, _, err := zk.Connect(
		zk.FormatServers(serversSplit),
		time.Duration(sessionTimeoutSec)*time.Second,
		zk.WithDialer(dialer),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to ZooKeeper: %w", err)
	}
	fmt.Printf("[DEBUG] Connected to ZooKeeper servers %s\n", serversSplit)

	if (username == "") != (password == "") {
		return nil, ErrUserPassBothOrNone
	}

	if username != "" {
		auth := "digest"
		credentials := fmt.Sprintf("%s:%s", username, password)
		err = conn.AddAuth(auth, []byte(credentials))
		if err != nil {
			return nil, fmt.Errorf("unable to add digest auth: %w", err)
		}
	}

	return &Client{
		zkConn: conn,
	}, nil
}

func newDialer(tlsConfig *TLSConfig) (zk.Dialer, error) { // #nosec G402
	tlsDialerConfig := &tls.Config{
		InsecureSkipVerify: tlsConfig.SkipVerify,
	}

	if tlsConfig.RootCertPath != "" {
		tlsDialerConfig.RootCAs = x509.NewCertPool()

		tlsRootCert, err := os.ReadFile(tlsConfig.RootCertPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read TLS root CA cert file: %w", err)
		}

		if !tlsDialerConfig.RootCAs.AppendCertsFromPEM(tlsRootCert) {
			return nil, ErrTLSParseCACert
		}
	}

	if tlsConfig.CertPath != "" || tlsConfig.KeyPath != "" {
		if tlsConfig.CertPath == "" || tlsConfig.KeyPath == "" {
			return nil, ErrTLSCertKeyBothOrNone
		}

		tlsCert, err := os.ReadFile(tlsConfig.CertPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read TLS client cert file: %w", err)
		}

		tlsKey, err := os.ReadFile(tlsConfig.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read TLS client key file: %w", err)
		}

		certificate, err := tls.X509KeyPair(tlsCert, tlsKey)
		if err != nil {
			return nil, fmt.Errorf("unable to parse TLS client X509 key pair: %w", err)
		}

		tlsDialerConfig.Certificates = []tls.Certificate{certificate}
	}

	return func(network, address string, timeout time.Duration) (net.Conn, error) {
		ctx := context.Background()

		dialer := &net.Dialer{Timeout: timeout}

		if tlsConfig.Enable {
			tlsDialer := &tls.Dialer{
				NetDialer: dialer,
				Config:    tlsDialerConfig,
			}

			return tlsDialer.DialContext(ctx, network, address)
		}

		return dialer.DialContext(ctx, network, address)
	}, nil
}

// NewClientFromEnv constructs a Client instance from environment variables.
//
// The only mandatory environment variable is EnvZooKeeperServer.
func NewClientFromEnv() (*Client, error) {
	zkServers, ok := os.LookupEnv(EnvZooKeeperServer)
	if !ok {
		return nil, NewMissingEnvVarError(EnvZooKeeperServer)
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

	tlsConfig := &TLSConfig{}

	tlsEnable, _ := os.LookupEnv(EnvZooKeeperTLSEnable)
	tlsConfig.Enable = tlsEnable == "true"

	tlsSkipVerify, _ := os.LookupEnv(EnvZooKeeperTLSSkipVerify)
	tlsConfig.SkipVerify = tlsSkipVerify == "true"

	tlsConfig.RootCertPath, _ = os.LookupEnv(EnvZooKeeperTLSRootCertPath)
	tlsConfig.CertPath, _ = os.LookupEnv(EnvZooKeeperTLSCertPath)
	tlsConfig.KeyPath, _ = os.LookupEnv(EnvZooKeeperTLSKeyPath)

	fmt.Println("[DEBUG] Creating Client from Environment Variables")
	return NewClient(
		zkServers,
		zkSessionInt,
		zkUsername,
		zkPassword,
		tlsConfig,
	)
}

// Create a ZNode at the given path.
//
// Note that any necessary ZNode parents will be created if absent.
func (c *Client) Create(path string, data []byte, acl []zk.ACL) (*ZNode, error) {
	if path[len(path)-1] == zNodePathSeparator {
		return nil, NewNonSeqZNodeCannotEndWithPathSeparatorError(path)
	}

	return c.doCreate(path, data, 0, acl)
}

// CreateSequential will create a ZNode at the given path, using the Sequential Node flag.
//
// See: https://zookeeper.apache.org/doc/r3.6.3/zookeeperProgrammers.html#Sequence+Nodes+--+Unique+Naming
//
// This will ensure unique naming within the same parent ZNode,
// by appending a monotonically increasing counter in the format `%010d`
// (that is 10 digits with 0 (zero) padding).
// Note that if the `path` ends in `/`, the ZNode name will be just the counter
// described above. For example:
//
//   - input path         -> `/this/is/a/path/`
//   - created znode path -> `/this/is/a/path/0000000001`
//
// Note also that any necessary ZNode parents will be created if absent.
func (c *Client) CreateSequential(path string, data []byte, acl []zk.ACL) (*ZNode, error) {
	return c.doCreate(path, data, zk.FlagSequence, acl)
}

func (c *Client) doCreate(
	path string,
	data []byte,
	createFlags int32,
	acl []zk.ACL,
) (*ZNode, error) {
	// Create any necessary parent for the ZNode we need to crete
	parentZNodes := listParentsInOrder(path)
	err := c.createEmptyZNodes(parentZNodes, 0, acl)
	if err != nil {
		return nil, err
	}

	// NOTE: Based on the `createFlags`, the path returned by `Create` can change (ex. sequential nodes)
	createdPath, err := c.zkConn.Create(path, data, createFlags, acl)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create ZNode '%s' (size: %d, createFlags: %d, acl: %v): %w",
			path,
			len(data),
			createFlags,
			acl,
			err,
		)
	}

	return c.Read(createdPath)
}

func listParentsInOrder(path string) []string {
	// Split the path one parent directory at a time
	parentPaths := []string{filepath.Dir(path)}
	for parentPaths[len(parentPaths)-1] != zNodeRootPath {
		parentPaths = append(parentPaths, filepath.Dir(parentPaths[len(parentPaths)-1]))
	}

	// Sort by increasing length (i.e. each parent before each child)
	sort.Strings(parentPaths)

	// Return all the parents, excluding `root`
	return parentPaths[1:]
}

func (c *Client) createEmptyZNodes(pathsInOrder []string, createFlags int32, acl []zk.ACL) error {
	for _, path := range pathsInOrder {
		exists, err := c.Exists(path)
		if err != nil {
			return err
		}

		// Will only create the znode if they don't already exist.
		//
		// NOTE: Terraform graph can sometimes decide to create multiple
		// ZNodes that share part of their path ancestry at the same time.
		// When that happens, we have contention in this area of code,
		// where a `path` that didn't exist above, it exists once we try
		// to create it.
		// For this reason, we avoid reporting an error if it is about
		// a ZNode already existing.
		if !exists {
			_, err := c.zkConn.Create(path, nil, createFlags, acl)
			if err != nil && !errors.Is(err, ErrZNodeAlreadyExists) {
				return fmt.Errorf(
					"failed to create parent ZNode '%s' (createFlags: %d, acl: %v): %w",
					path,
					createFlags,
					acl,
					err,
				)
			}
		}
	}

	return nil
}

// Read the ZNode at the given path.
func (c *Client) Read(path string) (*ZNode, error) {
	data, stat, err := c.zkConn.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read ZNode '%s': %w", path, err)
	}

	acls, _, err := c.zkConn.GetACL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ACLs for ZNode '%s': %w", path, err)
	}

	return &ZNode{
		Path: path,
		Stat: stat,
		Data: data,
		ACL:  acls,
	}, nil
}

// Update the ZNode at the given path, under the assumption that it is there.
//
// Will return an error if it doesn't already exist.
func (c *Client) Update(path string, data []byte, acl []zk.ACL) (*ZNode, error) {
	exists, err := c.Exists(path)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, NewCannotUpdateDoesNotExistError(path)
	}

	_, err = c.zkConn.SetACL(path, acl, matchAnyVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to update ZNode '%s' ACL: %w", path, err)
	}

	_, err = c.zkConn.Set(path, data, matchAnyVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to update ZNode '%s': %w", path, err)
	}

	return c.Read(path)
}

// Close the Client underlying connection.
func (c *Client) Close() {
	fmt.Println("[DEBUG] Closing underlying ZooKeeper connection")
	c.zkConn.Close()
}

// Delete the given ZNode.
//
// Note that will also delete any child ZNode, recursively.
func (c *Client) Delete(path string) error {
	children, _, err := c.zkConn.Children(path)
	if err != nil {
		return fmt.Errorf("failed to list children for ZNode '%s': %w", path, err)
	}

	for _, child := range children {
		childPath := fmt.Sprintf("%s%c%s", path, zNodePathSeparator, child)
		err = c.Delete(childPath)
		if err != nil {
			return fmt.Errorf("failed to delete child '%s' of ZNode '%s': %w", childPath, path, err)
		}
	}

	err = c.zkConn.Delete(path, matchAnyVersion)
	if err != nil {
		return fmt.Errorf("failed to delete ZNode '%s': %w", path, err)
	}
	return nil
}

// Exists checks for the existence of the given ZNode.
func (c *Client) Exists(path string) (bool, error) {
	exists, _, err := c.zkConn.Exists(path)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of ZNode '%s': %w", path, err)
	}

	return exists, nil
}

// RemoveSequentialSuffix takes the path to a sequential ZNode, maybe created via CreateSequential,
// and truncates the unique suffix.
//
// See: https://zookeeper.apache.org/doc/r3.6.3/zookeeperProgrammers.html#Sequence+Nodes+--+Unique+Naming
func RemoveSequentialSuffix(path string) string {
	return path[:len(path)-10]
}
