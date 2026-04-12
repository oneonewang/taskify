package tools

import (
	"context"

	"github.com/taskify/backend/internal/mcp"
	"github.com/taskify/backend/internal/repository"
)

// RegisterTasksTools 注册任务工具到MCP服务器
func RegisterTasksTools(server *mcp.Server) {
	// tasks.list - 列出任务
	server.RegisterTool("tasks.list", "列出项目中的任务", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectID, hasProject := args["project_id"]

		taskRepo := repository.NewTaskRepository()
		var tasks []interface{}

		// 简单实现：获取所有任务或按项目筛选
		if hasProject {
			taskList, err := taskRepo.GetByProject(uint(projectID.(float64)), "")
			if err != nil {
				return nil, err
			}
			for _, t := range taskList {
				tasks = append(tasks, map[string]interface{}{
					"id":          t.ID,
					"title":       t.Title,
					"description": t.Description,
					"status":      t.Status,
					"position":    t.Position,
					"assignee_id": t.AssigneeID,
					"project_id":  t.ProjectID,
					"created_at":  t.CreatedAt,
					"updated_at":  t.UpdatedAt,
				})
			}
		} else {
			tasks = []interface{}{}
		}

		return map[string]interface{}{
			"tasks": tasks,
		}, nil
	})

	// tasks.get - 获取任务详情
	server.RegisterTool("tasks.get", "获取任务详情", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		id, ok := args["id"]
		if !ok {
			return nil, nil
		}

		taskRepo := repository.NewTaskRepository()
		task, err := taskRepo.FindByID(uint(id.(float64)))
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"position":    task.Position,
			"assignee_id": task.AssigneeID,
			"project_id":  task.ProjectID,
			"created_at":  task.CreatedAt,
			"updated_at":  task.UpdatedAt,
		}, nil
	})

	// tasks.create - 创建任务
	server.RegisterTool("tasks.create", "创建新任务", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "任务创建功能需要用户认证",
		}, nil
	})

	// tasks.update - 更新任务
	server.RegisterTool("tasks.update", "更新任务信息", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "任务更新功能需要用户认证",
		}, nil
	})

	// tasks.delete - 删除任务
	server.RegisterTool("tasks.delete", "删除任务", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "任务删除功能需要用户认证",
		}, nil
	})

	// tasks.update_status - 更新任务状态
	server.RegisterTool("tasks.update_status", "更新任务状态", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "任务状态更新功能需要用户认证",
		}, nil
	})
}