package config

import (
	"context"
	"fmt"

	"log"

	"github.com/rulanugrh/cressida/internal/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SDatabase struct {
	DB *gorm.DB
	trace helper.IOpenTelemetry
}

func InitPostgreSQL() *SDatabase {
	return &SDatabase{trace: helper.NewOpenTelemetry()}
}

func (conn *SDatabase) DatabaseConnection() *gorm.DB {
	// span for detect connection into db
	span := conn.trace.StartTracer(context.Background(), "DatabaseConnection")
	defer span.End()

	cfg := GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error while connect to DB %s", err.Error())
	}

	conn.DB = db
	return db
}
