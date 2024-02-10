package models

type Category struct {
	ID    int64  `json:"id" form:"id" db:"id" uri:"category"`
	Title string `db:"title" json:"title" form:"title"`
}
