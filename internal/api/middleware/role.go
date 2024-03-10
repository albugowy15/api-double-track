package middleware

import (
	"net/http"

	"github.con/albugowy15/api-double-track/internal/pkg/utils"
	"github.con/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

func CheckAdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleClaim, err := jwt.GetJwtClaim(r, "role")
		if err != nil {
			utils.SendError(w, "token invalid", http.StatusUnauthorized)
			return
		}

		role := roleClaim.(string)
		if role != "admin" {
			utils.SendError(w, "tidak memiliki akses admin", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckStudentRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleClaim, err := jwt.GetJwtClaim(r, "role")
		if err != nil {
			utils.SendError(w, "token invalid", http.StatusUnauthorized)
			return
		}

		role := roleClaim.(string)
		if role != "admin" {
			utils.SendError(w, "tidak memiliki akses siswa", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
