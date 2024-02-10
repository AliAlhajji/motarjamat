package models

import "time"

type Post struct {
	ID         int64     `db:"id" json:"id" uri:"post"`
	Title      string    `db:"title" json:"title" form:"title"`
	Body       string    `db:"body" json:"body" form:"body"`
	UserID     string    `json:"userID" db:"user_id"`
	Username   string    `json:"username" db:"username"`
	Name       string    `json:"name" db:"name"`
	Link       string    `db:"link" json:"link" form:"link"`
	Date       time.Time `db:"date" json:"date"`
	Vocabs     []*Vocab  `json:"vocabs"`
	Categories []int64   `form:"categories[]"`
}
