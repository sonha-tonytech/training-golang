package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

var jwtSecret = []byte(os.Getenv("JWT_KEY"))

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VadidateLoginUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginUser := UserLogin{
			UserName: r.URL.Query().Get("userName"),
			Password: r.URL.Query().Get("password"),
		}

		validate := validator.New()

		err := validate.Struct(loginUser)

		if err != nil {
			http.Error(w, "Username or Password is missing", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VadidateRegisterUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const userContextKey contextKey = "registerUser"
		var registerUser User

		err := json.NewDecoder(r.Body).Decode(&registerUser)

		if err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		validate := validator.New()

		err = validate.Struct(registerUser)

		if err != nil {
			http.Error(w, "Some mandatory fields are missing", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, registerUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
