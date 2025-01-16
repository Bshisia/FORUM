package handlers

import (
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

type ProfileData struct {
	User  *utils.User
	Posts []utils.PostWithMetadata
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Check for valid session
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

	// Fetch user's posts
	posts, err := utils.GetPostsByUserID(db, user.ID)
	if err != nil {
		log.Printf("Error fetching posts: %v\n", err)
		http.Error(w, "Error loading user posts", http.StatusInternalServerError)
		return
	}

	// Prepare data for the profile page
	data := ProfileData{
		User:  user,
		Posts: posts,
	}

	// Parse and execute the profile template
	tmpl := template.Must(template.ParseFiles("static/templates/profile.html"))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
