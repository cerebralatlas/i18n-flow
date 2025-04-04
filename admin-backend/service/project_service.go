package service

import (
	"errors"
	"i18n-flow/model"
	"i18n-flow/model/db"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// ProjectService 项目服务
type ProjectService struct{}

// ProjectRequest 项目请求参数
type ProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(req ProjectRequest) (*ProjectResponse, error) {
	// 处理slug
	projectSlug := req.Slug
	if projectSlug == "" {
		projectSlug = slug.Make(req.Name)
	} else {
		projectSlug = slug.Make(projectSlug)
	}

	// 检查slug是否已存在
	var existingProject model.Project
	if err := db.DB.Where("slug = ?", projectSlug).First(&existingProject).Error; err == nil {
		return nil, errors.New("项目标识已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	project := model.Project{
		Name:        req.Name,
		Description: req.Description,
		Slug:        projectSlug,
		Status:      "active",
	}

	if err := db.DB.Create(&project).Error; err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Slug:        project.Slug,
		Status:      project.Status,
		CreatedAt:   project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetProjects 获取项目列表
func (s *ProjectService) GetProjects(page, pageSize int, keyword string) ([]ProjectResponse, int64, error) {
	var projects []model.Project
	var total int64

	query := db.DB.Model(&model.Project{})

	// 如果有关键字，添加搜索条件
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	// 格式化响应
	var response []ProjectResponse
	for _, p := range projects {
		response = append(response, ProjectResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Slug:        p.Slug,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, total, nil
}

// GetProjectByID 根据ID获取项目
func (s *ProjectService) GetProjectByID(id uint) (*ProjectResponse, error) {
	var project model.Project

	if err := db.DB.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Slug:        project.Slug,
		Status:      project.Status,
		CreatedAt:   project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(id uint, req ProjectRequest) (*ProjectResponse, error) {
	var project model.Project

	if err := db.DB.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 处理slug
	projectSlug := req.Slug
	if projectSlug == "" {
		projectSlug = slug.Make(req.Name)
	} else {
		projectSlug = slug.Make(projectSlug)
	}

	// 如果slug变更，检查是否与其他项目冲突
	if projectSlug != project.Slug {
		var existingProject model.Project
		if err := db.DB.Where("slug = ? AND id != ?", projectSlug, id).First(&existingProject).Error; err == nil {
			return nil, errors.New("项目标识已存在")
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	// 更新项目信息
	project.Name = req.Name
	project.Description = req.Description
	project.Slug = projectSlug

	if err := db.DB.Save(&project).Error; err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Slug:        project.Slug,
		Status:      project.Status,
		CreatedAt:   project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// DeleteProject 删除项目
func (s *ProjectService) DeleteProject(id uint) error {
	// 更新project的status
	db.DB.Model(&model.Project{}).Where("id = ?", id).Update("status", "archived")

	if err := db.DB.Delete(&model.Project{}, id).Error; err != nil {
		return err
	}
	return nil
}
