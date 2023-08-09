package config

import (
	"database/sql"
	"fmt"

	"github.com/wlrudi19/elastic-engine/config"
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
func (*database) LoadConfig() config.Config {
	return Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "projectrudi",
			Username: "rudilesmana",
			Password: "rudilesmana2023",
		},
	}
}

// TestConnection implements Database.
func (*database) TestConnection(config config.DatabaseConfig) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password,
	)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return db, err
}
