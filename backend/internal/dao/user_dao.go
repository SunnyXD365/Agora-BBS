package dao

import (
	"context"
	"database/sql"
	"errors"

	"agora-backend/internal/model"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) CreateUser(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (username, password_hash, email, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING id, role, status, created_at, updated_at
	`
	return d.db.QueryRowContext(ctx, query, u.Username, u.PasswordHash, u.Email, u.Avatar).
		Scan(&u.ID, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt)
}

func (d *UserDAO) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, username, password_hash, email, avatar, role, status, created_at, updated_at FROM users WHERE id = $1`
	u := &model.User{}
	err := d.db.QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Avatar, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (d *UserDAO) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, password_hash, email, avatar, role, status, created_at, updated_at FROM users WHERE username = $1`
	u := &model.User{}
	err := d.db.QueryRowContext(ctx, query, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Avatar, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}
