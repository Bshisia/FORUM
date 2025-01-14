package handlers

import (
	"forum/utils"
	"html/template"
	"net/http"
	"time"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("static/templates/signin.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user utils.User
		err := db.QueryRow(`
			SELECT id, password 
			FROM users 
			WHERE username = ?
		`, username).Scan(&user.ID, &user.PassWord)

		if err != nil || !utils.CheckPasswordsHash(password, user.PassWord) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		sessionToken, err := utils.CreateSession(db, user.ID)
		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

