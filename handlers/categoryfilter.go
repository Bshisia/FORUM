package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func FilterPostsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	posts, err := getPostsByCategory(category)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		log.Println("Error fetching posts by category:", err)
		return
	}

	for _, post := range posts {
		fmt.Fprintf(w, "Title: %s\nContent: %s\n\n", post.Title, post.Content)
	}
}

func getPostsByCategory(category string) ([]Post, error) {
	rows, err := db.Query("SELECT title, content FROM posts WHERE category = $1", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Title, &post.Content); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

type Post struct {
	Title   string
	Content string
}
