package middelware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"portfolio-website/internals/service"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func exampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretStr := os.Getenv("MY_JWT_SECRETKEY")
		if secretStr == "" {
			fmt.Println("MY_JWT_SECRETKEY is not set")
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			fmt.Println("authorization header is missing")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "Bearer") {
			fmt.Println("authorization header format must be 'Bearer {token}'")
			return
		}

		access_token := headerParts[1]

		token, err := jwt.ParseWithClaims(access_token, &service.AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(secretStr), nil
		})
		if err != nil {
			log.Fatal(err)
		} else if claims, ok := token.Claims.(*service.AccessTokenClaims); ok {
			if claims.ExpiresAt.After(time.Now()) {
				
			}
		} else {
			log.Fatal("unknown claims type, cannot proceed")
		}

		next.ServeHTTP(w, r)
	})
}
