package models

import (
	"time"
)

// PersonalAccessToken 个人访问令牌模型
type PersonalAccessToken struct {
	ID           uint       `gorm:"primaryKey"`                            // 令牌ID
	UserID       uint       `gorm:"index:idx_user_id;not null"`            // 所属用户ID
	TokenHash    string     `gorm:"uniqueIndex:idx_token_hash;not null"`  // 令牌 SHA-256 哈希
	TokenPrefix  string     `gorm:"size:20;not null"`                      // 令牌前缀（显示用，如 `tkf_a1b2`）
	Name         string     `gorm:"size:100;not null"`                     // 令牌名称（用户指定）
	Scope        string     `gorm:"size:500;not null"`                    // 权限范围（逗号分隔）
	CreatedAt    time.Time  `gorm:"autoCreateTime"`                        // 创建时间
	LastUsedAt   *time.Time `gorm:"index:idx_last_used"`                   // 最后使用时间
	ExpiresAt    time.Time  `gorm:"index:idx_expires_at;not null"`        // 过期时间
	RevokedAt    *time.Time                                             // 撤销时间（NULL = 未撤销）
}

// TableName 指定表名
func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}

// IsActive 检查令牌是否处于活跃状态
// 活跃状态：未撤销 且 未过期
func (t *PersonalAccessToken) IsActive() bool {
	if t.RevokedAt != nil {
		return false
	}
	return t.ExpiresAt.After(time.Now())
}

// IsRevoked 检查令牌是否已撤销
func (t *PersonalAccessToken) IsRevoked() bool {
	return t.RevokedAt != nil
}

// IsExpired 检查令牌是否已过期
func (t *PersonalAccessToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

// HasScope 检查令牌是否包含指定作用域
func (t *PersonalAccessToken) HasScope(requiredScope string) bool {
	if t.Scope == "" {
		return false
	}
	// 简单检查：scope 字段是否包含 requiredScope
	// 格式: "task:read,task:write,project:read"
	scopes := parseScope(t.Scope)
	for _, s := range scopes {
		if s == requiredScope {
			return true
		}
	}
	return false
}

// parseScope 解析作用域字符串
func parseScope(scope string) []string {
	if scope == "" {
		return []string{}
	}
	var result []string
	var current string
	for _, c := range scope {
		if c == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// TokenInfoResponse Token信息响应
type TokenInfoResponse struct {
	Active    bool      `json:"active"`
	UserID    uint      `json:"user_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Scope     string    `json:"scope,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	TokenName string    `json:"token_name,omitempty"`
}

// CreateTokenRequest 创建令牌请求
type CreateTokenRequest struct {
	Name          string `json:"name" binding:"required,max=100"`      // 令牌名称
	Scope         string `json:"scope" binding:"required,max=500"`    // 权限范围
	ExpiresInDays int    `json:"expires_in_days"`                      // 过期天数（默认90天）
}

// TokenCreatedResponse 令牌创建成功响应
type TokenCreatedResponse struct {
	ID          uint      `json:"id"`
	Token       string    `json:"token"`          // 仅在此响应中显示一次
	Name        string    `json:"name"`
	Scope       string    `json:"scope"`
	TokenPrefix string    `json:"token_prefix"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// TokenListItem 令牌列表项
type TokenListItem struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Scope       string     `json:"scope"`
	TokenPrefix string     `json:"token_prefix"`
	CreatedAt   time.Time  `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`
	Revoked     bool       `json:"revoked"`
}

// TokenListResponse 令牌列表响应
type TokenListResponse struct {
	Tokens []TokenListItem `json:"tokens"`
}