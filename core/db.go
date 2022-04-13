package core

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"

	// nolint
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	// nolint
	_ "github.com/lib/pq"
)

// Supported creators of database types.
const (
	MySQLType    = "mysql"
	PostgresType = "postgres"
)

// DB represents Database layer of the service.
type DB struct{}

// NewDB crates new Database instance.
func NewDB() *DB {
	return &DB{}
}

// GetConnection returns connection based on provided database type.
func (d DB) GetConnection(user, password, dbType, dbName, dbHost string) (*sqlx.DB, error) {

	var err error
	var connection *sqlx.DB
	switch dbType {
	case MySQLType:
		connection, err = d.getMySQLConnection(user, password, dbName, dbHost)
	case PostgresType:
		connection, err = d.getPostgresConnection(user, password, dbName, dbHost)
	default:
		return nil, errors.Errorf("unsupported database type. should be %s or %s", MySQLType, PostgresType)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "impossible to get %s database connection", dbType)
	}

	if err := connection.Ping(); err != nil {
		return nil, errors.Wrap(err, "impossible to reach database")
	}

	return connection, nil
}

func (d DB) getMySQLConnection(user, password, dbName, dbHost string) (*sqlx.DB, error) {
	return sqlx.Connect(MySQLType, fmt.Sprintf(
		`%s:%s@(%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci`,
		user,
		password,
		dbHost,
		dbName,
	))
}

func (d DB) getPostgresConnection(user, password, dbName, dbHost string) (*sqlx.DB, error) {
	return sqlx.Connect(PostgresType, fmt.Sprintf(
		`postgres://%s:%s@%s/%s?sslmode=disable`,
		user,
		url.QueryEscape(password),
		dbHost,
		dbName,
	))
}
