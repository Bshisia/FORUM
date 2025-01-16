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

func GetPostsByUserID(db *sql.DB, userID int) ([]PostWithMetadata, error) {
	query := `
		SELECT p.id, p.title, p.created_at, u.username 
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		if err := rows.Scan(&post.ID, &post.Title, &post.PostTime, &post.Username); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
