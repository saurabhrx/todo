package middleware

import (
	"context"
	"myTodo/database/dbHelper"
	"net/http"
)

type ContextKey string

const (
	userContext ContextKey = "userKey"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("session-id")
		if sessionID == "" {
			http.Error(w, "unauthorized user", http.StatusUnauthorized)
			return
		}

		userID, err := dbHelper.ValidateSession(sessionID)
		if err != nil {
			http.Error(w, "invalid or expired session", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContext, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) string {
	if user, ok := r.Context().Value(userContext).(string); ok && user != "" {
		return user
	}
	return ""
}
