package db

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

type CVitalDB interface {
	CVProfileDB
	UserDB
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
	SslMode  string
}

type PostgresDB struct {
	sqlxDB *sqlx.DB
	logger zerolog.Logger
}

func NewConnection(config DatabaseConfig, logger zerolog.Logger) (CVitalDB, error) {
	//TODO This doesn't throw an error if the DB connection isn't available?
	var connectionString = "host=" + config.Host +
		" port=" + config.Port +
		" user=" + config.User +
		" dbname=" + config.DbName +
		" password=" + config.Password +
		" sslmode=" + config.SslMode

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	postgresDB := &PostgresDB{
		sqlxDB: db,
		logger: logger,
	}

	err = runMigrations(postgresDB)
	if err != nil {
		return nil, fmt.Errorf("DB migrations failed: %v", err)
	}

	return postgresDB, nil
}

func runMigrations(db *PostgresDB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.sqlxDB.DB, "migrations"); err != nil {
		return err
	}
	return nil
}
