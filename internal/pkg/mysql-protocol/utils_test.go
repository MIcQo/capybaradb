package mysql_protocol

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandStringRunes(t *testing.T) {
	var randomString = RandStringRunes(8)

	assert.NotEmpty(t, randomString)
	assert.Len(t, randomString, 8)
}
