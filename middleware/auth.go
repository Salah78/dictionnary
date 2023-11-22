// middleware/authmiddleware.go
package middleware

import (
	"net/http"
	"strings"

)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Token d'authentification manquant", http.StatusUnauthorized)
			return
		}

		if !isValidToken(token) {
			http.Error(w, "Token d'authentification invalide", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidToken(token string) bool {
	// Pour cet exemple, nous supposons que le jeton est valide s'il n'est pas vide.
	return strings.TrimSpace(token) != ""
}
