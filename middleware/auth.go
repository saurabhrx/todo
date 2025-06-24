package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
	"time"
)

type ContextKey string

const (
	userContext ContextKey = "userKey"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func GenerateAccessToken(userID string) (string, error) {
	accessClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessJWT.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
func GenerateRefreshToken(userID string) (string, error) {
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshJWT.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return refreshToken, nil

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//session based authentication

		//sessionID := r.Header.Get("session-id")
		//if sessionID == "" {
		//	http.Error(w, "unauthorized user", http.StatusUnauthorized)
		//	return
		//}
		//
		//userID, err := dbHelper.ValidateSession(sessionID)
		//if err != nil {
		//	http.Error(w, "invalid or expired session", http.StatusUnauthorized)
		//	return
		//}

		// jwt based authentication

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or malformed token", http.StatusUnauthorized)
			return
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		accessClaims := &Claims{}
		token, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		// access token is valid
		if err == nil && token.Valid {
			ctx := context.WithValue(r.Context(), userContext, accessClaims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else { // access token invalid
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

	})
}

func UserContext(r *http.Request) string {
	if user, ok := r.Context().Value(userContext).(string); ok && user != "" {
		return user
	}
	return ""

}
