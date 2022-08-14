package db

import (
	"embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

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

func NewConnection(config DatabaseConfig) (*sqlx.DB, error) {
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

	return db, nil
}

func RunMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}
	return nil
}
