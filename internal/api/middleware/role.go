package middleware

import (
	"errors"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

var (
	ErrTokenInvalid  = errors.New("token invalid")
	ErrAdminAccess   = errors.New("tidak memiliki akses admin")
	ErrStudentAccess = errors.New("tidak memiliki akses siswa")
)

func CheckAdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleClaim, err := jwt.GetJwtClaim(r, "role")
		if err != nil {
			httputil.SendError(w, ErrTokenInvalid, http.StatusUnauthorized)
			return
		}

		role := roleClaim.(string)
		if role != "admin" {
			httputil.SendError(w, ErrAdminAccess, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckStudentRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleClaim, err := jwt.GetJwtClaim(r, "role")
		if err != nil {
			httputil.SendError(w, ErrTokenInvalid, http.StatusUnauthorized)
			return
		}

		role := roleClaim.(string)
		if role != "student" {
			httputil.SendError(w, ErrStudentAccess, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
