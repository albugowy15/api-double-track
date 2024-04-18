package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	userModel "github.com/albugowy15/api-double-track/internal/models/user"
	"github.com/albugowy15/api-double-track/internal/repositories/user"
	"github.com/albugowy15/api-double-track/internal/validator"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
)

// HanldeGetAdminProfile godoc
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
func HanldeGetAdminProfile(w http.ResponseWriter, r *http.Request) {
	userIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	adminId := userIdClaim.(string)
	admin, err := user.GetAdminById(adminId)
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
	var body userModel.UpdateAdminRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	err := validator.ValidateUpdateAdminRequest(body)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	_, err = user.GetAdminById(adminId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("profil admin tidak ditemukan"), http.StatusNotFound)
			return
		}
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	if err := user.UpdateAdminProfile(adminId, body); err != nil {
		log.Print(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendMessage(w, "berhasil memperbarui profil admin", http.StatusCreated)
}
