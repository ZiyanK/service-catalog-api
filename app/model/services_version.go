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
	queryInsertServiceVersion = `INSERT INTO service_versions(version, changelog, service_id) VALUES (:version, :changelog, :service_id)`

	queryCheckServiceVersionUsingVersion = `
	SELECT count(1)
	FROM service_versions sv
	JOIN services s ON s.service_id = sv.service_id
	JOIN users u ON s.user_uuid = u.user_uuid
	WHERE
		sv.version = :version
		AND s.service_id = :service_id
		AND u.user_uuid = :user_uuid`

	queryCheckServiceVersionUsingSVID = `
	SELECT count(1)
	FROM service_versions sv
	JOIN services s ON s.service_id = sv.service_id
	JOIN users u ON s.user_uuid = u.user_uuid
	WHERE
		sv.sv_id = :sv_id
		AND s.service_id = :service_id
		AND u.user_uuid = :user_uuid`

	queryDeleteServiceVersion = `DELETE FROM service_versions WHERE sv_id = :sv_id`
)

// ServiceVersion is a struct used to represent the `service_versions` table in the database
type ServiceVersion struct {
	SvID      int       `db:"sv_id" json:"sv_id"`
	Version   string    `db:"version" json:"version"`
	Changelog string    `db:"changelog" json:"changelog"`
	ServiceID int       `db:"service_id" json:"service_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateServiceVersion is used to create a new service version for a given service
func (sv *ServiceVersion) CreateServiceVersion(ctx context.Context, userUUID uuid.UUID) error {
	tx, err := db.Sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Check if service exists
	q, args, err := sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryCheckServiceByIDAndUserUUID, map[string]interface{}{
		"service_id": sv.ServiceID,
		"user_uuid":  userUUID,
	})
	if err != nil {
		log.Error("error building service fetch query", zap.Error(err))
		tx.Rollback()
		return err
	}

	var serviceCount int

	err = tx.GetContext(ctx, &serviceCount, q, args...)
	if err != nil && err != sql.ErrNoRows {
		log.Error("error querying user", zap.Error(err))
		return err
	}

	// If service was not found
	if serviceCount == 0 {
		tx.Rollback()
		log.Info("service does not exist")
		return errors.New("service does not exist")
	}

	// Check if service version with same version exists
	q, args, err = sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryCheckServiceVersionUsingVersion, map[string]interface{}{
		"version":    sv.Version,
		"service_id": sv.ServiceID,
		"user_uuid":  userUUID,
	})
	if err != nil {
		log.Error("error building service version check query", zap.Error(err))
		tx.Rollback()
		return err
	}

	var svCount int

	err = tx.GetContext(ctx, &svCount, q, args...)
	if err != nil && err != sql.ErrNoRows {
		log.Error("error querying service version", zap.Error(err))
		return err
	}

	// If service version with same version exists, svCount is 1
	if svCount == 1 {
		tx.Rollback()
		log.Info("service with same version exists")
		return errors.New("service version exists")
	}

	// If service version with same version doesn't exist
	q, args, err = sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryInsertServiceVersion, map[string]interface{}{
		"version":    sv.Version,
		"changelog":  sv.Changelog,
		"service_id": sv.ServiceID,
	})
	if err != nil {
		log.Error("error building service insert query", zap.Error(err))
		tx.Rollback()
		return err
	}

	result, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		log.Error("error inserting service", zap.Error(err))
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Error("Error while getting no. of rows affected", zap.Error(err))
		return err
	}

	// Service version created
	if rowsAffected == 0 {
		tx.Rollback()
		log.Info("no row were added")
		return errors.New("error while creating service")
	}

	tx.Commit()
	return nil
}

// DeleteServiceVersion is used to delete a particular service version for a given service
func DeleteServiceVersion(ctx context.Context, userUUID uuid.UUID, serviceID, svID int) error {
	tx, err := db.Sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Check if service_version belongs to user
	q, args, err := sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryCheckServiceVersionUsingSVID, map[string]interface{}{
		"sv_id":      svID,
		"service_id": serviceID,
		"user_uuid":  userUUID,
	})
	if err != nil {
		log.Error("error building service version check query", zap.Error(err))
		tx.Rollback()
		return err
	}

	var count int

	err = tx.GetContext(ctx, &count, q, args...)
	if err != nil && err != sql.ErrNoRows {
		log.Error("error querying service version", zap.Error(err))
		return err
	}

	if count == 0 {
		tx.Rollback()
		log.Info("service version does not exist")
		return errors.New("service version does not exist")
	}

	// If service version is present
	q, args, err = sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryDeleteServiceVersion, map[string]interface{}{
		"sv_id": svID,
	})
	if err != nil {
		log.Error("error building service version delete query", zap.Error(err))
		tx.Rollback()
		return err
	}

	result, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		log.Error("error deleting service", zap.Error(err))
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Error("Error while getting no. of rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 1 {
		tx.Commit()
		return nil
	}

	tx.Rollback()
	log.Info("no row were deleted")
	return errors.New("no rows deleted")
}
