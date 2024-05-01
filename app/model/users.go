package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	queryInsertUser = `INSERT INTO users(user_uuid, email, password) VALUES(:user_uuid, :email, :password) RETURNING *`

	queryCheckUserExist = `SELECT count(1) FROM users WHERE email = :email`

	queryGetUserByID = `
	SELECT u.user_uuid, u.email
	FROM users u
	WHERE u.user_uuid = :user_uuid`

	queryGetUserByEmail = `
	SELECT u.user_uuid, u.email, u.password
	FROM users u
	WHERE email = :email`

	queryUpdateUserByID = `
	UPDATE users SET email = :email, updated_at = NOW()
	WHERE user_uuid = :user_uuid`
)

// User is a struct used to represent the `users` table in the database
type User struct {
	UserUUID  uuid.UUID `db:"user_uuid" json:"-"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateUser is used to create a new user in the database
func (user *User) CreateUser(ctx context.Context) error {
	tx, err := db.Sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	q, args, err := sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryCheckUserExist, map[string]interface{}{
		"email": user.Email,
	})
	if err != nil {
		log.Error("error building user fetch query", zap.Error(err))
		tx.Rollback()
		return err
	}

	var count int

	err = tx.GetContext(ctx, &count, q, args...)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		log.Error("error querying user", zap.Error(err))
		return err
	}

	if count == 1 {
		tx.Rollback()
		log.Info("user with mail exists")
		return errors.New("mail exists")
	}

	// If mail doesn't exist
	q, args, err = sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryInsertUser, user)
	if err != nil {
		tx.Rollback()
		log.Error("error building user insert query", zap.Error(err))
		return err
	}

	result, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		tx.Rollback()
		log.Error("error inserting user", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Error("Error while getting no. of rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		log.Info("no row was updated")
		return err
	}

	tx.Commit()
	return nil
}

// GetUserByEmail is used to fetch a user using the email
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := db.NamedGetContext(ctx, &user, queryGetUserByEmail, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Error("Error while fetching user by email", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

// GetUserByID is used to fetch a user using the userUUID
func GetUserByID(ctx context.Context, userUUID uuid.UUID) (*User, error) {
	var user User

	err := db.NamedGetContext(ctx, &user, queryGetUserByID, map[string]interface{}{
		"user_uuid": userUUID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("No user found for id", zap.Any("uuid", userUUID))
			return nil, sql.ErrNoRows
		}
		log.Error("Error while fetching user by user_uuid", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

// UpdateUser is used to update the user email
func UpdateUser(ctx context.Context, updatedEmail string, userUUID uuid.UUID) error {
	tx, err := db.Sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	q, args, err := sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryGetUserByEmail, map[string]interface{}{
		"email": updatedEmail,
	})
	if err != nil {
		log.Error("error building user fetch query", zap.Error(err))
		tx.Rollback()
		return err
	}

	var user User

	err = tx.GetContext(ctx, &user, q, args...)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		log.Error("error querying user", zap.Error(err))
		return err
	}

	if user.Email != "" {
		tx.Rollback()
		log.Info("user with mail exists")
		return errors.New("mail exists")
	}

	// If mail doesn't exist
	q, args, err = sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryUpdateUserByID, map[string]interface{}{
		"email":     updatedEmail,
		"user_uuid": userUUID,
	})
	if err != nil {
		tx.Rollback()
		log.Error("error building user update query", zap.Error(err))
		return err
	}

	result, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		tx.Rollback()
		log.Error("error updating user", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Error("Error while getting no. of rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		log.Info("no row was updated")
		return err
	}

	tx.Commit()
	return nil
}
