package handlers

import (
	"database/sql"
	"forum/utils"
	"net/http"
	"strconv"
)

func VotePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := r.Context().Value("session").(*utils.Session)
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	vote, _ := strconv.Atoi(r.FormValue("vote")) // 1 for like, -1 for dislike

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Check if user already voted
	var existingVote int
	err = tx.QueryRow(`
        SELECT vote FROM post_votes
        WHERE user_id = ? AND post_id = ?
    `, session.UserID, postID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// Insert new vote
		_, err = tx.Exec(`
            INSERT INTO post_votes (user_id, post_id, vote)
            VALUES (?, ?, ?)
        `, session.UserID, postID, vote)
	} else if err == nil {
		// Update existing vote
		_, err = tx.Exec(`
            UPDATE post_votes
            SET vote = ?
            WHERE user_id = ? AND post_id = ?
        `, vote, session.UserID, postID)
	}

	if err != nil {
		tx.Rollback()
		http.Error(w, "Error processing vote", http.StatusInternalServerError)
		return
	}

	// Update post likes/dislikes count
	_, err = tx.Exec(`
        UPDATE posts
        SET likes = (SELECT COUNT(*) FROM post_votes WHERE post_id = ? AND vote = 1),
            dislikes = (SELECT COUNT(*) FROM post_votes WHERE post_id = ? AND vote = -1)
        WHERE id = ?
    `, postID, postID, postID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Error updating vote counts", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}
