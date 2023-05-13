package util

import (
	"errors"
	"sync/atomic"
)

var ErrGenID = errors.New("id overflow or negative")
var ErrClosed = errors.New("repository closed")
var ErrNotFound = errors.New("ad not found")

type UID struct {
	Id int64
}

func (u *UID) GenerateID() (int64, error) {
	newID := atomic.AddInt64(&u.Id, 1)
	if newID < 0 {
		return -1, ErrGenID
	}
	return newID, nil
}
