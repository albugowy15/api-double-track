package models

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}
