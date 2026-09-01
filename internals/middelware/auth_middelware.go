package middelware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"portfolio-website/internals/model"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	UserEmailKey   contextKey = "user_email"
	UserIsAdminKey contextKey = "user_is_admin"
)

func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}

func GetUserIsAdmin(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(UserIsAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}

// AuthMiddleware validates the access_token JWT and injects user claims into the request context.
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

		cookie, err := r.Cookie("access_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Println("[INFO] cookie not found")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(model.Response{
					Success: false,
					Message: "Cookie not found",
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

		accessToken := cookie.Value
		var claims model.AccessTokenClaims

		token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretStr), nil
		})

		if err != nil || !token.Valid {
			log.Printf("[WARN] Authentication failed: invalid or expired token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Invalid or expired token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserIsAdminKey, claims.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware verifies that the user is authenticated and is an admin (is_admin == true).
func AdminMiddleware(next http.Handler) http.Handler {
	return AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := GetUserIsAdmin(r.Context())
		if !isAdmin {
			log.Printf("[WARN] Authorization failed: admin privileges required for %s", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Forbidden: Only Admin is Allowed",
			})
			return
		}

		next.ServeHTTP(w, r)
	}))
}
