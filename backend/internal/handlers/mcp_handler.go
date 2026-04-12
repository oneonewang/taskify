package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPHandler MCP HTTP处理器
type MCPHandler struct {
	server       *mcp.Server
	sessionStore *SessionStore
}

// SessionStore 会话存储
type SessionStore struct {
	sessions sync.Map
}

// Session MCP会话
type Session struct {
	ID       string
	Server   *mcp.Server
	Writer   http.ResponseWriter
	closed   bool
	mu       sync.RWMutex
}

// NewMCPHandler 创建MCP处理器
func NewMCPHandler(server *mcp.Server) *MCPHandler {
	return &MCPHandler{
		server:       server,
		sessionStore: &SessionStore{},
	}
}

// HandleMCP 处理MCP请求
func (h *MCPHandler) HandleMCP(c *gin.Context) {
	if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "method_not_allowed",
		})
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed_to_read_body",
		})
		return
	}

	// 解析JSON-RPC请求
	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_jsonrpc_request",
		})
		return
	}

	// 获取会话ID
	sessionID := c.GetHeader("Mcp-Session-Id")

	// 处理请求
	ctx := context.Background()
	resp, err := h.handleRequest(ctx, sessionID, &req)
	if err != nil {
		// 返回错误
		errResp := JSONRPCErrorResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, resp)
}

// JSONRPCRequest JSON-RPC请求
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID     interface{}      `json:"id"`
	Method string           `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse JSON-RPC响应
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID     interface{}  `json:"id"`
	Result interface{}  `json:"result,omitempty"`
	Error  *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCErrorResponse JSON-RPC错误响应
type JSONRPCErrorResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID     interface{}  `json:"id"`
	Error  *JSONRPCError `json:"error"`
}

// JSONRPCError JSON-RPC错误
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *MCPHandler) handleRequest(ctx context.Context, sessionID string, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	// 处理 initialize 请求
	if req.Method == "initialize" {
		return h.handleInitialize(ctx, sessionID, req)
	}

	// 处理 ping 请求
	if req.Method == "ping" {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]interface{}{},
		}, nil
	}

	// 处理 tools/list 请求
	if req.Method == "tools/list" {
		return h.handleListTools(ctx, req)
	}

	// 处理 tools/call 请求
	if req.Method == "tools/call" {
		return h.handleCallTool(ctx, req)
	}

	// 处理 resources/list 请求
	if req.Method == "resources/list" {
		return h.handleListResources(ctx, req)
	}

	// 处理 resources/read 请求
	if req.Method == "resources/read" {
		return h.handleReadResource(ctx, req)
	}

	// 未知方法
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error: &JSONRPCError{
			Code:    -32601,
			Message: fmt.Sprintf("method not found: %s", req.Method),
		},
	}, nil
}

func (h *MCPHandler) handleInitialize(ctx context.Context, sessionID string, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	// 解析初始化参数
	var params InitializeParams
	if req.Params != nil {
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return nil, fmt.Errorf("invalid initialize params: %w", err)
		}
	}

	// 返回初始化结果
	result := InitializeResult{
		ProtocolVersion: "2025-06-18",
		Capabilities: ServerCapabilities{
			Tools: &ToolCapability{ListChanged: true},
			Resources: &ResourceCapability{
				Subscribe: true,
				ListChanged: true,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "taskify",
			Version: "1.0.0",
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

func (h *MCPHandler) handleListTools(ctx context.Context, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	// 获取服务器已注册的工具
	var tools []Tool

	result := ListToolsResult{
		Tools: tools,
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

func (h *MCPHandler) handleCallTool(ctx context.Context, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	var params CallToolParams
	if req.Params != nil {
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return nil, fmt.Errorf("invalid call tool params: %w", err)
		}
	}

	// 这里需要调用实际的工具处理器
	// 由于工具注册在服务器中，我们需要通过某种方式获取处理器
	// 暂时返回错误，实际实现需要工具注册表

	result := CallToolResult{
		Content: []TextContent{
			{Text: "tool called: " + params.Name},
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

func (h *MCPHandler) handleListResources(ctx context.Context, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	result := ListResourcesResult{
		Resources: []Resource{},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

func (h *MCPHandler) handleReadResource(ctx context.Context, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	var params ReadResourceParams
	if req.Params != nil {
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return nil, fmt.Errorf("invalid read resource params: %w", err)
		}
	}

	result := ReadResourceResult{
		Contents: []ResourceContents{},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

// InitializeParams 初始化参数
type InitializeParams struct {
	ProtocolVersion string `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo `json:"clientInfo"`
}

// ClientCapabilities 客户端能力
type ClientCapabilities struct {
	Roots *RootsCapability `json:"roots"`
}

// RootsCapability 根目录能力
type RootsCapability struct {
	ListChanged bool `json:"listChanged"`
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult 初始化结果
type InitializeResult struct {
	ProtocolVersion string            `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo        `json:"serverInfo"`
}

// ServerCapabilities 服务器能力
type ServerCapabilities struct {
	Tools     *ToolCapability      `json:"tools,omitempty"`
	Resources *ResourceCapability  `json:"resources,omitempty"`
}

// ToolCapability 工具能力
type ToolCapability struct {
	ListChanged bool `json:"listChanged"`
}

// ResourceCapability 资源能力
type ResourceCapability struct {
	Subscribe   bool `json:"subscribe"`
	ListChanged bool `json:"listChanged"`
}

// ServerInfo 服务器信息
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ListToolsResult 工具列表结果
type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

// Tool 工具定义
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallToolParams 调用工具参数
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult 调用工具结果
type CallToolResult struct {
	Content []TextContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// TextContent 文本内容
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ListResourcesResult 资源列表结果
type ListResourcesResult struct {
	Resources []Resource `json:"resources"`
}

// Resource 资源定义
type Resource struct {
	URI       string `json:"uri"`
	Name      string `json:"name"`
	MIMEType  string `json:"mimeType,omitempty"`
}

// ReadResourceParams 读取资源参数
type ReadResourceParams struct {
	URI string `json:"uri"`
}

// ReadResourceResult 读取资源结果
type ReadResourceResult struct {
	Contents []ResourceContents `json:"contents"`
}

// ResourceContents 资源内容
type ResourceContents struct {
	URI     string `json:"uri"`
	MIMEType string `json:"mimeType,omitempty"`
	Blob    string `json:"blob,omitempty"`
	Text    string `json:"text,omitempty"`
}

// ServeHTTP 实现 http.Handler 接口
func (h *MCPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 处理 CORS 预检请求
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Mcp-Session-Id")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 处理请求
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// 解析请求
	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
		return
	}

	// 获取会话ID
	sessionID := r.Header.Get("Mcp-Session-Id")

	// 处理请求
	ctx := context.Background()
	resp, err := h.handleRequest(ctx, sessionID, &req)
	if err != nil {
		errResp := JSONRPCErrorResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BufferResponseWriter 缓冲响应写入器
type BufferResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	code   int
}

func (w *BufferResponseWriter) Header() http.Header {
	return w.header
}

func (w *BufferResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *BufferResponseWriter) WriteHeader(code int) {
	w.code = code
}

func (w *BufferResponseWriter) Result() ([]byte, int) {
	return w.body.Bytes(), w.code
}