package entities

import "time"

type Ad struct {
	ID         int64
	Title      string
	Text       string
	AuthorID   int64
	Published  bool
	CreateDate time.Time
	UpdateDate time.Time
}
