package models

type Expectation struct {
	Id           string                     `db:"id"`
	StudentId    string                     `db:"student_id"`
	Expectations []ExpectationToAlternative `db:"expectations"`
}

type ExpectationToAlternative struct {
	ExpectationId   string `db:"expectation_id" json:"expectation_id"`
	AlternativeId   string `db:"alternative_id" json:"alternative_id"`
	AlternativeName string `db:"alternative_name" json:"alternative_name"`
	Id              int64  `db:"id" json:"id"`
	Rank            int    `db:"rank" json:"rank"`
}

type ExpectationData struct {
	AlternativeId int `json:"alternative_id"`
	Rank          int `json:"rank"`
}

type ExpectationRequest struct {
	Expectations []ExpectationData `json:"expectations"`
}

type ExpectationStatusResponse struct {
	Status string `json:"status"`
}
