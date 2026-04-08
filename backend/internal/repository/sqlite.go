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
		&models.Project{},
		&models.Task{},
		&models.Comment{},
	)
}

// SeedData 填充种子数据
func SeedData() error {
	// 填充用户数据
	for _, user := range models.SeedUsers {
		var existing models.User
		result := DB.First(&existing, user.ID)
		if result.RowsAffected == 0 {
			if err := DB.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	// 填充项目数据
	for _, project := range models.SeedProjects {
		var existing models.Project
		result := DB.First(&existing, project.ID)
		if result.RowsAffected == 0 {
			if err := DB.Create(&project).Error; err != nil {
				return err
			}
		}
	}

	log.Println("种子数据填充成功")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
