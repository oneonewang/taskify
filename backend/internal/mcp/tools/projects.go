package tools

import (
	"context"

	"github.com/taskify/backend/internal/mcp"
	"github.com/taskify/backend/internal/repository"
)

// RegisterProjectsTools 注册项目工具到MCP服务器
func RegisterProjectsTools(server *mcp.Server) {
	// projects.list - 列出项目
	server.RegisterTool("projects.list", "列出所有项目", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectRepo := repository.NewProjectRepository()
		var projects []interface{}

		if userID, hasUserID := args["user_id"]; hasUserID {
			projectList, err := projectRepo.GetByUser(uint(userID.(float64)))
			if err != nil {
				return nil, err
			}
			for _, p := range projectList {
				projects = append(projects, map[string]interface{}{
					"id":          p.ID,
					"name":        p.Name,
					"description": p.Description,
					"is_archived": p.IsArchived,
					"created_at":  p.CreatedAt,
					"updated_at":  p.UpdatedAt,
				})
			}
		} else {
			projectList, err := projectRepo.GetAll()
			if err != nil {
				return nil, err
			}
			for _, p := range projectList {
				projects = append(projects, map[string]interface{}{
					"id":          p.ID,
					"name":        p.Name,
					"description": p.Description,
					"is_archived": p.IsArchived,
					"created_at":  p.CreatedAt,
					"updated_at":  p.UpdatedAt,
				})
			}
		}

		return map[string]interface{}{
			"projects": projects,
		}, nil
	})

	// projects.get - 获取项目详情
	server.RegisterTool("projects.get", "获取项目详情", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		id, ok := args["id"]
		if !ok {
			return nil, nil
		}

		projectRepo := repository.NewProjectRepository()
		project, err := projectRepo.FindByID(uint(id.(float64)))
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":          project.ID,
			"name":        project.Name,
			"description": project.Description,
			"is_archived": project.IsArchived,
			"created_at":  project.CreatedAt,
			"updated_at":  project.UpdatedAt,
		}, nil
	})

	// projects.create - 创建项目
	server.RegisterTool("projects.create", "创建新项目", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		name, ok := args["name"].(string)
		if !ok || name == "" {
			return nil, nil
		}

		// 占位符实现 - 需要用户认证
		_ = name
		return map[string]interface{}{
			"message": "项目创建功能需要用户认证",
		}, nil
	})

	// projects.update - 更新项目
	server.RegisterTool("projects.update", "更新项目信息", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "项目更新功能需要用户认证",
		}, nil
	})

	// projects.delete - 删除项目
	server.RegisterTool("projects.delete", "删除项目", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "项目删除功能需要用户认证",
		}, nil
	})
}