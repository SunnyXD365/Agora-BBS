package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agora-backend/internal/model"
)

type UserDAO struct {
	db *sql.DB
}

func (d *UserDAO) CreateAdminLoginChallenge(ctx context.Context, id string, userID int64, codeHash string, expiresAt time.Time) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `DELETE FROM admin_login_challenges WHERE user_id=$1 AND consumed_at IS NULL`, userID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO admin_login_challenges(id,user_id,code_hash,expires_at) VALUES($1,$2,$3,$4)`, id, userID, codeHash, expiresAt); err != nil {
		return err
	}
	return tx.Commit()
}

func (d *UserDAO) DeleteAdminLoginChallenge(ctx context.Context, id string) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM admin_login_challenges WHERE id=$1`, id)
	return err
}

func (d *UserDAO) ConsumeAdminLoginChallenge(ctx context.Context, id, codeHash string) (int64, bool, error) {
	var userID int64
	var valid bool
	err := d.db.QueryRowContext(ctx, `
		UPDATE admin_login_challenges
		SET attempts=attempts+1, consumed_at=CASE WHEN code_hash=$2 THEN CURRENT_TIMESTAMP ELSE consumed_at END
		WHERE id=$1 AND consumed_at IS NULL AND expires_at>CURRENT_TIMESTAMP AND attempts<5
		RETURNING user_id,code_hash=$2`, id, codeHash).Scan(&userID, &valid)
	return userID, valid, err
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
