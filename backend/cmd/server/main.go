package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/config"
	"github.com/taskify/backend/internal/handlers"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/mcp"
	"github.com/taskify/backend/internal/mcp/tools"
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

	// 创建MCP服务器
	mcpServer := mcp.NewServer("taskify", "1.0.0")
	tools.RegisterTasksTools(mcpServer)
	tools.RegisterProjectsTools(mcpServer)
	tools.RegisterCommentsTools(mcpServer)

	// 创建MCP处理器
	mcpHandler := handlers.NewMCPHandler(mcpServer.GetServer())

	// 创建Token处理器
	tokenHandler := handlers.NewTokenHandler()

	// 注册路由
	projectHandler := handlers.NewProjectHandler()
	taskHandler := handlers.NewTaskHandler()
	commentHandler := handlers.NewCommentHandler()
	registerRoutes(r, projectHandler, taskHandler, commentHandler, mcpHandler, tokenHandler, "http://localhost:"+cfg.Port)

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// registerRoutes 注册所有路由
func registerRoutes(r *gin.Engine, projectHandler *handlers.ProjectHandler, taskHandler *handlers.TaskHandler, commentHandler *handlers.CommentHandler, mcpHandler *handlers.MCPHandler, tokenHandler *handlers.TokenHandler, issuerURL string) {
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
			adminUserHandler := handlers.NewAdminUserHandler()
			adminRequired.GET("/users", adminUserHandler.ListUsers)
			adminRequired.POST("/users/:id/reset-password", adminUserHandler.ResetPassword)
			adminRequired.POST("/users/:id/disable", adminUserHandler.DisableUser)
			adminRequired.POST("/users/:id/enable", adminUserHandler.EnableUser)
			adminRequired.POST("/users/import", adminUserHandler.ImportUsers)
			adminRequired.GET("/users/import/template", adminUserHandler.GetImportTemplate)
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
			// 审计日志
			auditLogHandler := handlers.NewAuditLogHandler()
			adminRequired.GET("/audit-logs", auditLogHandler.ListAuditLogs)
		}

		// 项目路由（需要登录）
		project := api.Group("/projects")
		project.Use(middleware.AuthRequired())
		{
			project.GET("", projectHandler.GetProjects)
			project.GET("/my", projectHandler.GetMyProjects)
			project.GET("/:id", projectHandler.GetProject)
			project.GET("/:id/tasks", taskHandler.GetTasks)
			project.GET("/:id/tasks/:taskId", taskHandler.GetTaskByProject)

			// 项目所有者操作
			projectOwner := project.Group("/:id")
			projectOwner.Use(middleware.RequireProjectOwner())
			{
				projectOwner.PUT("", projectHandler.UpdateProject)
				projectOwner.DELETE("", projectHandler.DeleteProject)
				projectOwner.POST("/archive", projectHandler.ArchiveProject)
				projectOwner.POST("/members/batch", membershipHandler.BatchAddProjectMembers)
				projectOwner.POST("/members", membershipHandler.AddProjectMember)
				projectOwner.PUT("/members/:user_id", membershipHandler.UpdateMemberRole)
				projectOwner.DELETE("/members/:user_id", membershipHandler.RemoveProjectMember)
			}

			// 项目成员操作（需要是项目成员）
			projectMember := project.Group("/:id")
			projectMember.Use(middleware.RequireProjectMember())
			{
				projectMember.GET("/members", membershipHandler.GetProjectMembers)
				projectMember.POST("/tasks", taskHandler.CreateTask)
				projectMember.GET("/my-membership", membershipHandler.GetMyProjectMembership)
			}

			// 项目创建
			project.POST("", projectHandler.CreateProject)
		}

		// 任务路由（需要登录）
		taskAuth := api.Group("/tasks")
		taskAuth.Use(middleware.AuthRequired())
		{
			taskAuth.PUT("/:id", taskHandler.UpdateTask)
			taskAuth.PUT("/:id/status", taskHandler.UpdateTaskStatus)
			taskAuth.DELETE("/:id", taskHandler.DeleteTask)

			// 评论路由
			taskAuth.GET("/:id/comments", commentHandler.GetComments)
			taskAuth.POST("/:id/comments", commentHandler.CreateComment)
			taskAuth.PUT("/:id/comments/:cid", commentHandler.UpdateComment)
			taskAuth.DELETE("/:id/comments/:cid", commentHandler.DeleteComment)
		}

		// SSE事件路由
		api.GET("/events", handlers.GetEvents)
	}

	// OAuth Discovery 端点（无需认证）
	discoveryHandler := handlers.NewDiscoveryHandler(issuerURL)
	r.GET("/.well-known/oauth-authorization-server", discoveryHandler.GetOAuthAuthorizationServer)
	r.GET("/.well-known/openid-configuration", discoveryHandler.GetOpenIDConfiguration)

	// OAuth Token 端点（公开）
	api.GET("/oauth/token/info", tokenHandler.GetTokenInfo)

	// OAuth Token 管理端点（需要会话认证 - 用于前端Web UI）
	oauth := api.Group("/oauth")
	oauth.Use(middleware.AuthRequired())
	{
		oauth.POST("/tokens", tokenHandler.CreateToken)
		oauth.GET("/tokens", tokenHandler.ListTokens)
		oauth.DELETE("/tokens/:id", tokenHandler.RevokeToken)
	}

	// MCP 端点（需要认证）
	r.POST("/mcp", middleware.PATAuthMiddleware(), mcpHandler.HandleMCP)
}
