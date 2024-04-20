package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func GetAdminByUsername(username string) (models.Admin, error) {
	admin := models.Admin{}
	err := db.AppDB.Get(&admin, "SELECT id, username, email, password, phone_number FROM admins WHERE username = $1", username)
	return admin, err
}

func GetAdminById(adminId string) (models.Admin, error) {
	admin := models.Admin{}
	err := db.AppDB.Get(&admin, "SELECT id, username, email, password, phone_number FROM admins WHERE id = $1", adminId)
	if err != nil {
		log.Println("db err: ", err)
		return admin, err
	}
	return admin, nil
}

func UpdateAdminProfile(adminId string, data models.UpdateAdminRequest) error {
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

func UpdateAdminPassword(adminId string, hashedPassword string) error {
	_, err := db.AppDB.Exec(`UPDATE admins SET password = $1 WHERE id = $2`, hashedPassword, adminId)
	if err != nil {
		log.Println("db err: ", err)
		return err
	}
	return nil
}
