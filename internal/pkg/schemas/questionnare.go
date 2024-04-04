package schemas

type Question struct {
	Code        string `json:"code"`
	Question    string `json:"question"`
	Category    string `json:"category"`
	Description int    `json:"description"`
	Id          int    `json:"id"`
	Number      int    `json:"number"`
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
	SchoolId                    string `json:"school_id"`
	Id                          int    `json:"id"`
	AlternativeId               int    `json:"alternative_id"`
	TotalOpenJobs               int    `json:"total_open_jobs"`
	EntrepreneurshipOpportunity int    `json:"entrepreneurship_opportunity"`
	Salary                      int    `json:"salary"`
}

type QuestionnareSettingRequest struct {
	AlternativeId               int `json:"alternative_id"`
	TotalOpenJobs               int `json:"total_open_jobs"`
	EntrepreneurshipOpportunity int `json:"entrepreneurship_opportunity"`
	Salary                      int `json:"salary"`
}
type QuestionnareSettingAlternative struct {
	Alternative                 string `json:"alternative"`
	Id                          int    `json:"id"`
	TotalOpenJobs               int    `json:"total_open_jobs"`
	EntrepreneurshipOpportunity int    `json:"entrepreneurship_opportunity"`
	Salary                      int    `db:"salary" json:"salary"`
}

type Answer struct {
	StudentId  string `json:"student_id"`
	Answer     string `json:"answer"`
	Id         int64  `json:"id"`
	QuestionId int    `json:"question_id"`
}

type QuestionnareStatusResponse struct {
	Status string `json:"status"`
}

type SubmitAnswerRequest struct {
	Answer string `json:"answer"`
	Id     int    `json:"id"`
	Number int    `json:"number"`
}
