package models

import (
	"github.com/guregu/null/v5"
)

type Question struct {
	Code        string      `db:"code" json:"code"`
	Question    string      `db:"question" json:"question"`
	Category    string      `db:"category" json:"category"`
	Description null.String `db:"description" json:"description"`
	Id          int         `db:"id" json:"id"`
	Number      int         `db:"number" json:"number"`
}

type QuestionnareSetting struct {
	SchoolId                    string     `db:"school_id" json:"school_id"`
	SchoolName                  string     `db:"school_name" json:"school_name"`
	AlternativeName             string     `db:"alternative_name" json:"alternative_name"`
	Id                          int        `db:"id" json:"id"`
	AlternativeId               int        `db:"alternative_id" json:"alternative_id"`
	TotalOpenJobs               null.Int16 `db:"total_open_jobs" json:"total_open_jobs"`
	EntrepreneurshipOpportunity null.Int16 `db:"enterpreneurship_opportunity" json:"enterpreneurship_opportunity"`
	Salary                      null.Int16 `db:"salary" json:"salary"`
}

type Answer struct {
	StudentId  string      `db:"student_id" json:"student_id"`
	AnswerStr  null.String `db:"answer_str" json:"answer_str"`
	Id         int64       `db:"id" json:"id"`
	QuestionId int         `db:"question_id" json:"question_id"`
	AnswerNum  null.Int16  `db:"answer_num" json:"answer_num"`
}
