package client_test

import (
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func initTest(t *testing.T) (*client.Client, *testifyAssert.Assertions) {
	assert := testifyAssert.New(t)

	client, err := client.NewClientFromEnv()
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
	znode, err := client.Create("/test/ClassicCRUD", []byte("one"))
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
	znode, err = client.Update("/test/ClassicCRUD", []byte("two"))
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

	noPrefixSeqZNode, err := client.CreateSequential("/test/CreateSequential/", []byte("seq"))
	assert.NoError(err)
	assert.Equal("/test/CreateSequential/0000000000", noPrefixSeqZNode.Path)

	prefixSeqZNode, err := client.CreateSequential("/test/CreateSequentialWithPrefix/prefix-", []byte("seq"))
	assert.NoError(err)
	assert.Equal("/test/CreateSequentialWithPrefix/prefix-0000000000", prefixSeqZNode.Path)

	// delete, recursively
	err = client.Delete("/test")
	assert.NoError(err)
}

func TestFailureWhenCreatingForNonSequentialZNodeEndingInSlash(t *testing.T) {
	client, assert := initTest(t)

	_, err := client.Create("/test/willFail/", nil)
	assert.Error(err)
	assert.Equal("non-sequential ZNode cannot have path '/test/willFail/' because it ends in '/'", err.Error())
}

func TestFailureWhenCreatingWhenZNodeAlreadyExists(t *testing.T) {
	client, assert := initTest(t)

	_, err := client.Create("/test/node", nil)
	assert.NoError(err)
	_, err = client.Create("/test/node", nil)
	assert.Error(err)
	assert.Equal("failed to create ZNode '/test/node' (size: 0, createFlags: 0, acl: [{15 world anyone}]): zk: node already exists", err.Error())

	err = client.Delete("/test")
	assert.NoError(err)
}

func TestFailureWithNonExistingZNodes(t *testing.T) {
	client, assert := initTest(t)

	_, err := client.Read("/does-not-exist")
	assert.Error(err)
	assert.Equal("failed to read ZNode '/does-not-exist': zk: node does not exist", err.Error())

	_, err = client.Update("/also-does-not-exist", nil)
	assert.Error(err)
	assert.Equal("failed to update ZNode '/also-does-not-exist': does not exist", err.Error())
}
