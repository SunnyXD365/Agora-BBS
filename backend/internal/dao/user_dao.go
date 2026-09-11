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
	query := `
		SELECT u.id, u.username, u.password_hash, u.email, u.avatar, u.role, u.status,
		       tp.trust_score, tp.unlock_level, tp.verified_read_seconds,
		       p.onboarding_statement, p.background_tag, p.onboarding_status,
		       u.created_at, u.updated_at
		FROM users u
		JOIN user_trust_profiles tp ON tp.user_id = u.id
		JOIN user_profiles p ON p.user_id = u.id
		WHERE u.id = $1`
	u := &model.User{}
	err := d.db.QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Avatar, &u.Role, &u.Status,
			&u.TrustScore, &u.UnlockLevel, &u.VerifiedReadSeconds, &u.OnboardingStatement,
			&u.BackgroundTag, &u.OnboardingStatus, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (d *UserDAO) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.password_hash, u.email, u.avatar, u.role, u.status,
		       tp.trust_score, tp.unlock_level, tp.verified_read_seconds,
		       p.onboarding_statement, p.background_tag, p.onboarding_status,
		       u.created_at, u.updated_at
		FROM users u
		JOIN user_trust_profiles tp ON tp.user_id = u.id
		JOIN user_profiles p ON p.user_id = u.id
		WHERE u.username = $1`
	u := &model.User{}
	err := d.db.QueryRowContext(ctx, query, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Avatar, &u.Role, &u.Status,
			&u.TrustScore, &u.UnlockLevel, &u.VerifiedReadSeconds, &u.OnboardingStatement,
			&u.BackgroundTag, &u.OnboardingStatus, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}
