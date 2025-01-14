package handlers

import (
	"forum/utils"
	"net/http"
	"text/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*utils.Session)
	if session == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	user, err := utils.GetUserByID(db, session.UserID)
	if err != nil {
		http.Error(w, "Error loading user data", http.StatusInternalServerError)
		return
	}

	posts, err := utils.GetAllPosts(db)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		User:  user,
		Posts: posts,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, data)
}
