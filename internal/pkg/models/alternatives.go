package models

type Alternative struct {
	Alternative string `db:"alternative" json:"alternative"`
	Description string `db:"description" json:"description"`
	Id          int    `db:"id" json:"id"`
}
