package routes

import (
	"forum/handlers"
	"net/http"
)

// setupRoutes handles the route registration
func SetupRoutes(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/signup", handlers.SignUpHandler)
	mux.HandleFunc("/signin", handlers.SignInHandler)
	// mux.HandleFunc("/signout", handlers.SignOutHandler)
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/timeline", handlers.TimelineHandler)
	mux.HandleFunc("/post/create", handlers.CreatePostHandler)
	mux.HandleFunc("/post/", handlers.ViewPostHandler)
	mux.HandleFunc("/category/", handlers.CategoryHandler)
	mux.HandleFunc("/posts/created", handlers.UserPostsHandler)
	mux.HandleFunc("/posts/liked", handlers.LikedPostsHandler)
	mux.HandleFunc("/post/vote", handlers.VotePostHandler)
	mux.HandleFunc("/comment/vote", handlers.VoteCommentHandler)
	mux.HandleFunc("/post/comment", handlers.CreateCommentHandler)
}
