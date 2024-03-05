package user

type Student struct {
	Id          string `db:"id" json:"id"`
	Username    string `db:"username" json:"username"`
	Password    string `db:"password" json:"-"`
	Fullname    string `db:"fullname" json:"fullname"`
	Nisn        string `db:"nisn" json:"nisn"`
	Email       string `db:"email" json:"email"`
	PhoneNumber string `db:"phone_number" json:"phone_number"`
}
