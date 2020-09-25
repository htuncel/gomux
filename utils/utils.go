package utils

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-playground/validator/v10"

	"main/configs"
)

var (
	// Validate is a global variable to validate structs
	Validate *validator.Validate
)

func init() {
	Validate = validator.New()
}

// VerifyToken to verify the jwt token validity
func VerifyToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return configs.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

// GenerateToken to generate jwt token
// TODO will take user information later to sign with user data
func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":       0,
		"identity": "00000000000",
		"title":    "",
		"iat":      time.Now().Unix(),
		"exp":      (time.Now().Local().Add(time.Second * time.Duration(86400)).Unix()),
	})

	return token.SignedString(configs.Secret)
}

// GetToken to extract jwt token from header
func GetToken(w http.ResponseWriter, r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Replace(reqToken, "Bearer ", "", 1)
	return splitToken
}

// GetClientIP to extract client ip from header
func GetClientIP(w http.ResponseWriter, r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if parts := strings.Split(xff, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return r.RemoteAddr
	}

	return host
}
