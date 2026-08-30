package middelware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"portfolio-website/internals/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		secretStr := os.Getenv("MY_JWT_SECRETKEY")
		if secretStr == "" {
			log.Println("[ERROR] MY_JWT_SECRETKEY is not set in environment")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal Server Error",
			})
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Println("[INFO] Authentication failed: session_token cookie not found")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(model.Response{
					Success: false,
					Message: "Authentication required: session cookie not found",
				})
				return
			}
			log.Printf("[ERROR] Failed to retrieve session cookie: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal Server Error",
			})
			return
		}

		access_token := cookie.Value

		token, err := jwt.Parse(access_token, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretStr), nil
		})

		if err != nil || !token.Valid {
			log.Printf("[WARN] Authentication failed: invalid token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("[WARN] Authentication failed: unknown claims type")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Unknown claims type",
			})
			return
		}

		isValid := claims.VerifyExpiresAt(time.Now().Unix(), true)
		if !isValid {
			log.Println("[WARN] Authentication failed: token is expired")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Token is expired",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
