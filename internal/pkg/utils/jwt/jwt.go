package jwt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

func SetupAuth(secret string) {
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
}

func CreateToken(claim JWTClaim) string {
	jwtClaims := map[string]interface{}{
		"user_id":  claim.UserId,
		"username": claim.Username,
		"role":     claim.Role,
		"email":    claim.Email,
	}
	jwtauth.SetIssuedNow(jwtClaims)
	_, token, err := GetAuth().Encode(jwtClaims)
	if err != nil {
		log.Fatalf("error create token: %v", err)
	}
	return token
}
