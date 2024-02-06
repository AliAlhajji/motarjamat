package models

import (
	"time"
)

const RoleAdmin string = "admin"
const RoleUser string = "user"

type User struct {
	ID       int64     `db:"id" json:"id" uri:"id"`
	UUID     string    `db:"uuid" json:"uuid" uri:"uuid"`
	Email    string    `db:"email" form:"email" json:"email"`
	Password string    `db:"password" form:"password"`
	Username string    `db:"username" form:"username" json:"username" uri:"username"`
	Name     string    `db:"name" form:"name" json:"name"`
	Role     string    `db:"role" json:"role"`
	JoinDate time.Time `db:"join_date" json:"join_date"`
}
