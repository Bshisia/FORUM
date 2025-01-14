// handlers/homehandler.go
package handlers

import (
	"forum/utils"
	"html/template"
	"net/http"
)

type HomeData struct {
	Posts           []utils.PostWithMetadata
	Categories      []utils.Category
	User            *utils.User
	CurrentCategory utils.Category
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	session, _ := r.Context().Value("session").(*utils.Session)
	var user *utils.User
	if session != nil {
		user, _ = utils.GetUserByID(db, session.UserID)
	}

	categories, err := utils.GetAllCategories(db)
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	posts, err := utils.GetAllPosts(db)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		Posts:      posts,
		Categories: categories,
		User:       user,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/home.html"))
	tmpl.Execute(w, data)
}
