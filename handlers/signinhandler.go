package handlers

import (
	"database/sql"
	"forum/utils"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type SignInData struct {
	Username      string
	UsernameError string
	PasswordError string
	GeneralError  string
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/templates/signin.html"))

	if r.Method == "GET" {
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		data := SignInData{}

		// Get form values
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))

		// Validate input
		if username == "" {
			data.UsernameError = "Username is required"
		}
		if password == "" {
			data.PasswordError = "Password is required"
		}

		// If there are validation errors, return early
		if data.UsernameError != "" || data.PasswordError != "" {
			data.Username = username // Preserve the username input
			tmpl.Execute(w, data)
			return
		}

		// Query the database
		var user utils.User
		err := db.QueryRow(`
            SELECT id, password 
            FROM users 
            WHERE username = ?
        `, username).Scan(&user.ID, &user.PassWord)

		// Handle database errors
		if err != nil {
			if err == sql.ErrNoRows {
				data.GeneralError = "Invalid username or password"
			} else {
				data.GeneralError = "An error occurred. Please try again later."
			}
			data.Username = username // Preserve the username input
			tmpl.Execute(w, data)
			return
		}

		// Verify password
		if !utils.CheckPasswordsHash(password, user.PassWord) {
			data.GeneralError = "Invalid username or password"
			data.Username = username // Preserve the username input
			tmpl.Execute(w, data)
			return
		}

		// Create session
		sessionToken, err := utils.CreateSession(db, user.ID)
		if err != nil {
			data.GeneralError = "An error occurred. Please try again later."
			data.Username = username // Preserve the username input
			tmpl.Execute(w, data)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true, // Enable for HTTPS
			SameSite: http.SameSiteStrictMode,
		})

		// Redirect to timeline page
		http.Redirect(w, r, "/timeline", http.StatusSeeOther)
	}
}


// func SignOutHandler(w http.ResponseWriter, r *http.Request) {
// 	// Clear the session from database if needed
// 	if cookie, err := r.Cookie("session_token"); err == nil {
// 		utils.DeleteSession(db, cookie.Value)
// 	}

// 	// Clear the cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "session_token",
// 		Value:    "",
// 		Path:     "/",
// 		Expires:  time.Now().Add(-time.Hour),
// 		HttpOnly: true,
// 		Secure:   true, // Enable for HTTPS
// 		SameSite: http.SameSiteStrictMode,
// 	})

// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
