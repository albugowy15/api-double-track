package user

import (
	"github.con/albugowy15/api-double-track/internal/pkg/db"
	"github.con/albugowy15/api-double-track/internal/pkg/models/user"
)

type AdminRepository struct{}

var adminRepository *AdminRepository

func GetAdminRepository() *AdminRepository {
	if adminRepository == nil {
		adminRepository = &AdminRepository{}
	}
	return adminRepository
}

func (r *AdminRepository) GetAdminByUsername(username string) (user.Admin, error) {
	admin := user.Admin{}
	err := db.GetDb().Get(&admin, "SELECT id, username, email, password, phone_number FROM admins WHERE username = $1", username)
	if err != nil {
		return admin, err
	}
	return admin, nil
}
