// utils/database_utils.go
package utils

import (
	"database/sql"
)

// GetUserByID retrieves a user by their ID
func GetUserByID(db *sql.DB, userID int) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT id, email, username, name, password 
		FROM users 
		WHERE id = ?
	`, userID).Scan(&user.ID, &user.Email, &user.UserName, &user.Name, &user.PassWord)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllCategories retrieves all categories
func GetAllCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(&cat.ID, &cat.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetAllPosts retrieves all posts with metadata
func GetAllPosts(db *sql.DB) ([]PostWithMetadata, error) {
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
		ORDER BY p.post_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanPosts(db, rows)
}

// GetUserPosts retrieves all posts by a specific user
func GetUserPosts(db *sql.DB, userID int) ([]PostWithMetadata, error) {
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
		WHERE p.user_id = ?
		ORDER BY p.post_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanPosts(db, rows)
}

// GetLikedPosts retrieves all posts liked by a specific user
func GetLikedPosts(db *sql.DB, userID int) ([]PostWithMetadata, error) {
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
		JOIN post_votes pv ON p.id = pv.post_id
		WHERE pv.user_id = ? AND pv.vote = 1
		ORDER BY p.post_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanPosts(db, rows)
}

func GetCategoriesWithCount(db *sql.DB) ([]CategoryWithCount, error) {
	rows, err := db.Query(`
		SELECT 
			c.id,
			c.name,
			COUNT(pc.post_id) as post_count
		FROM categories c
		LEFT JOIN post_categories pc ON c.id = pc.category_id
		GROUP BY c.id, c.name
		ORDER BY post_count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryWithCount
	for rows.Next() {
		var cat CategoryWithCount
		err := rows.Scan(&cat.ID, &cat.Name, &cat.PostCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetPostsByCategory retrieves all posts for a specific category
func GetPostsByCategory(db *sql.DB, categoryID int) ([]PostWithMetadata, error) {
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
		return nil, err
	}
	defer rows.Close()

	return ScanPosts(db, rows)
}