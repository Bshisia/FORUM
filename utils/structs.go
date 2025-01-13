package utils

import "time"

type user struct {
	ID int
	name string
	Email string
	UserName string
	PassWord string
}

type Post struct {
	ID int
	UserID int
	Title string
	Content string
	PostTime time.Time
	Likes int
	Dislikes int
	Comments int
}

type Comments struct {
	ID int
	PostID int
	UserID int
	Content string
	CommentTime time.Time
	Likes int
	Dislikes int
}

type Category struct {
	ID int
	Name string
}



