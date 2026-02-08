package repository

import (
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	db *sql.DB
}
type User struct {
    ID        int
    Name      string
    Email     string
    Password  string
    CreatedAt time.Time
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(
	ctx context.Context,
	name, email, password string,
) error {

	query := `
		INSERT INTO users (name, email, password, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		name,
		email,
		password,
		time.Now(),
	)

	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    query := `
        SELECT id, name, email, password, created_at 
        FROM users 
        WHERE email = $1
    `

    var user User
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.CreatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &user, nil
}
func (r *UserRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)`

	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}
