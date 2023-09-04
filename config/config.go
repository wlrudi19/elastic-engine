package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func LoanConfig() Config {
	return Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "projectrudi",
			Username: "rudilesmana",
			Password: "rudilesmana2023",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "rudilesmana2023",
			DBIndex:  1,
		},
	}
}

func ConnectConfig(config DatabaseConfig, redisConfig RedisConfig) (*sql.DB, *redis.Client, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Host, config.Port, config.Name, config.Username, config.Password,
	)

	connDB, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
		return nil, nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       1,
	})

	_, err = redisClient.Ping().Result()

	if err != nil {
		log.Fatalf("error connecting to redis: %v", err)
		return nil, nil, err
	}

	return connDB, redisClient, nil
}
