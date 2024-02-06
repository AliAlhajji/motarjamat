package models

type Vocab struct {
	ID      string `db:"id" json:"id" uri:"id" binding:"required,uuid"`
	English string `db:"english" json:"english" uri:"english"`
	Arabic  string `db:"arabic" json:"arabic" uri:"arabic"`
	Meaning string `db:"meaning" json:"meaning"`
}
