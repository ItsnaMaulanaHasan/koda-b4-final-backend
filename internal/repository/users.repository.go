package repository

import (
	"backend-koda-shortlink/internal/database"
	"backend-koda-shortlink/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (fullname, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := database.DB.QueryRow(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.Password,
	).Scan(&user.Id)

	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, profile_photo, fullname, email, password FROM users WHERE email = $1`

	rows, err := database.DB.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetById(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, profile_photo, fullname, email
		FROM users
		WHERE id = $1
	`

	rows, err := database.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := database.DB.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UpdateCreatedByAndUpdatedBy(ctx context.Context, userId int) error {
	query := `UPDATE users SET created_by = $1, updated_by = $1 WHERE id = $1`
	_, err := database.DB.Exec(ctx, query, userId)
	return err
}
