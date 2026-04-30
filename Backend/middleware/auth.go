package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"context"
	
)

var jwtKey = []byte("secret_key")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1. Get Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// 2. Extract token (remove "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
	claims, ok := token.Claims.(jwt.MapClaims)
if !ok {
	http.Error(w, "Invalid claims", http.StatusUnauthorized)
	return
}

email, ok := claims["email"].(string)
if !ok {
	http.Error(w, "Email not found in token", http.StatusUnauthorized)
	return
}
ctx := r.Context()
ctx = context.WithValue(ctx, "userEmail", email)
r = r.WithContext(ctx)

		// 4. Continue to actual handler
		next(w, r)
	}
}
