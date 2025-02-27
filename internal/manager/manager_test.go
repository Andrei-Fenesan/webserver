package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConnectionManagerShouldCreateManagerWithGivenPort(t *testing.T) {
	assert := assert.New(t)

	cm := NewConcurrentConnectionManger(nil, 8081)

	assert.Equal(uint32(8081), cm.port)
}
