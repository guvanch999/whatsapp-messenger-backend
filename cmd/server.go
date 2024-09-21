package cmd

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/config"
	"github.com/medium-messenger/messenger-backend/internal/database"
	"github.com/medium-messenger/messenger-backend/internal/validator"
	supa "github.com/nedpals/supabase-go"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"log"
)

type Server struct {
	Echo                *echo.Echo
	Config              *config.Schema
	Database            *gorm.DB
	Supabase            *supa.Client
	SecretManagerClient *secretmanager.Client
}

func NewServer() *Server {
	cfg := config.GetConfig()

	supabase := supa.CreateClient(cfg.SupabaseUrl, cfg.SupabasApiKey)
	db, err := database.Connect(cfg.PostgresUri, cfg.DisableAutoMigration)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Validator = validator.NewValidator()

	ctx := context.Background()
	secretManagerClient, err := secretmanager.NewClient(
		ctx,
		option.WithCredentialsJSON([]byte(cfg.SecretManagerCredentials)),
	)
	if err != nil {
		log.Fatalf("failed to setup secret manager client: %v", err.Error())
	}

	return &Server{
		Echo:                e,
		Config:              cfg,
		Database:            db,
		Supabase:            supabase,
		SecretManagerClient: secretManagerClient,
	}
}
