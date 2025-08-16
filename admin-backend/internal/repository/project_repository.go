package repository

import (
	"context"
	"errors"
	"i18n-flow/internal/domain"

	"gorm.io/gorm"
)

// ProjectRepository 项目仓储实现
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓储实例
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// GetByID 根据ID获取项目
func (r *ProjectRepository) GetByID(ctx context.Context, id uint) (*domain.Project, error) {
	var project domain.Project
	if err := r.db.WithContext(ctx).First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

// GetBySlug 根据Slug获取项目
func (r *ProjectRepository) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	var project domain.Project
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

// GetAll 获取所有项目（分页）
func (r *ProjectRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Project, int64, error) {
	var projects []*domain.Project
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&domain.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// Create 创建项目
func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// Update 更新项目
func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// Delete 删除项目
func (r *ProjectRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Project{}, id).Error
}
