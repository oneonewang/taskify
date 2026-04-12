package tools

import (
	"context"

	"github.com/taskify/backend/internal/mcp"
	"github.com/taskify/backend/internal/repository"
)

// RegisterCommentsTools 注册评论工具到MCP服务器
func RegisterCommentsTools(server *mcp.Server) {
	// comments.list - 列出评论
	server.RegisterTool("comments.list", "列出任务的评论", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		taskID, ok := args["task_id"]
		if !ok {
			return nil, nil
		}

		commentRepo := repository.NewCommentRepository()
		comments, err := commentRepo.GetByTask(uint(taskID.(float64)))
		if err != nil {
			return nil, err
		}

		var result []interface{}
		for _, c := range comments {
			result = append(result, map[string]interface{}{
				"id":         c.ID,
				"content":    c.Content,
				"task_id":    c.TaskID,
				"user_id":    c.UserID,
				"created_at": c.CreatedAt,
				"updated_at": c.UpdatedAt,
			})
		}

		return map[string]interface{}{
			"comments": result,
		}, nil
	})

	// comments.get - 获取评论详情
	server.RegisterTool("comments.get", "获取评论详情", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		id, ok := args["id"]
		if !ok {
			return nil, nil
		}

		commentRepo := repository.NewCommentRepository()
		comment, err := commentRepo.FindByID(uint(id.(float64)))
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"task_id":    comment.TaskID,
			"user_id":    comment.UserID,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
		}, nil
	})

	// comments.create - 创建评论
	server.RegisterTool("comments.create", "创建评论", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "评论创建功能需要用户认证",
		}, nil
	})

	// comments.update - 更新评论
	server.RegisterTool("comments.update", "更新评论", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "评论更新功能需要用户认证",
		}, nil
	})

	// comments.delete - 删除评论
	server.RegisterTool("comments.delete", "删除评论", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		// 占位符实现 - 需要用户认证
		return map[string]interface{}{
			"message": "评论删除功能需要用户认证",
		}, nil
	})
}