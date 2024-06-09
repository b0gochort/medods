package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func GetEmail(guid string, db *sqlx.DB) (string, error) {
	var email string

	query := `SELECT email FROM users WHERE id = $1;`
	if err := db.QueryRow(query, guid).Scan(&email); err != nil {
		return "", fmt.Errorf("db.QueryRow: %w", err)
	}

	return email, nil
}
