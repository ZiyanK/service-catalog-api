package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	queryInsertService = `INSERT INTO services(name, description, user_uuid) VALUES(:name, :description, :user_uuid)`

	queryCheckServiceByNameAndUserUUID = `
	SELECT COUNT(1) FROM services s
	WHERE s.name = :name AND s.user_uuid = :user_uuid`

	queryCheckServiceByIDAndUserUUID = `
	SELECT COUNT(1) FROM services s
	WHERE s.service_id = :service_id AND s.user_uuid = :user_uuid`

	queryGetService = `
	SELECT s.service_id, s.name, s.description, COALESCE(sv.sv_id, 0) as sv_id, COALESCE(sv.version,'') as version, COALESCE(sv.changelog,'') as changelog
	FROM services s
	LEFT JOIN service_versions sv ON sv.service_id = s.service_id
	WHERE s.user_uuid = :user_uuid AND s.service_id = :service_id`

	queryUpdateService = `
	UPDATE services s SET name = :name, description = :description, updated_at = NOW()
	WHERE s.service_id = :service_id
	RETURNING *`

	queryDeleteService         = `DELETE FROM services WHERE service_id = :service_id`
	queryDeleteServiceVersions = `DELETE FROM service_versions WHERE service_id = :service_id`
)

// Service is a struct used to represent the `services` table in the database
type Service struct {
	ServiceID     int       `db:"service_id" json:"service_id"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	UserUUID      uuid.UUID `db:"user_uuid" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	VersionsCount int       `db:"versions_count" json:"versions_count"`
}

