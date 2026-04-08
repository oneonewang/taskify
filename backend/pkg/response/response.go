package response

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功响应结构
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage 返回成功响应带消息
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(200, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, httpStatus int, code string, message string) {
	c.JSON(httpStatus, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorWithDetails 返回错误响应带详情
func ErrorWithDetails(c *gin.Context, httpStatus int, code string, message string, details interface{}) {
	c.JSON(httpStatus, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// 预定义错误码
const (
	CodeValidationError = "VALIDATION_ERROR"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeNotFound         = "NOT_FOUND"
	CodeInternalError    = "INTERNAL_ERROR"
)
