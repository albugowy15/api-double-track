package schemas

type Student struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	Fullname    string `json:"fullname"`
	Nisn        string `json:"nisn"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type StudentProfile struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Nisn        string `json:"nisn"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	School      string `json:"school"`
}

type AddStudentRequest struct {
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Nisn        string `json:"nisn"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type UpdateStudentRequest struct {
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Nisn        string `json:"nisn"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type DeleteStudentRequest struct {
	Id string `json:"id"`
}
