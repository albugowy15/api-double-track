package user

type Admin struct {
	Id          string `db:"id"`
	Username    string `db:"username"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	PhoneNumber string `db:"phone_number"`
}
