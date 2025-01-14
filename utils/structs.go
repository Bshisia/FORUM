package utils

import "time"

type User struct {
	ID       int
	Name     string
	Email    string
	UserName string
	PassWord string
}

type Post struct {
	ID       int
	UserID   int
	Title    string
	Content  string
	PostTime time.Time
	Likes    int
	Dislikes int
	Comments int
}

type Comment struct {
	ID          int
	PostID      int
	UserID      int
	Content     string
	CommentTime time.Time
	Likes       int
	Dislikes    int
}

type Category struct {
	ID   int
	Name string
}

type PostWithMetadata struct {
	Post
	Username     string
	Categories   []Category
	CommentCount int
}

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

type PostPageData struct {
	Post     PostWithMetadata
	Comments []Comment
	User     *User
}
