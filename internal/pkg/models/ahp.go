package models

type AHP struct {
	StudentId        string  `db:"student_id"`
	Id               int32   `db:"id"`
	ConsistencyIndex float32 `db:"consistency_index"`
}

type AHPToAlternatives struct {
	Score         float32 `db:"score"`
	AhpId         int32   `db:"ahp_id"`
	AlternativeId int     `db:"alternative_id"`
}
