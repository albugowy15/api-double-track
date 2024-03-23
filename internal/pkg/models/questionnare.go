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

type QuestionResponse struct {
	Question string   `json:"question"`
	Type     string   `json:"type"`
	MinText  string   `json:"min_text"`
	MaxText  string   `json:"max_text"`
	Options  []string `json:"options"`
	Id       int      `json:"id"`
	Number   int      `json:"number"`
}

type QuestionnareSetting struct {
	SchoolId                    string     `db:"school_id" json:"school_id"`
	Id                          int        `db:"id" json:"id"`
	AlternativeId               int        `db:"alternative_id" json:"alternative_id"`
	TotalOpenJobs               null.Int16 `db:"total_open_jobs" json:"total_open_jobs"`
	EntrepreneurshipOpportunity null.Int16 `db:"entrepreneurship_opportunity" json:"entrepreneurship_opportunity"`
	Salary                      null.Int16 `db:"salary" json:"salary"`
}

type QuestionnareSettingAlternative struct {
	Alternative                 string     `db:"alternative" json:"alternative"`
	Id                          int        `db:"id" json:"id"`
	TotalOpenJobs               null.Int16 `db:"total_open_jobs" json:"total_open_jobs"`
	EntrepreneurshipOpportunity null.Int16 `db:"entrepreneurship_opportunity" json:"entrepreneurship_opportunity"`
	Salary                      null.Int16 `db:"salary" json:"salary"`
}

type Answer struct {
	StudentId  string      `db:"student_id" json:"student_id"`
	Answer     null.String `db:"answer" json:"answer"`
	Id         int64       `db:"id" json:"id"`
	QuestionId int         `db:"question_id" json:"question_id"`
}

type SubmitAnswerRequest struct {
	Id     int    `json:"id"`
	Number int    `json:"number"`
	Answer string `json:"answer"`
}
