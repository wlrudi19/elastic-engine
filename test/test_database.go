package test

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wlrudi19/elastic-engine/models"
)

func TestConnection(config models.DatabaseConfig) (*sql.DB, error) {
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
