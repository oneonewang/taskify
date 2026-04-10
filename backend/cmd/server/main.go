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

	// 初始化会话中间件
	sessionCfg := config.DefaultSessionConfig()
	config.InitSession(r, sessionCfg)

	// 注册中间件
	r.Use(middleware.CORS())

	// 注册路由
	projectHandler := handlers.NewProjectHandler()
	registerRoutes(r, projectHandler)

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// registerRoutes 注册所有路由
func registerRoutes(r *gin.Engine, projectHandler *handlers.ProjectHandler) {
	api := r.Group("/api")
	{
		// 认证路由（无需登录）
		authHandler := handlers.NewAuthHandler()
		userHandler := handlers.NewUserHandler()
		roleHandler := handlers.NewRoleHandler()
		permissionHandler := handlers.NewPermissionHandler()
		membershipHandler := handlers.NewMembershipHandler()
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// 需要登录的路由
		authRequired := api.Group("")
		authRequired.Use(middleware.AuthRequired())
		{
			// 用户路由 - 只有本人可以访问
			authRequired.GET("/users/me", userHandler.GetCurrentUser)
			authRequired.PUT("/users/me", userHandler.UpdateProfile)
			authRequired.PUT("/users/me/password", userHandler.ChangePassword)
		}

		// 管理员路由
		adminRequired := api.Group("/admin")
		adminRequired.Use(middleware.AuthRequired(), middleware.RequireAdmin())
		{
			// 管理员用户管理
			adminRequired.GET("/users", handlers.GetUsers)
			// 角色管理
			adminRequired.GET("/roles", roleHandler.GetRoles)
			adminRequired.GET("/roles/:id", roleHandler.GetRole)
			adminRequired.POST("/roles", roleHandler.CreateRole)
			adminRequired.PUT("/roles/:id", roleHandler.UpdateRole)
			adminRequired.DELETE("/roles/:id", roleHandler.DeleteRole)
			adminRequired.PUT("/roles/:id/permissions", roleHandler.SetRolePermissions)
			// 权限管理
			adminRequired.GET("/permissions", permissionHandler.GetPermissions)
			// 用户角色分配（管理员用）
			adminRequired.GET("/users/:id/roles", membershipHandler.GetUserRoles)
			adminRequired.POST("/users/:id/roles", membershipHandler.AssignSystemRole)
			adminRequired.DELETE("/users/:id/roles/:role_id", membershipHandler.RemoveSystemRole)
			adminRequired.GET("/users/:id/project-memberships", membershipHandler.GetUserProjectMemberships)
		}

		// 项目路由
		project := api.Group("/projects")
		{
			// 公开的查看路由（需要登录）
			project.GET("", handlers.GetProjects)
			project.GET("/:id", handlers.GetProject)
			project.GET("/:id/tasks", handlers.GetTasks)

			// 需要认证的路由
			projectAuth := project.Group("")
			projectAuth.Use(middleware.AuthRequired(), middleware.RequireProjectMember())
			{
				projectAuth.POST("/:id/tasks", handlers.CreateTask)
			}

			// 项目所有者操作
			projectOwner := project.Group("/:id")
			projectOwner.Use(middleware.AuthRequired(), middleware.RequireProjectOwner())
			{
				projectOwner.PUT("", projectHandler.UpdateProject)
				projectOwner.DELETE("", projectHandler.DeleteProject)
				projectOwner.POST("/archive", projectHandler.ArchiveProject)
				projectOwner.GET("/members", membershipHandler.GetProjectMembers)
				projectOwner.POST("/members", membershipHandler.AddProjectMember)
				projectOwner.PUT("/members/:user_id", membershipHandler.UpdateMemberRole)
				projectOwner.DELETE("/members/:user_id", membershipHandler.RemoveProjectMember)
			}

			// 项目创建（需要登录）
			project.POST("", middleware.AuthRequired(), projectHandler.CreateProject)
		}

		// 任务路由（需要登录）
		taskAuth := api.Group("/tasks")
		taskAuth.Use(middleware.AuthRequired())
		{
			taskAuth.PUT("/:id", handlers.UpdateTask)
			taskAuth.PUT("/:id/status", handlers.UpdateTaskStatus)
			taskAuth.DELETE("/:id", handlers.DeleteTask)

			// 评论路由
			taskAuth.GET("/:id/comments", handlers.GetComments)
			taskAuth.POST("/:id/comments", handlers.CreateComment)
			taskAuth.PUT("/:id/comments/:cid", handlers.UpdateComment)
			taskAuth.DELETE("/:id/comments/:cid", handlers.DeleteComment)
		}

		// SSE事件路由
		api.GET("/events", handlers.GetEvents)
	}
}