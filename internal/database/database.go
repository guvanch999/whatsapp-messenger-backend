package database

import (
	"context"
	. "github.com/medium-messenger/messenger-backend/internal/modules/api-keys/models"
	. "github.com/medium-messenger/messenger-backend/internal/modules/contact-list/model"
	. "github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	. "github.com/medium-messenger/messenger-backend/internal/modules/organization/models"
	. "github.com/medium-messenger/messenger-backend/internal/modules/templates/models"
	. "github.com/medium-messenger/messenger-backend/internal/modules/user-providers/model"
	. "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Connect(dsn string, disableAutoMigration bool) (*gorm.DB, error) {
	database, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Warn),
			NowFunc: func() time.Time {
				utc, _ := time.LoadLocation("")
				return time.Now().In(utc)
			},
			PrepareStmt:     false,
			CreateBatchSize: 100,
			Dialector:       postgres.New(postgres.Config{DSN: dsn}),
		},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(40)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)

	sqlDB.Conn(context.Background())

	if !disableAutoMigration {
		database.AutoMigrate(
			&UserInfo{},
			&UserContact{},
			&ContactList{},
			&Template{},
			&Organization{},
			&UserProvider{},
			&ApiKey{},
		)
	}

	log.Println("Connected to PostgresSql")
	return database, nil
}
