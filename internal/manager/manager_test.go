package manager

import (
	"testing"
	connection_preparer "webserver/internal/connection-preparer"

	"github.com/stretchr/testify/assert"
)

func TestCreateConnectionManagerShouldCreateManagerWithGivenPort(t *testing.T) {
	cm := NewConcurrentConnectionManger(nil, &connection_preparer.PlainConnectionPreparer{}, 8081)

	assert.Equal(t, uint32(8081), cm.port)
}
