package handlers

import (
	"forum/utils"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	session, _ := r.Context().Value("session").(*utils.Session)
	var user *utils.User
	if session != nil {
		user, _ = utils.GetUserByID(db, session.UserID)
	}

	rows, err := db.Query(`
		SELECT 
			p.id,
			p.user_id,
			p.title,
			p.content,
			p.post_at,
			p.likes,
			p.dislikes,
			u.username,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, postID)
	if err != nil {
		http.Error(w, "Error loading post", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts, err := utils.ScanPosts(db, rows)
	if err != nil || len(posts) == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Get comments
	rows, err = db.Query(`
		SELECT 
			c.id,
			c.post_id,
			c.user_id,
			c.content,
			c.comment_at,
			c.likes,
			c.dislikes,
			u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.comment_at DESC
	`, postID)
	if err != nil {
		http.Error(w, "Error loading comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []utils.Comment
	for rows.Next() {
		var comment utils.Comment
		var username string
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CommentTime,
			&comment.Likes,
			&comment.Dislikes,
			&username,
		)
		if err != nil {
			http.Error(w, "Error scanning comments", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	data := utils.PostPageData{
		Post:     posts[0],
		Comments: comments,
		User:     user,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/post.html"))
	tmpl.Execute(w, data)
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/category/"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	session, _ := r.Context().Value("session").(*utils.Session)
	var user *utils.User
	if session != nil {
		user, _ = utils.GetUserByID(db, session.UserID)
	}

	// Get category info
	var currentCategory utils.Category
	err = db.QueryRow("SELECT id, name FROM categories WHERE id = ?", categoryID).
		Scan(&currentCategory.ID, &currentCategory.Name)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// Get posts in this category
	rows, err := db.Query(`
		SELECT 
			p.id,
			p.user_id,
			p.title,
			p.content,
			p.post_at,
			p.likes,
			p.dislikes,
			u.username,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY p.post_at DESC
	`, categoryID)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts, err := utils.ScanPosts(db, rows)
	if err != nil {
		http.Error(w, "Error scanning posts", http.StatusInternalServerError)
		return
	}

	categories, err := utils.GetAllCategories(db)
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		Posts:           posts,
		Categories:      categories,
		User:            user,
		CurrentCategory: currentCategory,
	}

	tmpl := template.Must(template.ParseFiles("static/templates/home.html"))
	tmpl.Execute(w, data)
}
