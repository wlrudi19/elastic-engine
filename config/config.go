package config

import database "github.com/wlrudi19/elastic-engine/models"

func LoadConfig() database.Config {
	return database.Config{
		Database: database.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "projectrudi",
			Username: "rudilesmana",
			Password: "",
		},
	}
}
