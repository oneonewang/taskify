package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server MCP服务器封装
type Server struct {
	server    *mcp.Server
	tools     map[string]ToolHandler
	resources map[string]ResourceHandler
	logger    *slog.Logger
	mu        sync.RWMutex
}

// ToolHandler 工具处理器函数类型
type ToolHandler func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// ResourceHandler 资源处理器函数类型
type ResourceHandler func(ctx context.Context, uri string) (interface{}, error)

// Implementation 服务器实现信息
type Implementation struct {
	Name    string
	Version string
}

// NewServer 创建MCP服务器
func NewServer(name, version string) *Server {
	logger := slog.Default()

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    name,
		Version: version,
	}, &mcp.ServerOptions{
		Logger: logger,
	})

	return &Server{
		server:    mcpServer,
		tools:     make(map[string]ToolHandler),
		resources: make(map[string]ResourceHandler),
		logger:    logger,
	}
}

// RegisterTool 注册工具
func (s *Server) RegisterTool(name, description string, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tools[name] = handler

	// 使用 SDK 的 AddTool 注册
	mcp.AddTool[map[string]interface{}, any](s.server, &mcp.Tool{
		Name:        name,
		Description: description,
		InputSchema: json.RawMessage(`{"type": "object", "properties": {}}`),
	}, func(ctx context.Context, req *mcp.CallToolRequest, args map[string]interface{}) (result *mcp.CallToolResult, output any, err error) {
		if req.Params.Arguments != nil {
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{Text: "invalid arguments: " + err.Error()}},
					IsError: true,
				}, nil, nil
			}
		}

		res, err := handler(ctx, args)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
				IsError: true,
			}, nil, nil
		}

		return nil, res, nil
	})
}

// RegisterResource 注册资源
func (s *Server) RegisterResource(uri, name, mimeType string, handler ResourceHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.resources[uri] = handler

	s.server.AddResource(&mcp.Resource{
		URI:      uri,
		Name:     name,
		MIMEType: mimeType,
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		result, err := handler(ctx, uri)
		if err != nil {
			return nil, err
		}

		content, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resource: %w", err)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:     uri,
					MIMEType: mimeType,
					Blob:    content,
				},
			},
		}, nil
	})
}

// GetServer 获取底层MCP服务器
func (s *Server) GetServer() *mcp.Server {
	return s.server
}