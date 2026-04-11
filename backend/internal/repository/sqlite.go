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
	err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.ProjectMembership{},
		&models.AuditLog{},
		&models.Project{},
		&models.Task{},
		&models.Comment{},
	)
	if err != nil {
		return err
	}

	// 创建复合索引（AutoMigrate 不会自动创建多列索引）
	if err := createCompositeIndexes(); err != nil {
		return err
	}

	return nil
}

// createCompositeIndexes 创建复合索引
func createCompositeIndexes() error {
	// 复合索引：优化 project_memberships JOIN 查询
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_pm_project_user_role
		ON project_memberships(project_id, user_id, role_id)
	`).Error; err != nil {
		return err
	}

	// 复合索引：优化 tasks 表看板查询 (project_id + status)
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tasks_project_status
		ON tasks(project_id, status)
	`).Error; err != nil {
		return err
	}

	// 索引：tasks.assignee_id
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tasks_assignee
		ON tasks(assignee_id)
	`).Error; err != nil {
		return err
	}

	// 索引：comments.task_id
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_comments_task_id
		ON comments(task_id)
	`).Error; err != nil {
		return err
	}

	// 索引：projects.is_archived
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_projects_archived
		ON projects(is_archived)
	`).Error; err != nil {
		return err
	}

	// 唯一约束：project_memberships (user_id, project_id)
	if err := DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_pm_user_project
		ON project_memberships(user_id, project_id)
	`).Error; err != nil {
		return err
	}

	return nil
}

// SeedData 填充种子数据（由seed.go提供）

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
