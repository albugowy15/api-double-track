package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/pkg/httpx"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var tokenAuth *jwtauth.JWTAuth

func Init(secret string) {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func GetAuth() *jwtauth.JWTAuth {
	return tokenAuth
}

func GetJwtClaim(r *http.Request, key string) (interface{}, error) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	value, ok := claims[key]
	if !ok {
		return "", fmt.Errorf("invalid token")
	}
	return value, nil
}

type JWTClaim struct {
	UserId   string
	Username string
	Role     string
	Email    string
	SchoolId string
}

func CreateToken(claim JWTClaim) string {
	jwtClaims := map[string]interface{}{
		"user_id":   claim.UserId,
		"username":  claim.Username,
		"role":      claim.Role,
		"email":     claim.Email,
		"school_id": claim.SchoolId,
	}
	jwtauth.SetIssuedNow(jwtClaims)
	// jwtauth.SetExpiryIn(jwtClaims, time.Duration(15*time.Minute))
	_, token, err := GetAuth().Encode(jwtClaims)
	if err != nil {
		log.Fatalf("error create token: %v", err)
	}
	return token
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			httpx.SendError(w, err, http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			httpx.SendError(w, errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
