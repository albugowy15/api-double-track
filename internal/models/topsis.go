package models

type TOPSIS struct {
	StudentId string `db:"student_id"`
	Id        int32  `db:"id"`
}

type TOPSISToAlternatives struct {
	Score         float32 `db:"score"`
	TopsisId      int32   `db:"topsis_id"`
	AlternativeId int     `db:"alternative_id"`
}
