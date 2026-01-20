package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/auth"
)

type contextKey string

const AuthContextKey contextKey = "auth"

type AuthInfo struct {
	Authenticated bool
	Method        string // "api_key" or "jwt"
}

func Auth(apiKeyHash string, jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authInfo := AuthInfo{}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey != "" {
				if auth.ValidateAPIKey(apiKey, apiKeyHash) {
					authInfo.Authenticated = true
					authInfo.Method = "api_key"
				} else {
					response.Unauthorized(w, "invalid API key")
					return
				}
			}

			if !authInfo.Authenticated {
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					token := strings.TrimPrefix(authHeader, "Bearer ")
					claims, err := jwtManager.ValidateToken(token)
					if err != nil {
						response.Unauthorized(w, "invalid token")
						return
					}
					if !claims.IsAccessToken() {
						response.Unauthorized(w, "invalid token type")
						return
					}
					authInfo.Authenticated = true
					authInfo.Method = "jwt"
				}
			}

			if !authInfo.Authenticated {
				response.Unauthorized(w, "authentication required")
				return
			}

			ctx := context.WithValue(r.Context(), AuthContextKey, authInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalAuth(apiKeyHash string, jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authInfo := AuthInfo{}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey != "" && auth.ValidateAPIKey(apiKey, apiKeyHash) {
				authInfo.Authenticated = true
				authInfo.Method = "api_key"
			}

			if !authInfo.Authenticated {
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					token := strings.TrimPrefix(authHeader, "Bearer ")
					claims, err := jwtManager.ValidateToken(token)
					if err == nil && claims.IsAccessToken() {
						authInfo.Authenticated = true
						authInfo.Method = "jwt"
					}
				}
			}

			ctx := context.WithValue(r.Context(), AuthContextKey, authInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthInfo(ctx context.Context) AuthInfo {
	if info, ok := ctx.Value(AuthContextKey).(AuthInfo); ok {
		return info
	}
	return AuthInfo{}
}

func IsAuthenticated(ctx context.Context) bool {
	return GetAuthInfo(ctx).Authenticated
}
