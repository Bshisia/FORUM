package middleware

import (
	"context"
	"database/sql"
	"forum/utils"
	"net/http"
	"time"
)

func oooAuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			session, err := utils.GetSession(db, cookie.Value)
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   "",
					Path:    "/",
					Expires: time.Now(),
				})
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
