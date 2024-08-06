package config

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/helper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
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

	// use plugin for prometheus
	db.Use(prometheus.New(prometheus.Config{
		DBName: "cresida_db",
		RefreshInterval: 15,
		HTTPServerPort: 3000,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.Postgres{
				VariableNames: []string{"threads_running"},
				Prefix: "db_stats_",
			},
		},
	}))

	// set max connection pool
	sql, err := db.DB()
	if err != nil {
		log.Printf("Error while set DB %s", err.Error())
	}

	// max open connection
	sql.SetMaxOpenConns(100)
	// max idle connection
	sql.SetMaxIdleConns(2)
	// max lifetime
	sql.SetConnMaxLifetime(time.Since(time.Now().Add(30 * time.Minute)))

	conn.DB = db
	return db
}


func (conn *SDatabase) Migration() {
	// migration all struct
	conn.DB.AutoMigrate(&domain.Role{}, &domain.User{}, &domain.Vehicle{}, &domain.Transporter{}, &domain.Driver{}, &domain.Order{}, &domain.Transaction{})
}

func (conn *SDatabase) Seeder() {
	// seeder for role

	cfg := GetConfig()

	roleAdmin := domain.Role{ Name: "Admin", Description: "This is role administrator"}
	roleDrive := domain.Role{ Name: "Driver", Description: "This is role driver"}
	roleUser := domain.Role{ Name: "User", Description: "This is role user"}

	findRoleAdmin := conn.DB.Where("name = ?", roleAdmin.Name)
	if findRoleAdmin.RowsAffected < 0 {
		findRoleAdmin.Create(&roleAdmin)
	}

	findRoleDriver := conn.DB.Where("name = ?", roleDrive.Name)
	if findRoleDriver.RowsAffected < 0 {
		findRoleDriver.Create(&roleDrive)
	}

	findRoleUser := conn.DB.Where("name = ?", roleUser.Name)
	if findRoleUser.RowsAffected < 0 {
		findRoleUser.Create(&roleUser)
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), 14)
	userAdmin := domain.User{
		FName: "Admin",
		LName: "Cressida",
		Email: cfg.Admin.Email,
		Password: string(hashPassword),
		RoleID: 1,
		Address: "-",
		Phone: "-",
	}

	if findUserAdmin := conn.DB.Where("email = ?", cfg.Admin.Email); findUserAdmin.RowsAffected < 0 {
		findUserAdmin.Create(&userAdmin)
	}
}
