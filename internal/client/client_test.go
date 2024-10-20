package client_test

import (
	"testing"

	"github.com/go-zookeeper/zk"
	testifyAssert "github.com/stretchr/testify/assert"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func initTest(t *testing.T) (*client.Client, *testifyAssert.Assertions) {
	assert := testifyAssert.New(t)

	client, err := client.DefaultPool().GetClientFromEnv()
	assert.NoError(err)

	return client, assert
}

func TestClassicCRUD(t *testing.T) {
	client, assert := initTest(t)

	// confirm not exists yet
	znodeExists, err := client.Exists("/test/ClassicCRUD")
	assert.NoError(err)
	assert.False(znodeExists)

	// create
	znode, err := client.Create("/test/ClassicCRUD", []byte("one"), zk.WorldACL(zk.PermAll))
	assert.NoError(err)
	assert.Equal("/test/ClassicCRUD", znode.Path)
	assert.Equal([]byte("one"), znode.Data)

	// confirm exists
	znodeExists, err = client.Exists("/test/ClassicCRUD")
	assert.NoError(err)
	assert.True(znodeExists)

	// read
	znode, err = client.Read("/test/ClassicCRUD")
	assert.NoError(err)
	assert.Equal("/test/ClassicCRUD", znode.Path)
	assert.Equal([]byte("one"), znode.Data)

	// update
	znode, err = client.Update("/test/ClassicCRUD", []byte("two"), zk.WorldACL(zk.PermAll))
	assert.NoError(err)
	assert.Equal("/test/ClassicCRUD", znode.Path)
	assert.Equal([]byte("two"), znode.Data)

	// delete
	err = client.Delete("/test/ClassicCRUD")
	assert.NoError(err)

	// confirm not exists
	znodeExists, err = client.Exists("/test/ClassicCRUD")
	assert.NoError(err)
	assert.False(znodeExists)

	// confirm container still exists
	znodeExists, err = client.Exists("/test")
	assert.NoError(err)
	assert.True(znodeExists)

	// delete container
	err = client.Delete("/test")
	assert.NoError(err)
}

func TestCreateSequential(t *testing.T) {
	client, assert := initTest(t)

	noPrefixSeqZNode, err := client.CreateSequential("/test/CreateSequential/", []byte("seq"), zk.WorldACL(zk.PermAll))
	assert.NoError(err)
	assert.Equal("/test/CreateSequential/0000000000", noPrefixSeqZNode.Path)

	prefixSeqZNode, err := client.CreateSequential("/test/CreateSequentialWithPrefix/prefix-", []byte("seq"), zk.WorldACL(zk.PermAll))
	assert.NoError(err)
	assert.Equal("/test/CreateSequentialWithPrefix/prefix-0000000000", prefixSeqZNode.Path)

	// delete, recursively
	err = client.Delete("/test")
	assert.NoError(err)
}

func TestDigestAuthenticationSuccess(t *testing.T) {
	t.Setenv(client.EnvZooKeeperUsername, "username")
	t.Setenv(client.EnvZooKeeperPassword, "password")
	client, assert := initTest(t)

	// Create a ZNode accessible only by the given user
	acl := zk.DigestACL(zk.PermAll, "username", "password")
	znode, err := client.Create("/auth-test/DigestAuthentication", []byte("data"), acl)
	assert.NoError(err)
	assert.Equal("/auth-test/DigestAuthentication", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Make sure it's accessible
	znode, err = client.Read("/auth-test/DigestAuthentication")
	assert.NoError(err)
	assert.Equal("/auth-test/DigestAuthentication", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Cleanup
	err = client.Delete("/auth-test/DigestAuthentication")
	assert.NoError(err)
	err = client.Delete("/auth-test")
	assert.NoError(err)
}

func TestFailureWhenReadingZNodeWithIncorrectAuth(t *testing.T) {
	// Create client authenticated as foo user
	t.Setenv(client.EnvZooKeeperUsername, "foo")
	t.Setenv(client.EnvZooKeeperPassword, "password")
	fooClient, assert := initTest(t)

	// Create a ZNode accessible only by foo user
	acl := zk.DigestACL(zk.PermAll, "foo", "password")
	znode, err := fooClient.Create("/auth-fail-test/AccessibleOnlyByFoo", []byte("data"), acl)
	assert.NoError(err)
	assert.Equal("/auth-fail-test/AccessibleOnlyByFoo", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Make sure it's accessible by foo user
	znode, err = fooClient.Read("/auth-fail-test/AccessibleOnlyByFoo")
	assert.NoError(err)
	assert.Equal("/auth-fail-test/AccessibleOnlyByFoo", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Create client authenticated as bar user
	t.Setenv(client.EnvZooKeeperUsername, "bar")
	t.Setenv(client.EnvZooKeeperPassword, "password")
	barClient, err := client.DefaultPool().GetClientFromEnv()
	assert.NoError(err)

	// The node should be inaccessible by bar user
	_, err = barClient.Read("/auth-fail-test/AccessibleOnlyByFoo")
	assert.EqualError(err, "failed to read ZNode '/auth-fail-test/AccessibleOnlyByFoo': zk: not authenticated")

	// Cleanup
	err = fooClient.Delete("/auth-fail-test/AccessibleOnlyByFoo")
	assert.NoError(err)
	err = fooClient.Delete("/auth-fail-test")
	assert.NoError(err)
}

func TestFailureWhenCreatingForNonSequentialZNodeEndingInSlash(t *testing.T) {
	client, assert := initTest(t)

	_, err := client.Create("/test/willFail/", nil, zk.WorldACL(zk.PermAll))
	assert.Error(err)
	assert.Equal("non-sequential ZNode cannot have path '/test/willFail/' because it ends in '/'", err.Error())
}

func TestFailureWhenCreatingWhenZNodeAlreadyExists(t *testing.T) {
	client, assert := initTest(t)

	_, err := client.Create("/test/node", nil, zk.WorldACL(zk.PermAll))
	assert.NoError(err)
	_, err = client.Create("/test/node", nil, zk.WorldACL(zk.PermAll))
	assert.Error(err)
	assert.Equal("failed to create ZNode '/test/node' (size: 0, createFlags: 0, acl: [{31 world anyone}]): zk: node already exists", err.Error())

	err = client.Delete("/test")
	assert.NoError(err)
}

func TestFailureWithNonExistingZNodes(t *testing.T) {
	client, assert := initTest(t)

	_, err := client.Read("/does-not-exist")
	assert.Error(err)
	assert.Equal("failed to read ZNode '/does-not-exist': zk: node does not exist", err.Error())

	_, err = client.Update("/also-does-not-exist", nil, zk.WorldACL(zk.PermAll))
	assert.Error(err)
	assert.Equal("failed to update ZNode '/also-does-not-exist': does not exist", err.Error())
}
