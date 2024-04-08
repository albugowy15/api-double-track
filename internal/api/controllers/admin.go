package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	userModel "github.com/albugowy15/api-double-track/internal/pkg/models/user"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories/user"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/albugowy15/api-double-track/internal/pkg/validator"
)

// GetAdminProfile godoc
//
//	@Summary		Get admin profile
//	@Description	Get admin profile
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httputil.DataJsonResponse{data=schemas.Admin}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		404				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/admin/profile [get]
func GetAdminProfile(w http.ResponseWriter, r *http.Request) {
	userIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	adminId := userIdClaim.(string)
	a := user.GetAdminRepository()
	admin, err := a.GetAdminById(adminId)
	if err != nil {
		if err == sql.ErrNoRows {
			httputil.SendError(w, errors.New("profil admin tidak ditemukan"), http.StatusNotFound)
			return
		}
		log.Println(err)
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httputil.SendData(w, admin, http.StatusOK)
}

// UpdateAdminProfile godoc
//
//	@Summary		Update admin profile
//	@Description	Update admin profile
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.UpdateAdminRequest	true	"Update admin profile request body"
//	@Success		201				{object}	httputil.MessageJsonResponse
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/admin/profile [patch]
func UpdateAdminProfile(w http.ResponseWriter, r *http.Request) {
	userIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	adminId := userIdClaim.(string)
	var body userModel.UpdateAdminRequest
	httputil.GetBody(w, r, &body)
	err := validator.ValidateUpdateAdminRequest(body)
	if err != nil {
		httputil.SendError(w, err, http.StatusBadRequest)
		return
	}

	a := user.GetAdminRepository()
	_, err = a.GetAdminById(adminId)
	if err != nil {
		if err == sql.ErrNoRows {
			httputil.SendError(w, errors.New("profil admin tidak ditemukan"), http.StatusNotFound)
			return
		}
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	if err := a.UpdateAdminProfile(adminId, body); err != nil {
		log.Print(err)
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httputil.SendMessage(w, "berhasil memperbarui profil admin", http.StatusCreated)
}
