package user

import (
	"log"

	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models/user"
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
	return admin, err
}

func (r *AdminRepository) GetAdminById(adminId string) (user.Admin, error) {
	admin := user.Admin{}
	err := db.GetDb().Get(&admin, "SELECT id, username, email, phone_number FROM admins WHERE id = $1", adminId)
	return admin, err
}

func (r *AdminRepository) UpdateAdminProfile(adminId string, data user.UpdateAdminRequest) error {
	tx, err := db.GetDb().Beginx()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE admins SET username = $1 WHERE id = $2", data.Username, adminId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(data.Email) > 0 {
		_, err = tx.Exec("UPDATE admins SET email = $1 WHERE id = $2", data.Email, adminId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(data.PhoneNumber) > 0 {
		_, err = tx.Exec("UPDATE admins SET phone_number = $1 WHERE id = $2", data.PhoneNumber, adminId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
