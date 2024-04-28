package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// NamedSelectContext is used to fetch multiple rows from the database
func (d *Database) NamedSelectContext(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := sqlx.BindNamed(sqlx.BindType(d.Sqlx.DriverName()), query, arg)
	if err != nil {
		return err
	}

	return d.Sqlx.SelectContext(ctx, dest, q, args...)
}

// NamedGetContext is used to fetch a single row from the database
func (d *Database) NamedGetContext(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := sqlx.BindNamed(sqlx.BindType(d.Sqlx.DriverName()), query, arg)
	if err != nil {
		return err
	}

	return d.Sqlx.GetContext(ctx, dest, q, args...)
}

// NamedExecContext is used to execute a sql query on the database
func (d *Database) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	q, args, err := sqlx.BindNamed(sqlx.BindType(d.Sqlx.DriverName()), query, arg)
	if err != nil {
		return nil, err
	}

	return d.Sqlx.ExecContext(ctx, q, args...)
}

// NamedExecContextReturnID is used to execute a sql query on the database and return the ID
func (d *Database) NamedExecContextReturnID(ctx context.Context, query string, arg interface{}, ID interface{}) error {
	q, args, err := sqlx.BindNamed(sqlx.BindType(d.Sqlx.DriverName()), query, arg)
	if err != nil {
		return err
	}

	return d.Sqlx.QueryRowx(q, args...).Scan(ID)
}

// NamedExecContextReturnRow is used to execute a sql query on the database and return the row
func (d *Database) NamedExecContextReturnRow(ctx context.Context, query string, arg interface{}, obj interface{}) error {
	q, args, err := sqlx.BindNamed(sqlx.BindType(d.Sqlx.DriverName()), query, arg)
	if err != nil {
		return err
	}

	log.Info("q", zap.Any("q", q), zap.Any("args", args))

	return d.Sqlx.QueryRowx(q, args...).StructScan(obj)
}
