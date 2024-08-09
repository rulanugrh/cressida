package config

import (
	"fmt"
	"time"

	"log"

	"github.com/rulanugrh/cressida/internal/entity/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

type SDatabase struct {
	DB *gorm.DB
}

func InitPostgreSQL() *SDatabase {
	return &SDatabase{}
}

func (conn *SDatabase) DatabaseConnection() *gorm.DB {
	cfg := GetConfig()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
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
		RefreshInterval: 5,
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


func Migration(db *gorm.DB) {
	// migration all struct
	err := db.AutoMigrate(&domain.Role{}, &domain.User{}, &domain.Vehicle{}, &domain.Transporter{}, &domain.Driver{}, &domain.Order{}, &domain.Transaction{})
	fmt.Println(err)
}

func Seeder(db *gorm.DB) error {
	// seeder for role

	cfg := GetConfig()

	roleAdmin := domain.Role{ Name: "Admin", Description: "This is role administrator"}
	roleDrive := domain.Role{ Name: "Driver", Description: "This is role driver"}
	roleUser := domain.Role{ Name: "User", Description: "This is role user"}

	findRoleAdmin := db.Where("name = ?", roleAdmin.Name)
	if findRoleAdmin.RowsAffected < 1 {
		db.Create(&roleAdmin)
	}

	findRoleDriver := db.Where("name = ?", roleDrive.Name)
	if findRoleDriver.RowsAffected < 1 {
		db.Create(&roleDrive)
	}

	findRoleUser := db.Where("name = ?", roleUser.Name)
	if findRoleUser.RowsAffected < 1 {
		db.Create(&roleUser)
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

	if findUserAdmin := db.Where("email = ?", cfg.Admin.Email); findUserAdmin.RowsAffected < 1 {
		db.Create(&userAdmin)
	}
	return nil
}
