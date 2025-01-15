package handlers

import (
	"database/sql"
	"forum/utils"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type SignUpErrors struct {
    NameError     string
    EmailError    string
    UsernameError string
    PasswordError string
    GeneralError  string
}

type SignUpData struct {
    Errors SignUpErrors
    Name     string
    Email    string
    Username string
}

// Store db connection globally for the package
var db *sql.DB

// InitDB initializes the database connection for all handlers
func InitDB(database *sql.DB) {
	db = database
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl := template.Must(template.ParseFiles("static/templates/signup.html"))
        tmpl.Execute(w, &SignUpData{})
        return
    }

    if r.Method == "POST" {
        data := SignUpData{
            Name:     r.FormValue("name"),
            Email:    r.FormValue("email"),
            Username: r.FormValue("username"),
        }
        errors := SignUpErrors{}
        hasError := false

        // Validate name
        if data.Name == "" {
            errors.NameError = "Name is required"
            hasError = true
        }

        // Validate email
        if !utils.ValidateEmail(data.Email) {
            errors.EmailError = "Invalid email format"
            hasError = true
        }

        // Validate username
        if !utils.ValidateUsername(data.Username) {
            errors.UsernameError = "Username must be between 3 and 30 characters"
            hasError = true
        }

        // Validate password
        password := r.FormValue("password")
        confirmPassword := r.FormValue("confirm_password")

        if !utils.ValidatePassword(password) {
            errors.PasswordError = "Password must be at least 8 characters"
            hasError = true
        } else if password != confirmPassword {
            errors.PasswordError = "Passwords do not match"
            hasError = true
        }

        if hasError {
            data.Errors = errors
            tmpl := template.Must(template.ParseFiles("static/templates/signup.html"))
            tmpl.Execute(w, data)
            return
        }

        // Hash password and create user
        hashedPassword, err := utils.HashPassword(password)
        if err != nil {
            errors.GeneralError = "Internal Server Error"
            data.Errors = errors
            tmpl := template.Must(template.ParseFiles("static/templates/signup.html"))
            tmpl.Execute(w, data)
            return
        }

        _, err = db.Exec(`
            INSERT INTO users (email, username, password, name)
            VALUES (?, ?, ?, ?)
        `, data.Email, data.Username, hashedPassword, data.Name)

        if err != nil {
            errors.GeneralError = "Username or email already exists"
            data.Errors = errors
            tmpl := template.Must(template.ParseFiles("static/templates/signup.html"))
            tmpl.Execute(w, data)
            return
        }

        http.Redirect(w, r, "/signin", http.StatusSeeOther)
    }
}