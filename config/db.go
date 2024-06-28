package config

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"log"
)

type SDatabase struct {
	DB *gorm.DB
}

func InitPostgreSQL() *SDatabase {
	return &SDatabase{}
}

func (conn *SDatabase) DatabaseConnection() *gorm.DB {
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
