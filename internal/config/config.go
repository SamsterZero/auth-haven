package config

import (
	"os"
)

type Config struct {
	DBUrl     string
	GRPCPort  string
	JWTSecret string
}

func Load() (*Config, error) {
	return &Config{
		DBUrl:     getEnv("DB_URL", "postgres://user:pass@localhost:5432/auth"),
		GRPCPort:  getEnv("GRPC_PORT", ":50051"),
		JWTSecret: getEnv("JWT_SECRET", "supersecret"),
	}, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
