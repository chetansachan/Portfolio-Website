package middelware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"portfolio-website/internals/model"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		secretStr := os.Getenv("MY_JWT_SECRETKEY")
		if secretStr == "" {
			fmt.Println("MY_JWT_SECRETKEY is not set")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(model.LoginResponse{
				Success: false,
				Message: "Internal Server Error",
			})
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			fmt.Println("authorization header is missing")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Authorization header is missing",
			})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "Bearer") {
			fmt.Println("Invalid Header Format")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Invalid Header Format",
			})
			return
		}

		access_token := headerParts[1]
		if access_token == "" || access_token == header {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Invalid Header Format",
			})
			return
		}

		token, err := jwt.Parse(access_token, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				fmt.Println("Unexpected Signing Method")
				w.WriteHeader(http.StatusBadGateway)
				json.NewEncoder(w).Encode(model.Response{
					Success: false,
					Message: "Internal Servor Error",
				})
			}
			return []byte(secretStr), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)
		isValid := claims.VerifyExpiresAt(time.Now().Unix(), true) 


		if err != nil || !token.Valid {
			fmt.Println("Invalid Token")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Invalid Token",
			})
			return 
		} else if !ok {
			fmt.Println("Unknown claims type, cannot proceed")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Unknown claims type",
			})
			return 
		} else if !isValid {
			fmt.Println("Token is Expired")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Token is Expired",
			})
			return 
		}

		next.ServeHTTP(w, r)
	})
}
