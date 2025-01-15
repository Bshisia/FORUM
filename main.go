package main

import (
    "forum/handlers"
    "forum/middleware"
    "forum/utils"
    "log"
    "net/http"
)

func main() {
    // Initialize database
    db := utils.InitialiseDB()
    defer db.Close()

    // Initialize handlers with database connection
    handlers.InitDB(db)

    // Create ServeMux for routing
    mux := http.NewServeMux()

    // Rest of your routes remain the same
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    mux.HandleFunc("/signup", handlers.SignUpHandler)
    mux.HandleFunc("/signin", handlers.SignInHandler)
    // mux.HandleFunc("/signout", handlers.SignOutHandler)
    mux.HandleFunc("/", handlers.HomeHandler)
    mux.HandleFunc("/post/create", handlers.CreatePostHandler)
    mux.HandleFunc("/post/", handlers.ViewPostHandler)
    mux.HandleFunc("/category/", handlers.CategoryHandler)
    mux.HandleFunc("/posts/created", handlers.UserPostsHandler)
    mux.HandleFunc("/posts/liked", handlers.LikedPostsHandler)
    mux.HandleFunc("/post/vote", handlers.VotePostHandler)
    mux.HandleFunc("/comment/vote", handlers.VoteCommentHandler)
    mux.HandleFunc("/post/comment", handlers.CreateCommentHandler)

    // Apply middleware
    handler := middleware.LoggingMiddleware(
        middleware.AuthMiddleware(db)(mux),
    )

    log.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatal(err)
    }
}