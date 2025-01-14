package handlers

import (
	"forum/utils"
	"html/template"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("static/templates/create-post.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		session, _ := r.Context().Value("session").(*utils.Session)
		if session == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["categories"]

		// Insert post
		result, err := db.Exec(`
            INSERT INTO posts (user_id, title, content)
            VALUES (?, ?, ?)
        `, session.UserID, title, content)

		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		postID, _ := result.LastInsertId()

		// Insert categories
		for _, catID := range categories {
			categoryID, _ := strconv.Atoi(catID)
			_, err = db.Exec(`
                INSERT INTO post_categories (post_id, category_id)
                VALUES (?, ?)
            `, postID, categoryID)
		}

		http.Redirect(w, r, "/post/"+strconv.FormatInt(postID, 10), http.StatusSeeOther)
	}
}
