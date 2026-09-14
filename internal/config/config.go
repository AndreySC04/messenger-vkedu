package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PoolConfig struct {
	MaxConns        string
	MinConns        string
	MaxConnLifetime string
	MaxConnIdleTime string
}

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	Pool       PoolConfig
}

type Config struct {
	AppPort string
	DB      DBConfig
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists && value != "" {
		return value
	}
	return defaultValue
}

func ConfigLoad() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("LOG - .env FILE NOT FOUND")
	}
	config := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DB: DBConfig{
			DBHost:     getEnv("POSTGRES_HOST", "127.0.0.1"),
			DBPort:     getEnv("POSTGRES_PORT", "5434"),
			DBUser:     getEnv("POSTGRES_USER", "messenger_user"),
			DBPassword: getEnv("POSTGRES_PASSWORD", "messenger_password"),
			DBName:     getEnv("POSTGRES_DB", "messenger_db"),
			SSLMode:    getEnv("POSTGRES_SSLMODE", "disable"),
			Pool: PoolConfig{
				MaxConns:        getEnv("POSTGRES_POOL_MAX_CONNS", "25"),
				MinConns:        getEnv("POSTGRES_POOL_MIN_CONNS", "5"),
				MaxConnLifetime: getEnv("POSTGRES_POOL_MAX_CONN_LIFETIME", "1h"),
				MaxConnIdleTime: getEnv("POSTGRES_POOL_MAX_CONN_IDLE_TIME", "30m"),
			},
		},
	}
	log.Println("LOG - CONFIG LOADED")

	return config, nil
}

func (db DBConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=%s&pool_min_conns=%s&pool_max_conn_lifetime=%s&pool_max_conn_idle_time=%s",
		db.DBUser, db.DBPassword, db.DBHost, db.DBPort, db.DBName, db.SSLMode, db.Pool.MaxConns, db.Pool.MinConns, db.Pool.MaxConnLifetime, db.Pool.MaxConnIdleTime,
	)
}
