package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	MCI        float64
}

func Load() *Config {
	mciStr := os.Getenv("MCI")
	mci, err := strconv.ParseFloat(mciStr, 64)
	if err != nil || mci <= 0 {
		mci = 3932
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "salary_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		MCI:        mci,
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
