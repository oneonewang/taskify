package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/config"
	"github.com/taskify/backend/internal/handlers"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/repository"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置Gin模式
	gin.SetMode(cfg.GinMode)

	// 初始化数据库
	dbPath := "taskify.db"
	if err := repository.InitSQLite(dbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建Gin实例
	r := gin.Default()

	// 注册中间件
	r.Use(middleware.CORS())

	// 注册路由
	registerRoutes(r)

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// registerRoutes 注册所有路由
func registerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 用户路由
		api.GET("/users", handlers.GetUsers)

		// 项目路由
		api.GET("/projects", handlers.GetProjects)
		api.GET("/projects/:id", handlers.GetProject)
		api.GET("/projects/:id/tasks", handlers.GetTasks)
		api.POST("/projects/:id/tasks", handlers.CreateTask)

		// 任务路由
		api.PUT("/tasks/:id", handlers.UpdateTask)
		api.PUT("/tasks/:id/status", handlers.UpdateTaskStatus)
		api.DELETE("/tasks/:id", handlers.DeleteTask)

		// 评论路由
		api.GET("/tasks/:id/comments", handlers.GetComments)
		api.POST("/tasks/:id/comments", handlers.CreateComment)
		api.PUT("/tasks/:id/comments/:cid", handlers.UpdateComment)
		api.DELETE("/tasks/:id/comments/:cid", handlers.DeleteComment)

		// SSE事件路由
		api.GET("/events", handlers.GetEvents)
	}
}
