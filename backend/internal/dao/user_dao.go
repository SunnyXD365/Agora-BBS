package dao

import (
	"context"
	"database/sql"
	"errors"

	"Agora-BBS/internal/model"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

// CreateUser 写入新用户
func (d *UserDAO) CreateUser(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (username, password_hash, email, avatar, role, status, trust_score, unlock_level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return d.db.QueryRowContext(
		ctx, query,
		u.Username, u.PasswordHash, u.Email, u.Avatar, u.Role, u.Status, u.TrustScore, u.UnlockLevel,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

// GetUserByUsername 根据用户名查询用户
func (d *UserDAO) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, email, avatar, role, status, trust_score, unlock_level, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	u := &model.User{}
	err := d.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Avatar,
		&u.Role, &u.Status, &u.TrustScore, &u.UnlockLevel, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetUserByID 根据用户 ID 查询用户
func (d *UserDAO) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, email, avatar, role, status, trust_score, unlock_level, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	u := &model.User{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Avatar,
		&u.Role, &u.Status, &u.TrustScore, &u.UnlockLevel, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
