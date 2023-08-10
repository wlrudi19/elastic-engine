package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database interface {
	LoadConfig() Config
	TestConnection(config DatabaseConfig) (*sql.DB, error)
}

type database struct {
}

func NewLoad() Database {
	return &database{}
}

// LoadConfig implements Database.
func (d *database) LoadConfig() Config {
	return Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "",
			Username: "",
			Password: "",
		},
	}
}

// TestConnection implements Database.
func (d *database) TestConnection(config DatabaseConfig) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Host, config.Port, config.Name, config.Username, config.Password,
	)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return db, err
}
