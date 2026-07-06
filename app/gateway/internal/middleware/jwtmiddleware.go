package middleware

import (
	"net/http"
	"strings"

	"amigo-api/app/gateway/internal/config"

	"github.com/dgrijalva/jwt-go"
)

type JWTMiddleware struct {
	Config config.Auth
}

func NewJWTMiddleware(c config.Auth) *JWTMiddleware {
	return &JWTMiddleware{
		Config: c,
	}
}

// Whitelist of paths that don't require token validation
var whitelist = map[string]bool{
	"/api/user/third_login":      true,
	"/api/user/login":            true,
	"/api/admin/login":           true,
	"/api/admin/refresh":         true,
	"/api/sdk/message/send_code": true,
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"code":401,"msg":"` + msg + `"}`))
}

func (m *JWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 白名单前缀匹配，容错尾部斜杠
		for prefix := range whitelist {
			if strings.HasPrefix(r.URL.Path, prefix) {
				next(w, r)
				return
			}
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeUnauthorized(w, "Authorization header is required")
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthorized(w, "Authorization header format must be Bearer {token}")
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.Config.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			writeUnauthorized(w, "Invalid or expired token")
			return
		}

		// Token is valid, proceed to next handler
		next(w, r)
	}
}
