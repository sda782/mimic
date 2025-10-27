package auth

import (
	"mimic/backend/database"
	"mimic/backend/types"
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := database.GetUserByToken(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(types.WithUser(r.Context(), user))
		next.ServeHTTP(w, r)
	}
}
