package models

type AHP struct {
	StudentId        string  `db:"student_id"`
	Id               int32   `db:"id"`
	ConsistencyRatio float32 `db:"consistency_ratio"`
}

type AHPToAlternatives struct {
	Score         float32 `db:"score"`
	AhpId         int32   `db:"ahp_id"`
	AlternativeId int     `db:"alternative_id"`
}

type TOPSIS struct {
	StudentId string `db:"student_id"`
	Id        int32  `db:"id"`
}

type TOPSISToAlternatives struct {
	Score         float32 `db:"score"`
	TopsisId      int32   `db:"topsis_id"`
	AlternativeId int     `db:"alternative_id"`
}
