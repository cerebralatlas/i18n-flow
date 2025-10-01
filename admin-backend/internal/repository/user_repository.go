package repository

import (
	"context"
	"errors"
	"i18n-flow/internal/domain"

	"gorm.io/gorm"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetAll 获取用户列表
func (r *UserRepository) GetAll(ctx context.Context, limit, offset int, keyword string) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.User{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}
