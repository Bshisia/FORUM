package handlers

import (
	"forum/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type CreatePostData struct {
	Categories []utils.Category
	User       *utils.User
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Check for valid session
	session, ok := r.Context().Value("session").(*utils.Session)
	if !ok || session == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		// Fetch categories for the form
		categories, err := utils.GetAllCategories(db)
		if err != nil {
			log.Printf("Error fetching categories: %v\n", err)
			http.Error(w, "Error loading categories", http.StatusInternalServerError)
			return
		}

		// Get user data
		user, err := utils.GetUserByID(db, session.UserID)
		if err != nil {
			log.Printf("Error fetching user: %v\n", err)
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		// Prepare template data
		data := CreatePostData{
			Categories: categories,
			User:       user,
		}

		// Parse and execute template
		tmpl, err := template.ParseFiles("static/templates/createpost.html")
		if err != nil {
			log.Printf("Error parsing template: %v\n", err)
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error executing template: %v\n", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == "POST" {
		// Parse the form
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v\n", err)
			http.Error(w, "Error processing form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["categories"] // This gets all selected categories

		// Validate input
		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error starting transaction: %v\n", err)
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		// Insert post
		result, err := tx.Exec(`
			INSERT INTO posts (user_id, title, content)
			VALUES (?, ?, ?)
		`, session.UserID, title, content)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting post: %v\n", err)
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		postID, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			log.Printf("Error getting last insert ID: %v\n", err)
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		// Insert categories
		for _, catID := range categories {
			categoryID, err := strconv.Atoi(catID)
			if err != nil {
				tx.Rollback()
				log.Printf("Error converting category ID: %v\n", err)
				continue
			}

			_, err = tx.Exec(`
				INSERT INTO post_categories (post_id, category_id)
				VALUES (?, ?)
			`, postID, categoryID)
			if err != nil {
				tx.Rollback()
				log.Printf("Error inserting category: %v\n", err)
				http.Error(w, "Error creating post", http.StatusInternalServerError)
				return
			}
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v\n", err)
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		// Redirect to the new post
		return
	}

	// If neither GET nor POST
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}