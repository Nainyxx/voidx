package repository

import (
	"context"
	"database/sql"
	"fmt"
	"uuid"

	"example.com/m/models"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) UserRepository {
	return &postgresRepo{db: db}
}

func ConnectDB(host, port, user, password, dbname string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// GET
func (repo *postgresRepo) GetUserProfileByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	query := `
        SELECT user_id, username, name, surname, phone, description, avatar_image_url, created_at, updated_at
        FROM user_profiles
        WHERE user_id = $1
    `

	profile := &models.UserProfile{}
	err := repo.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.Username,
		&profile.Name,
		&profile.Surname,
		&profile.Phone,
		&profile.Description,
		&profile.AvatarImageURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile not found")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return profile, nil
}

// POST
func (repo *postgresRepo) CreateUserProfile(ctx context.Context, profile *models.UserProfile) error {
	query := `
        INSERT INTO user_profiles (user_id, username, name, surname, phone, description, avatar_image_url, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := repo.db.ExecContext(ctx, query,
		profile.UserID,
		profile.Username,
		profile.Name,
		profile.Surname,
		profile.Phone,
		profile.Description,
		profile.AvatarImageURL,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

// PATCH
func (repo *postgresRepo) UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error {
	query := `
        UPDATE user_profiles
        SET username = $1, name = $2, surname = $3, phone = $4,
            description = $5, avatar_image_url = $6, updated_at = $7
        WHERE user_id = $8
    `

	result, err := repo.db.ExecContext(ctx, query,
		profile.Username, profile.Name, profile.Surname, profile.Phone,
		profile.Description, profile.AvatarImageURL, profile.UpdatedAt,
		profile.UserID,
	)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profile not found")
	}

	return nil
}

// DELETE
func (repo *postgresRepo) DeleteUserProfile(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM user_profiles WHERE user_id = $1
    `

	result, err := repo.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profile not found")
	}

	return nil
}
