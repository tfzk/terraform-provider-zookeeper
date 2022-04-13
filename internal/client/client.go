package client

import (
	"fmt"
	"path/filepath"
	"sort"
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
}

const (
	ServersStringSeparator = ","
	ZNodeRootPath          = "/"
	ZNodePathSeparator     = '/'

	// MatchAnyVersion is used when submitting an update/delete request.
	// Providing `version = -1` means that the operation will match any
	// version of the ZNode found.
	MatchAnyVersion = -1
)

// NewClient constructs a new `Client` instance.
func NewClient(servers string, sessionTimeoutSec int) (*Client, error) {
	serversSplit := strings.Split(servers, ServersStringSeparator)

	conn, _, err := zk.Connect(zk.FormatServers(serversSplit), time.Duration(sessionTimeoutSec)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to ZooKeeper: %w", err)
	}

	return &Client{
		zkConn: conn,
	}, nil
}

// Create a ZNode at the given path.
//
// Note that any necessary ZNode parents will be created if absent.
func (c *Client) Create(path string, data []byte) (ZNode, error) {
	// TODO Make ACL configurable
	acl := zk.WorldACL(zk.PermRead | zk.PermWrite | zk.PermCreate | zk.PermDelete)

	if path[len(path)-1] == ZNodePathSeparator {
		return ZNode{}, fmt.Errorf("non-sequential ZNode cannot have path '%s' because it ends in '%c'", path, ZNodePathSeparator)
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
//   input path         -> `/this/is/a/path/`
//   created znode path -> `/this/is/a/path/0000000001`
//
// Note also that any necessary ZNode parents will be created if absent.
func (c *Client) CreateSequential(path string, data []byte) (ZNode, error) {
	// TODO Make ACL configurable
	acl := zk.WorldACL(zk.PermRead | zk.PermWrite | zk.PermCreate | zk.PermDelete)

	return c.doCreate(path, data, zk.FlagSequence, acl)
}

func (c *Client) doCreate(path string, data []byte, createFlags int32, acl []zk.ACL) (ZNode, error) {
	// Create any necessary parent for the ZNode we need to crete
	parentZNodes := listParentsInOrder(path)
	err := c.createEmptyZNodes(parentZNodes, 0, acl)
	if err != nil {
		return ZNode{}, err
	}

	// NOTE: Based on the `createFlags`, the path returned by `Create` can change (ex. sequential nodes)
	createdPath, err := c.zkConn.Create(path, data, createFlags, acl)
	if err != nil {
		return ZNode{}, fmt.Errorf("failed to create ZNode '%s' (size: %d, createFlags: %d, acl: %v): %w", path, len(data), createFlags, acl, err)
	}

	return c.Read(createdPath)
}

func listParentsInOrder(path string) []string {
	// Split the path one parent directory at a time
	parentPaths := []string{filepath.Dir(path)}
	for parentPaths[len(parentPaths)-1] != ZNodeRootPath {
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

		// Will only create the znode if they don't already exist
		if !exists {
			_, err := c.zkConn.Create(path, nil, createFlags, acl)
			if err != nil {
				return fmt.Errorf("failed to create parent ZNode '%s' (createFlags: %d, acl: %v): %w", path, createFlags, acl, err)
			}
		}
	}

	return nil
}

// Read the ZNode at the given path.
func (c *Client) Read(path string) (ZNode, error) {
	data, stat, err := c.zkConn.Get(path)
	if err != nil {
		return ZNode{}, fmt.Errorf("failed to read ZNode '%s': %w", path, err)
	}

	return ZNode{
		Path: path,
		Stat: stat,
		Data: data,
	}, nil
}

// Update the ZNode at the given path, under the assumption that it is there.
//
// Will return an error if it doesn't already exist.
func (c *Client) Update(path string, data []byte) (ZNode, error) {
	exists, err := c.Exists(path)
	if err != nil {
		return ZNode{}, err
	}

	if !exists {
		return ZNode{}, fmt.Errorf("failed to update ZNode '%s': does not exist", path)
	}

	_, err = c.zkConn.Set(path, data, MatchAnyVersion)
	if err != nil {
		return ZNode{}, fmt.Errorf("failed to update ZNode '%s': %w", path, err)
	}

	return c.Read(path)
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
		childPath := fmt.Sprintf("%s%c%s", path, ZNodePathSeparator, child)
		err = c.Delete(childPath)
		if err != nil {
			return fmt.Errorf("failed to delete child '%s' of ZNode '%s': %w", childPath, path, err)
		}
	}

	err = c.zkConn.Delete(path, MatchAnyVersion)
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

// StatAsMap is an helper that returns the zk.Stat contained to ZNode.
func (z *ZNode) StatAsMap() map[string]int64 {
	return map[string]int64{
		"czxid":          z.Stat.Czxid,
		"mzxid":          z.Stat.Mzxid,
		"ctime":          z.Stat.Ctime,
		"mtime":          z.Stat.Mtime,
		"version":        int64(z.Stat.Version),
		"cversion":       int64(z.Stat.Cversion),
		"aversion":       int64(z.Stat.Aversion),
		"ephemeralOwner": z.Stat.EphemeralOwner,
		"dataLength":     int64(z.Stat.DataLength),
		"numChildren":    int64(z.Stat.NumChildren),
		"pzxid":          z.Stat.Pzxid,
	}
}
