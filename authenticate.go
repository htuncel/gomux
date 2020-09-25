package main

import (
	"net/http"
	"strings"

	"main/utils"
)

// Authenticate is jwt middleware
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")

		if len(token) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}

		token = strings.Replace(token, "Bearer ", "", 1)
		_, errVerifyToken := utils.VerifyToken(token)

		if errVerifyToken != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + errVerifyToken.Error()))
			return
		}

		next.ServeHTTP(w, r)
	})
}
