package auth

import "database/sql"

type Repository interface {
	Create(*User) error
	GetByEmail(string) (*User, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(u *User) error {
	return r.db.QueryRow(`
		INSERT INTO users (email, password_hash, role)
		VALUES ($1,$2,$3)
		RETURNING id, created_at
	`, u.Email, u.PasswordHash, u.Role).Scan(&u.ID, &u.CreatedAt)
}

func (r *postgresRepository) GetByEmail(email string) (*User, error) {
	var u User
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, role, created_at
		FROM users WHERE email=$1
	`, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	return &u, err
}
