// handlers/commenthandler.go
package handlers

import (
	"database/sql"
	"forum/utils"
	"net/http"
	"strconv"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := r.Context().Value("session").(*utils.Session)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	// Insert comment
	_, err = db.Exec(`
        INSERT INTO comments (post_id, user_id, content)
        VALUES (?, ?, ?)
    `, postID, session.UserID, content)

	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}

func VoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := r.Context().Value("session").(*utils.Session)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
	vote, _ := strconv.Atoi(r.FormValue("vote")) // 1 for like, -1 for dislike

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Check if user already voted
	var existingVote int
	err = tx.QueryRow(`
        SELECT vote FROM comment_votes
        WHERE user_id = ? AND comment_id = ?
    `, session.UserID, commentID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// Insert new vote
		_, err = tx.Exec(`
            INSERT INTO comment_votes (user_id, comment_id, vote)
            VALUES (?, ?, ?)
        `, session.UserID, commentID, vote)
	} else if err == nil {
		// Update existing vote
		_, err = tx.Exec(`
            UPDATE comment_votes
            SET vote = ?
            WHERE user_id = ? AND comment_id = ?
        `, vote, session.UserID, commentID)
	}

	if err != nil {
		tx.Rollback()
		http.Error(w, "Error processing vote", http.StatusInternalServerError)
		return
	}

	// Update comment likes/dislikes count
	_, err = tx.Exec(`
        UPDATE comments
        SET likes = (SELECT COUNT(*) FROM comment_votes WHERE comment_id = ? AND vote = 1),
            dislikes = (SELECT COUNT(*) FROM comment_votes WHERE comment_id = ? AND vote = -1)
        WHERE id = ?
    `, commentID, commentID, commentID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Error updating vote counts", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}
