package repository

import (
	"log"

	"github.com/taskify/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitSQLite 初始化SQLite数据库
func InitSQLite(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移
	err = AutoMigrate()
	if err != nil {
		return err
	}

	// 填充种子数据
	err = SeedData()
	if err != nil {
		return err
	}

	log.Println("SQLite数据库初始化成功")
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.ProjectMembership{},
		&models.AuditLog{},
		&models.Project{},
		&models.Task{},
		&models.Comment{},
		&models.PersonalAccessToken{},
	)
}

// SeedData 填充种子数据（由seed.go提供）

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
