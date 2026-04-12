package repository

import (
	"time"

	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// TokenRepository 令牌数据访问层
type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository 创建令牌仓库
func NewTokenRepository() *TokenRepository {
	return &TokenRepository{db: GetDB()}
}

// Create 创建令牌
func (r *TokenRepository) Create(token *models.PersonalAccessToken) error {
	return r.db.Create(token).Error
}

// GetByHash 根据哈希查找令牌（不检查状态）
func (r *TokenRepository) GetByHash(hash string) (*models.PersonalAccessToken, error) {
	var token models.PersonalAccessToken
	err := r.db.Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByID 根据ID查找令牌
func (r *TokenRepository) GetByID(id uint) (*models.PersonalAccessToken, error) {
	var token models.PersonalAccessToken
	err := r.db.First(&token, id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// ListByUser 列出用户的所有令牌
func (r *TokenRepository) ListByUser(userID uint) ([]models.PersonalAccessToken, error) {
	var tokens []models.PersonalAccessToken
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// Revoke 撤销令牌
func (r *TokenRepository) Revoke(id uint) error {
	now := time.Now()
	return r.db.Model(&models.PersonalAccessToken{}).
		Where("id = ?", id).
		Update("revoked_at", now).Error
}

// UpdateLastUsed 更新最后使用时间
func (r *TokenRepository) UpdateLastUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&models.PersonalAccessToken{}).
		Where("id = ?", id).
		Update("last_used_at", now).Error
}

// Delete 删除令牌
func (r *TokenRepository) Delete(id uint) error {
	return r.db.Delete(&models.PersonalAccessToken{}, id).Error
}

// GetActiveTokenByHash 获取活跃令牌（未撤销且未过期）
func (r *TokenRepository) GetActiveTokenByHash(hash string) (*models.PersonalAccessToken, error) {
	var token models.PersonalAccessToken
	now := time.Now()
	err := r.db.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", hash, now).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}