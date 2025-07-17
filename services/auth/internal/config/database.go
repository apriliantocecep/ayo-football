package config

import (
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/entity"
	"github.com/apriliantocecep/posfin-blog/shared"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewDatabase(vaultClient *shared.VaultClient) *gorm.DB {
	var gormDB *gorm.DB

	secret := utils.GetVaultSecretAuthSvc(vaultClient)

	dsn := secret["DATABASE_URL"].(string)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(100)
	conn.SetConnMaxLifetime(time.Duration(300) * time.Second)

	// Auto Migrate
	err = db.AutoMigrate(
		&entity.User{},
	)
	if err != nil {
		log.Fatalf("failed to migrate the entities: %v", err)
	}

	gormDB = db

	return gormDB
}
