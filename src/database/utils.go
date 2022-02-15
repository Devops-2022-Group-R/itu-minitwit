package database

import "database/sql"

// Convenience method to look up the id for a username.
func GetUserId(username string, db *sql.DB) *int64 {
	row := db.QueryRow("select user_id from user where username = ?", username)

	var userId int64
	err := row.Scan(&userId)
	if err != nil {
		return nil
	}

	return &userId
}