// ServiceWithVersions is a struct used to get the given service and all of it's version
type ServiceWithVersions struct {
	ServiceID   int    `db:"service_id" json:"service_id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	SvID        int    `db:"sv_id" json:"sv_id"`
	Version     string `db:"version" json:"version"`
	Changelog   string `db:"changelog" json:"changelog"`
}

// CreateService is used to create a new service for a user
func (service *Service) CreateService(ctx context.Context) error {
	tx, err := db.Sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	q, args, err := sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryCheckServiceByNameAndUserUUID, map[string]interface{}{
		"name":      service.Name,
		"user_uuid": service.UserUUID,
	})
	if err != nil {
		log.Error("error building service fetch query", zap.Error(err))
		tx.Rollback()
		return err
	}

	var count int

	err = tx.GetContext(ctx, &count, q, args...)
	if err != nil && err != sql.ErrNoRows {
		log.Error("error querying user", zap.Error(err))
		return err
	}

	if count > 0 {
		tx.Rollback()
		log.Info("service with same name exists")
		return errors.New("service exists")
	}

	// If service doesn't exist
	q, args, err = sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryInsertService, map[string]interface{}{
		"name":        service.Name,
		"description": service.Description,
		"user_uuid":   service.UserUUID,
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

	if rowsAffected == 1 {
		tx.Commit()
		return nil
	}

	tx.Rollback()
	log.Info("no row was updated")
	return err
}

// GetServices is used to fetch all the services present ofr a given user
func GetServices(ctx context.Context, userUUID uuid.UUID, limit, offset int, serviceName, orderBy string) ([]Service, error) {
	var services []Service

	// This query was required to be initialized here to add the order by (ASC,DESC) clause
	// Reason: sqlx does not permit to pass keywords as args
	querySelectServices := `
	SELECT s.service_id, s.name, s.description, s.created_at, s.updated_at, COUNT(sv.service_id) AS versions_count
	FROM services s
	LEFT JOIN service_versions sv ON s.service_id = sv.service_id
	WHERE
		s.user_uuid = :user_uuid AND s.name LIKE :name
	GROUP BY s.service_id
	ORDER BY s.created_at %s
	LIMIT :limit OFFSET :offset`

	params := map[string]interface{}{
		"user_uuid": userUUID,
		"name":      "%%",
		"limit":     10,
		"offset":    0,
	}

	if limit > 0 {
		params["limit"] = limit
	}
	if offset > 0 {
		params["offset"] = offset
	}
	if len(serviceName) > 0 {
		params["name"] = fmt.Sprintf("%%%v%%", serviceName)
	}

	// To add order by
	if orderBy == "DESC" {
		querySelectServices = fmt.Sprintf(querySelectServices, "DESC")
	} else {
		querySelectServices = fmt.Sprintf(querySelectServices, "ASC")
	}

	err := db.NamedSelectContext(ctx, &services, querySelectServices, params)
	if err != nil {
		log.Error("Error while fetching services", zap.Error(err))
		return nil, err
	}

	return services, nil
}

// GetService is used to get a paritcular service with all it's versions
func GetService(ctx context.Context, serviceID int, userUUID uuid.UUID) ([]ServiceWithVersions, error) {
	var service []ServiceWithVersions

	err := db.NamedSelectContext(ctx, &service, queryGetService, map[string]interface{}{
		"user_uuid":  userUUID,
		"service_id": serviceID,
	})
	if err != nil {
		log.Error("Error while fetching service", zap.Error(err))
		return nil, err
	}

	return service, nil
}

// UpdateService is used to update the service name and description of a given service
func (service *Service) UpdateService(ctx context.Context) error {
	tx, err := db.Sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	q, args, err := sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryCheckServiceByIDAndUserUUID, map[string]interface{}{
		"service_id": service.ServiceID,
		"user_uuid":  service.UserUUID,
	})
	if err != nil {
		log.Error("error building service fetch query", zap.Error(err))
		tx.Rollback()
		return err
	}

	var count int

	err = tx.GetContext(ctx, &count, q, args...)
	if err != nil && err != sql.ErrNoRows {
		log.Error("error querying user", zap.Error(err))
		return err
	}

	// If service was not found
	if count == 0 {
		tx.Rollback()
		log.Info("service does not exist")
		return errors.New("service does not exist")
	}

	// If service is found
	result, err := db.NamedExecContext(ctx, queryUpdateService, map[string]interface{}{
		"service_id":  service.ServiceID,
		"name":        service.Name,
		"description": service.Description,
	})
	if err != nil {
		tx.Rollback()
		log.Error("Error while updating service", zap.Error(err))
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
		return errors.New("no rows updated")
	}

	tx.Commit()
	return nil
}

// DeleteService is used to delete a given service and all it's version
func (service *Service) DeleteService(ctx context.Context) error {
	tx, err := db.Sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	q, args, err := sqlx.BindNamed(sqlx.BindType(db.Sqlx.DriverName()), queryCheckServiceByIDAndUserUUID, map[string]interface{}{
		"service_id": service.ServiceID,
		"user_uuid":  service.UserUUID,
	})
	if err != nil {
		tx.Rollback()
		log.Error("error building service fetch query", zap.Error(err))
		return err
	}

	var count int

	err = tx.GetContext(ctx, &count, q, args...)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		log.Error("error querying user", zap.Error(err))
		return err
	}

	// If service was not found
	if count == 0 {
		tx.Rollback()
		log.Info("service does not exist")
		return errors.New("service does not exist")
	}

	// If service is found

	// Deleting service versions
	// Deleting versions first due to foreign key constraint
	_, err = db.NamedExecContext(ctx, queryDeleteServiceVersions, map[string]interface{}{
		"service_id": service.ServiceID,
	})
	if err != nil {
		tx.Rollback()
		log.Error("Error while deleting service", zap.Error(err))
		return err
	}

	// Deleting service
	result, err := db.NamedExecContext(ctx, queryDeleteService, map[string]interface{}{
		"service_id": service.ServiceID,
	})
	if err != nil {
		tx.Rollback()
		log.Error("Error while deleting service", zap.Error(err))
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
		return errors.New("no rows deleted")
	}

	tx.Commit()
	return nil
}
