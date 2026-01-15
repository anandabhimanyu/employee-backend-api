package auth

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(email, passwordHash string) error {
	_, err := r.db.Exec(
		`INSERT INTO users (email, password_hash) VALUES ($1, $2)`,
		email,
		passwordHash,
	)
	return err
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	var u User

	err := r.db.QueryRow(
		`SELECT id, email, password_hash, created_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return &u, err
}
