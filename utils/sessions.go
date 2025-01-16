package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"
)

func GenerateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func CreateSession(db *sql.DB, userID int) (string, error) {
	sessionToken := GenerateSessionToken()
	ExpiresAt := time.Now().Add(24 * time.Hour)

	_, err := db.Exec(`
	INSERT INTO sessions(id, user_id, expires_at)
	VALUES (?, ?, ?)
	`, sessionToken, userID, ExpiresAt)
	return sessionToken, err
}

func GetSession(db *sql.DB, sessionToken string) (*Session, error) {
	var session Session
	err := db.QueryRow(`
        SELECT id, user_id, expires_at
        FROM sessions
        WHERE id = ? AND expires_at > ?
    `, sessionToken, time.Now()).Scan(&session.ID, &session.UserID, &session.ExpiresAt)

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteSession(db *sql.DB, sessionToken string) error {
	if sessionToken == "" {
		return errors.New("session token cannot be empty")
	}

	// Update column name from 'token' to 'id'
	stmt, err := db.Prepare("DELETE FROM sessions WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(sessionToken)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}
