package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort int
	MCI        float64
}

func Load() (*Config, error) {
	// Parse database port
	dbPortStr := getEnv("DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil || dbPort < 1 || dbPort > 65535 {
		return nil, fmt.Errorf("invalid DB_PORT: must be 1-65535, got %s", dbPortStr)
	}

	// Parse server port
	serverPortStr := getEnv("SERVER_PORT", "8080")
	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil || serverPort < 1 || serverPort > 65535 {
		return nil, fmt.Errorf("invalid SERVER_PORT: must be 1-65535, got %s", serverPortStr)
	}

	// Parse MCI
	mciStr := os.Getenv("MCI")
	mci, err := strconv.ParseFloat(mciStr, 64)
	if err != nil || mci <= 0 {
		mci = 3932.0
	}

	// SSL mode: require by default (more secure), or disable explicitly if needed
	sslMode := getEnv("DB_SSL_MODE", "require")
	if sslMode != "require" && sslMode != "disable" && sslMode != "prefer" {
		return nil, fmt.Errorf("invalid DB_SSL_MODE: must be require|disable|prefer, got %s", sslMode)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "salary_db"),
		DBSSLMode:  sslMode,
		ServerPort: serverPort,
		MCI:        mci,
	}, nil
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
