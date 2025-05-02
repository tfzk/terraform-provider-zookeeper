package client_test

import (
	"testing"

	"github.com/go-zookeeper/zk"
	testifyAssert "github.com/stretchr/testify/assert"
	testifyRequire "github.com/stretchr/testify/require"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func initTest(
	t *testing.T,
) (*client.Client, *testifyAssert.Assertions, *testifyRequire.Assertions) {
	t.Helper()
	assert := testifyAssert.New(t)
	require := testifyRequire.New(t)

	client, err := client.NewClientFromEnv()
	require.NoError(err)
	require.NoError(err)

	return client, assert, require
}

func TestClassicCRUD(t *testing.T) {
	client, assert, require := initTest(t)
	defer client.Close()

	// confirm not exists yet
	znodeExists, err := client.Exists("/test/ClassicCRUD")
	require.NoError(err)
	assert.False(znodeExists)

	// create
	znode, err := client.Create("/test/ClassicCRUD", []byte("one"), zk.WorldACL(zk.PermAll))
	require.NoError(err)
	assert.Equal("/test/ClassicCRUD", znode.Path)
	assert.Equal([]byte("one"), znode.Data)

	// confirm exists
	znodeExists, err = client.Exists("/test/ClassicCRUD")
	require.NoError(err)
	assert.True(znodeExists)

	// read
	znode, err = client.Read("/test/ClassicCRUD")
	require.NoError(err)
	assert.Equal("/test/ClassicCRUD", znode.Path)
	assert.Equal([]byte("one"), znode.Data)

	// update
	znode, err = client.Update("/test/ClassicCRUD", []byte("two"), zk.WorldACL(zk.PermAll))
	require.NoError(err)
	assert.Equal("/test/ClassicCRUD", znode.Path)
	assert.Equal([]byte("two"), znode.Data)

	// delete
	err = client.Delete("/test/ClassicCRUD")
	require.NoError(err)

	// confirm not exists
	znodeExists, err = client.Exists("/test/ClassicCRUD")
	require.NoError(err)
	assert.False(znodeExists)

	// confirm container still exists
	znodeExists, err = client.Exists("/test")
	require.NoError(err)
	assert.True(znodeExists)

	// delete container
	err = client.Delete("/test")
	require.NoError(err)
}

func TestCreateSequential(t *testing.T) {
	client, assert, require := initTest(t)
	defer client.Close()

	noPrefixSeqZNode, err := client.CreateSequential(
		"/test/CreateSequential/",
		[]byte("seq"),
		zk.WorldACL(zk.PermAll),
	)
	require.NoError(err)
	assert.Equal("/test/CreateSequential/0000000000", noPrefixSeqZNode.Path)

	prefixSeqZNode, err := client.CreateSequential(
		"/test/CreateSequentialWithPrefix/prefix-",
		[]byte("seq"),
		zk.WorldACL(zk.PermAll),
	)
	require.NoError(err)
	assert.Equal("/test/CreateSequentialWithPrefix/prefix-0000000000", prefixSeqZNode.Path)

	// delete, recursively
	err = client.Delete("/test")
	require.NoError(err)
}

func TestDigestAuthenticationSuccess(t *testing.T) {
	t.Setenv(client.EnvZooKeeperUsername, "username")
	t.Setenv(client.EnvZooKeeperPassword, "password")
	client, assert, require := initTest(t)
	defer client.Close()

	// Create a ZNode accessible only by the given user
	acl := zk.DigestACL(zk.PermAll, "username", "password")
	znode, err := client.Create("/auth-test/DigestAuthentication", []byte("data"), acl)
	require.NoError(err)
	assert.Equal("/auth-test/DigestAuthentication", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Make sure it's accessible
	znode, err = client.Read("/auth-test/DigestAuthentication")
	require.NoError(err)
	assert.Equal("/auth-test/DigestAuthentication", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Cleanup
	err = client.Delete("/auth-test/DigestAuthentication")
	require.NoError(err)
	err = client.Delete("/auth-test")
	require.NoError(err)
}

func TestFailureWhenReadingZNodeWithIncorrectAuth(t *testing.T) {
	// Create client authenticated as foo user
	t.Setenv(client.EnvZooKeeperUsername, "foo")
	t.Setenv(client.EnvZooKeeperPassword, "password")
	fooClient, assert, require := initTest(t)
	defer fooClient.Close()

	// Create a ZNode accessible only by foo user
	acl := zk.DigestACL(zk.PermAll, "foo", "password")
	znode, err := fooClient.Create("/auth-fail-test/AccessibleOnlyByFoo", []byte("data"), acl)
	require.NoError(err)
	assert.Equal("/auth-fail-test/AccessibleOnlyByFoo", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Make sure it's accessible by foo user
	znode, err = fooClient.Read("/auth-fail-test/AccessibleOnlyByFoo")
	require.NoError(err)
	assert.Equal("/auth-fail-test/AccessibleOnlyByFoo", znode.Path)
	assert.Equal([]byte("data"), znode.Data)
	assert.Equal(acl, znode.ACL)

	// Create client authenticated as bar user
	t.Setenv(client.EnvZooKeeperUsername, "bar")
	t.Setenv(client.EnvZooKeeperPassword, "password")
	barClient, err := client.NewClientFromEnv()
	require.NoError(err)
	defer barClient.Close()

	// The node should be inaccessible by bar user
	_, err = barClient.Read("/auth-fail-test/AccessibleOnlyByFoo")
	require.EqualError(
		err,
		"failed to read ZNode '/auth-fail-test/AccessibleOnlyByFoo': zk: not authenticated",
	)

	// Cleanup
	err = fooClient.Delete("/auth-fail-test/AccessibleOnlyByFoo")
	require.NoError(err)
	err = fooClient.Delete("/auth-fail-test")
	require.NoError(err)
}

func TestFailureWhenCreatingForNonSequentialZNodeEndingInSlash(t *testing.T) {
	client, assert, require := initTest(t)
	defer client.Close()

	_, err := client.Create("/test/willFail/", nil, zk.WorldACL(zk.PermAll))
	require.Error(err)
	assert.Equal(
		"non-sequential ZNode cannot have path '/test/willFail/' because it ends in '/'",
		err.Error(),
	)
}

func TestFailureWhenCreatingWhenZNodeAlreadyExists(t *testing.T) {
	client, assert, require := initTest(t)
	defer client.Close()

	_, err := client.Create("/test/node", nil, zk.WorldACL(zk.PermAll))
	require.NoError(err)
	_, err = client.Create("/test/node", nil, zk.WorldACL(zk.PermAll))
	require.Error(err)
	assert.Equal(
		"failed to create ZNode '/test/node' (size: 0, createFlags: 0, acl: [{31 world anyone}]): zk: node already exists",
		err.Error(),
	)

	err = client.Delete("/test")
	require.NoError(err)
}

func TestFailureWithNonExistingZNodes(t *testing.T) {
	client, assert, require := initTest(t)
	defer client.Close()

	_, err := client.Read("/does-not-exist")
	require.Error(err)
	assert.Equal("failed to read ZNode '/does-not-exist': zk: node does not exist", err.Error())

	_, err = client.Update("/also-does-not-exist", nil, zk.WorldACL(zk.PermAll))
	require.Error(err)
	assert.Equal("failed to update ZNode '/also-does-not-exist': does not exist", err.Error())
}
