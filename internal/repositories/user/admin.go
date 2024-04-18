package user

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models/user"
)

func GetAdminByUsername(username string) (user.Admin, error) {
	admin := user.Admin{}
	err := db.AppDB.Get(&admin, "SELECT id, username, email, password, phone_number FROM admins WHERE username = $1", username)
	return admin, err
}

func GetAdminById(adminId string) (user.Admin, error) {
	admin := user.Admin{}
	err := db.AppDB.Get(&admin, "SELECT id, username, email, phone_number FROM admins WHERE id = $1", adminId)
	return admin, err
}

func UpdateAdminProfile(adminId string, data user.UpdateAdminRequest) error {
	tx, err := db.AppDB.Beginx()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE admins SET username = $1 WHERE id = $2", data.Username, adminId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE admins SET email = $1 WHERE id = $2", data.Email, adminId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE admins SET phone_number = $1 WHERE id = $2", data.PhoneNumber, adminId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
