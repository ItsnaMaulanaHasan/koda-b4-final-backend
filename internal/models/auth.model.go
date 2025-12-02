package models

import (
	"backend-koda-shortlink/internal/database"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Register struct {
	Id       int    `json:"id"`
	FullName string `form:"fullname" json:"fullName"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"-"`
}

type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type QueryLogin struct {
	Id       int    `db:"id"`
	Password string `db:"password"`
}

func RegisterUser(bodyRegister *Register) (bool, string, error) {
	isSuccess := false
	message := ""

	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		message = "Failed to start database transaction"
		return isSuccess, message, err
	}
	defer tx.Rollback(ctx)

	// insert data to users
	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (fullname, email, password)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		bodyRegister.FullName, bodyRegister.Email, bodyRegister.Password,
	).Scan(&bodyRegister.Id)
	if err != nil {
		message = "Internal server error while inserting new user"
		return isSuccess, message, err
	}

	// update created_by and updated_by
	_, err = tx.Exec(ctx, `UPDATE users SET created_by = $1, updated_by = $1 WHERE id = $1`, bodyRegister.Id)
	if err != nil {
		message = "Internal server error while update created_by and updated_by"
		return isSuccess, message, err
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		message = "Failed to commit transaction"
		return isSuccess, message, err
	}

	isSuccess = true
	message = "User registered successfully"
	return isSuccess, message, nil
}

func CheckUserEmail(email string) (bool, error) {
	exists := false
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email,
	).Scan(&exists)

	if err != nil {
		return exists, err
	}

	return exists, nil
}

func GetUserByEmail(bodyLogin *Login) (QueryLogin, string, error) {
	message := ""
	user := QueryLogin{}
	rows, err := database.DB.Query(context.Background(),
		"SELECT id, password FROM users WHERE email = $1",
		bodyLogin.Email,
	)
	if err != nil {
		message = "Failed to fetch user from database"
		return user, message, err
	}

	user, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[QueryLogin])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, message, err
		}
		message = "Failed to process user data"
		return user, message, err
	}

	return user, message, nil
}
