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
func (r *ProjectRepository) GetAll(ctx context.Context, limit, offset int, keyword string) ([]*domain.Project, int64, error) {
	var projects []*domain.Project
	var total int64

	// 构建基础查询条件，添加软删除过滤
	baseQuery := r.db.WithContext(ctx).Model(&domain.Project{}).Where("deleted_at IS NULL")
	
	// 构建搜索条件
	var query *gorm.DB
	if keyword != "" {
		// 优化搜索：优先匹配名称，然后是slug，最后是描述
		query = baseQuery.Where("(name LIKE ? OR slug LIKE ? OR description LIKE ?)", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	} else {
		query = baseQuery
	}

	// 优化：使用相同的查询条件计算总数
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回
	if total == 0 {
		return []*domain.Project{}, 0, nil
	}

	// 获取分页数据，添加排序优化查询性能
	if err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
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
