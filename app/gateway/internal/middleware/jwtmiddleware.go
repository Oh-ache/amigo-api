package middleware

import (
	"net/http"
	"strings"

	"amigo-api/app/gateway/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JWTMiddleware struct {
	Config config.Auth
}

func NewJWTMiddleware(c config.Auth) *JWTMiddleware {
	return &JWTMiddleware{
		Config: c,
	}
}

type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Whitelist of paths that don't require token validation
var whitelist = []string{
	"/api/user/third_login",
}

func (m *JWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if path is in whitelist
		for _, path := range whitelist {
			if r.URL.Path == path {
				next(w, r)
				return
			}
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.OkJsonCtx(r.Context(), w, CommonResp{
				Code: 401,
				Msg:  "Authorization header is required",
				Data: nil,
			})
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			httpx.OkJsonCtx(r.Context(), w, CommonResp{
				Code: 401,
				Msg:  "Authorization header format must be Bearer {token}",
				Data: nil,
			})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.Config.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			httpx.OkJsonCtx(r.Context(), w, CommonResp{
				Code: 401,
				Msg:  "Invalid or expired token",
				Data: nil,
			})
			return
		}

		// Token is valid, proceed to next handler
		next(w, r)
	}
}
