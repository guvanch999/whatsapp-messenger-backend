package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"path/filepath"
)

type Schema struct {
	Environment              string `env:"GO_ENV"`
	Port                     string `env:"PORT"`
	AppUrl                   string `env:"APP_URL"`
	SupabaseUrl              string `env:"SUPABASE_URL"`
	SupabasApiKey            string `env:"SUPABASE_KEY"`
	PostgresUri              string `env:"POSTGRES_URI"`
	BranchName               string `env:"BRANCH_NAME"`
	SecretManagerCredentials string `env:"GOOGLE_CREDENTIALS"`
	DisableAutoMigration     bool   `env:"DISABLE_AUTO_MIGRATION" envDefault:"false"`
	SecretKeyForHash         string `env:"SECRET_KEY_FOR_HASH"`
}

var cfg Schema

func GetConfig() *Schema {
	_ = godotenv.Load(filepath.Join(".env"))

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}

	return &cfg
}
