package models

import "github.com/guregu/null/v5"

type Alternative struct {
	Alternative string      `db:"alternative" json:"alternative"`
	Description null.String `db:"description" json:"description"`
	Id          int         `db:"id" json:"id"`
}
