package models

type Organization struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
