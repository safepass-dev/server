package middlewares

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/logging"
	"github.com/safepass/server/pkg/models"
)

type AuthMiddleware struct {
	logger *logging.Logger
	config config.Config
}

func NewAuthMiddleware(logger *logging.Logger, config config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
		config: config,
	}
}

func httpError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := models.Response{
		Status:     code,
		StatusText: http.StatusText(code),
		Data:       map[string]string{"message": error},
	}

	json.NewEncoder(w).Encode(response)
}

func (m *AuthMiddleware) AuthMiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpError(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpError(w, "Authorization header format must be 'Bearer <token>'", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		secretKey, err := m.config.GetJWTSecretKey()
		if err != nil {
			httpError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		claims, err := validateToken(token, secretKey)
		if err != nil {
			httpError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(tokenString string, secretKey *ecdsa.PrivateKey) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return &secretKey.PublicKey, nil
	})

	fmt.Println("Girdi0")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Girdi1")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	fmt.Println("Girdi2")

	return nil, fmt.Errorf("Invalid token")
}
