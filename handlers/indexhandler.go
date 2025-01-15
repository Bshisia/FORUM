// handlers/indexhandler.go
package handlers

import (
	"forum/utils"
	"log"
	"net/http"
	"text/template"
)

type TimelineData struct {
	User            *utils.User
	Posts           []utils.PostWithMetadata
	Categories      []utils.CategoryWithCount
	CurrentCategory *utils.Category
}

func TimelineHandler(w http.ResponseWriter, r *http.Request) {
	// Check for a valid session
	session, ok := r.Context().Value("session").(*utils.Session)
	if !ok || session == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	// Fetch user data
	user, err := utils.GetUserByID(db, session.UserID)
	if err != nil {
		log.Printf("Error fetching user: %v\n", err)
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	// If user is nil, redirect to signin
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	// Fetch categories with post counts
	categories, err := utils.GetCategoriesWithCount(db)
	if err != nil {
		log.Printf("Error fetching categories: %v\n", err)
		// Continue with empty categories rather than throwing an error
		categories = []utils.CategoryWithCount{}
	}

	// Fetch posts
	posts, err := utils.GetAllPosts(db)
	if err != nil {
		log.Printf("Error fetching posts: %v\n", err)
		// Continue with empty posts rather than throwing an error
		posts = []utils.PostWithMetadata{}
	}

	// Prepare data for the timeline page
	data := TimelineData{
		User:       user,
		Posts:      posts,
		Categories: categories,
	}

	// Parse both the base template and any required template
	tmpl, err := template.ParseFiles("static/templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	// Execute template with error handling
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}
