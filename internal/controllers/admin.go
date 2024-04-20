package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/internal/validator"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/crypto"
	"github.com/albugowy15/api-double-track/pkg/httpx"
	"golang.org/x/crypto/bcrypt"
)

// HandleGetAdminProfile godoc
//
//	@Summary		Get admin profile
//	@Description	Get admin profile
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.Admin}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		404				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/admin/profile [get]
func HandleGetAdminProfile(w http.ResponseWriter, r *http.Request) {
	userIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	adminId := userIdClaim.(string)
	admin, err := repositories.GetAdminById(adminId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("profil admin tidak ditemukan"), http.StatusNotFound)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendData(w, admin, http.StatusOK)
}

// HandlePatchAdminProfile godoc
//
//	@Summary		Update admin profile
//	@Description	Update admin profile
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.UpdateAdminRequest	true	"Update admin profile request body"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/admin/profile [patch]
func HandlePatchAdminProfile(w http.ResponseWriter, r *http.Request) {
	userIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	adminId := userIdClaim.(string)
	var body models.UpdateAdminRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	err := validator.ValidateUpdateAdminRequest(body)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	_, err = repositories.GetAdminById(adminId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("profil admin tidak ditemukan"), http.StatusNotFound)
			return
		}
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	if err := repositories.UpdateAdminProfile(adminId, body); err != nil {
		log.Print(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendMessage(w, "berhasil memperbarui profil admin", http.StatusCreated)
}

// HandlePatchAdminChangePassword godoc
//
//	@Summary		Change admin password
//	@Description	Change admin password
//	@Tags			Admin
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		models.ChangePasswordRequest	true	"Change admin password request body"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/admin/change-password [patch]
func HandlePatchAdminChangePassword(w http.ResponseWriter, r *http.Request) {
	userIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	adminId := userIdClaim.(string)

	var body models.ChangePasswordRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, httpx.ErrDecodeJsonBody, http.StatusBadRequest)
		return
	}

	if err := validator.ValidateChangePassword(body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	admin, err := repositories.GetAdminById(adminId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("data admin tidak ditemukan"), http.StatusNotFound)
			return
		}
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.OldPassword)); err != nil {
		log.Println("hashedPass:", admin.Password)
		log.Println(err)
		httpx.SendError(w, errors.New("password lama salah"), http.StatusBadRequest)
		return
	}

	hashedNewPassword, err := crypto.HashStr(body.NewPassword)
	if err != nil {
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusNotFound)
		return
	}

	if err := repositories.UpdateAdminPassword(adminId, hashedNewPassword); err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusNotFound)
		return
	}

	httpx.SendMessage(w, "berhasil mengubah password admin", http.StatusCreated)
}
