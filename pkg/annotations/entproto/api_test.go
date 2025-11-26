package entproto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPIMethod(t *testing.T) {
	assert.Equal(t, (APIGet | APIDelete).Methods(), []APIMethod{APIGet, APIDelete})
}
