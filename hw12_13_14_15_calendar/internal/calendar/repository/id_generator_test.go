package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateEventId(t *testing.T) {
	id1 := GenerateEventID()
	id2 := GenerateEventID()
	assert.Greater(t, len(id1.String()), 0)
	assert.Greater(t, len(id2.String()), 0)
	assert.NotEqual(t, id1, id2)
}
