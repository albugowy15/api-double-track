package models

import "github.com/guregu/null/v5"

type Statistic struct {
	RegisteredStudents    int64      `db:"registered_students" json:"registered_students"`
	QuestionnareCompleted int64      `db:"questionnare_completed" json:"questionnare_completed"`
	ConsistencyAvg        null.Float `db:"consistency_avg" json:"consistency_avg"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
