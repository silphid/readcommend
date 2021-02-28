package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Top-level configs, private because the whole struct should
// not be passed down the call chain, only its individual values.
type config struct {
	Env         string
	LogLevel    string        `envconfig:"LOG_LEVEL" default:"info"`
	BaseURL     string        `envconfig:"BASE_URL" default:"http://localhost:5000"`
	Port        uint16        `envconfig:"PORT" default:"5000"`
	DBUrl       string        `envconfig:"DB_URL" default:"postgres://postgres:password123@localhost:5432/readcommend?sslmode=disable"`
	GracePeriod time.Duration `envconfig:"GRACE_PERIOD" default:"10s"`
}

func getConfig() config {
	// Optionally load env vars from a `.env.xxx` file
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	godotenv.Load(".env." + env)
	godotenv.Load()

	// Let Mr. Hightower work his magic!!
	cfg := config{Env: env}
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	return cfg
}
