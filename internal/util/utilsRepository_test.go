package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUID_GenerateID(t *testing.T) {
	uid := UID{Id: -1}
	id, err := uid.GenerateID()
	assert.Equal(t, int64(0), id)
	assert.NoError(t, err)
}
func TestUID_GenerateID_WrongUID(t *testing.T) {
	uid := UID{Id: -2}
	_, err := uid.GenerateID()
	assert.ErrorIs(t, ErrGenID, err)
}
