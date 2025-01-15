// handlers/userhandler.go
package handlers

import (
	"forum/utils"
	"html/template"
	"net/http"
)

func UserPostsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*utils.Session)
	if session == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	posts, err := utils.GetUserPosts(db, session.UserID)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	categories, err := utils.GetAllCategories(db)
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	user, _ := utils.GetUserByID(db, session.UserID)

	data := HomeData{
		Posts:      posts,
		Categories: categories,
		User:       user,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, data)
}

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*utils.Session)
	if session == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	posts, err := utils.GetLikedPosts(db, session.UserID)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	categories, err := utils.GetAllCategories(db)
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	user, _ := utils.GetUserByID(db, session.UserID)

	data := HomeData{
		Posts:      posts,
		Categories: categories,
		User:       user,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, data)
}