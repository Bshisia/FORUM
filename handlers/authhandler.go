package handlers

import (
	"database/sql"
	"forum/utils"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

// Store db connection globally for the package
var db *sql.DB

// InitDB initializes the database connection for all handlers
func InitDB(database *sql.DB) {
	db = database
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("static/templates/signup.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		username := r.FormValue("username") // Fixed: getting username from correct field
		password := r.FormValue("password")

		// Added validation
		if !utils.ValidateEmail(email) || !utils.ValidateUsername(username) || !utils.ValidatePassword(password) {
			http.Error(w, "Invalid Input", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
            INSERT INTO users (email, username, password, name)
            VALUES (?, ?, ?, ?)
        `, email, username, hashedPassword, name)

		if err != nil {
			http.Error(w, "Username or email already exists", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
}
