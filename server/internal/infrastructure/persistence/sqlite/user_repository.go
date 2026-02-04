package sqlite

import (
	"context"
	"database/sql"
	"time"

	"autostrike/internal/domain/entity"
)

// UserRepository implements repository.UserRepository using SQLite
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new SQLite user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	// Default to active
	if !user.IsActive {
		user.IsActive = true
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, username, email, password_hash, role, is_active, last_login_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.ID, user.Username, user.Email, user.PasswordHash, user.Role, user.IsActive, user.LastLoginAt, user.CreatedAt, user.UpdatedAt)

	return err
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET username = ?, email = ?, password_hash = ?, role = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`, user.Username, user.Email, user.PasswordHash, user.Role, user.IsActive, user.UpdatedAt, user.ID)

	return err
}

// Delete deletes a user (hard delete - use Deactivate for soft delete)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	return err
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &lastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &lastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE email = ?
	`, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &lastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindAll finds all users (including inactive)
func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, username, email, password_hash, role, is_active, last_login_at, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// FindActive finds all active users
func (r *UserRepository) FindActive(ctx context.Context) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, username, email, password_hash, role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE is_active = 1 ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET last_login_at = ?, updated_at = ? WHERE id = ?
	`, now, now, id)
	return err
}

// Deactivate soft-deletes a user by setting is_active to false
func (r *UserRepository) Deactivate(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET is_active = 0, updated_at = ? WHERE id = ?
	`, now, id)
	return err
}

// Reactivate re-enables a deactivated user
func (r *UserRepository) Reactivate(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET is_active = 1, updated_at = ? WHERE id = ?
	`, now, id)
	return err
}

// CountByRole counts active users with a specific role
func (r *UserRepository) CountByRole(ctx context.Context, role entity.UserRole) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users WHERE role = ? AND is_active = 1
	`, role).Scan(&count)
	return count, err
}

func (r *UserRepository) scanUsers(rows *sql.Rows) ([]*entity.User, error) {
	var users []*entity.User

	for rows.Next() {
		user := &entity.User{}
		var lastLoginAt sql.NullTime

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &lastLoginAt, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
