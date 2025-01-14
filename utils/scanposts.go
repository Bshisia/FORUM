package utils

import (
	"database/sql"
)

// scanPosts processes database rows and returns slice of PostWithMetadata
func ScanPosts(db *sql.DB, rows *sql.Rows) ([]PostWithMetadata, error) {
	var posts []PostWithMetadata

	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.PostTime,
			&post.Likes,
			&post.Dislikes,
			&post.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		// Get categories for this post
		catRows, err := db.Query(`
			SELECT c.id, c.name 
			FROM categories c
			JOIN post_categories pc ON c.id = pc.category_id
			WHERE pc.post_id = ?
		`, post.ID)
		if err != nil {
			return nil, err
		}
		defer catRows.Close()

		var categories []Category
		for catRows.Next() {
			var cat Category
			err := catRows.Scan(&cat.ID, &cat.Name)
			if err != nil {
				return nil, err
			}
			categories = append(categories, cat)
		}

		post.Categories = categories
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
