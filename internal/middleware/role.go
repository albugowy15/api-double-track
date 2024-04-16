package middleware

import (
	"errors"
	"net/http"

	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
)

var (
	ErrTokenInvalid  = errors.New("token invalid")
	ErrAdminAccess   = errors.New("tidak memiliki akses admin")
	ErrStudentAccess = errors.New("tidak memiliki akses siswa")
)

func CheckAdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleClaim, err := auth.GetJwtClaim(r, "role")
		if err != nil {
			httpx.SendError(w, ErrTokenInvalid, http.StatusUnauthorized)
			return
		}

		role := roleClaim.(string)
		if role != "admin" {
			httpx.SendError(w, ErrAdminAccess, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckStudentRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleClaim, err := auth.GetJwtClaim(r, "role")
		if err != nil {
			httpx.SendError(w, ErrTokenInvalid, http.StatusUnauthorized)
			return
		}

		role := roleClaim.(string)
		if role != "student" {
			httpx.SendError(w, ErrStudentAccess, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
