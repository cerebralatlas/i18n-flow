package service

import (
	"context"
	"i18n-flow/internal/domain"
	"strings"

	"github.com/gosimple/slug"
)

// ProjectService 项目服务实现
type ProjectService struct {
	projectRepo       domain.ProjectRepository
	userRepo          domain.UserRepository
	projectMemberRepo domain.ProjectMemberRepository
}

// NewProjectService 创建项目服务实例
func NewProjectService(
	projectRepo domain.ProjectRepository,
	userRepo domain.UserRepository,
	projectMemberRepo domain.ProjectMemberRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo:       projectRepo,
		userRepo:          userRepo,
		projectMemberRepo: projectMemberRepo,
	}
}

// Create 创建项目
func (s *ProjectService) Create(ctx context.Context, req domain.CreateProjectRequest) (*domain.Project, error) {
	// 生成slug
	projectSlug := slug.Make(req.Name)
	if projectSlug == "" {
		return nil, domain.ErrInvalidSlug
	}

	// 检查slug是否已存在
	existingProject, err := s.projectRepo.GetBySlug(ctx, projectSlug)
	if err == nil && existingProject != nil {
		return nil, domain.ErrProjectExists
	}

	// 创建项目
	project := &domain.Project{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Slug:        projectSlug,
		Status:      "active",
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetByID 根据ID获取项目
func (s *ProjectService) GetByID(ctx context.Context, id uint) (*domain.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

// GetAll 获取所有项目
func (s *ProjectService) GetAll(ctx context.Context, limit, offset int, keyword string) ([]*domain.Project, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.projectRepo.GetAll(ctx, limit, offset, keyword)
}

// Update 更新项目
func (s *ProjectService) Update(ctx context.Context, id uint, req domain.UpdateProjectRequest) (*domain.Project, error) {
	// 获取现有项目
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		project.Name = strings.TrimSpace(req.Name)
		// 如果名称改变，重新生成slug
		newSlug := slug.Make(req.Name)
		if newSlug != project.Slug {
			// 检查新slug是否已存在
			existingProject, err := s.projectRepo.GetBySlug(ctx, newSlug)
			if err == nil && existingProject != nil && existingProject.ID != project.ID {
				return nil, domain.ErrProjectExists
			}
			project.Slug = newSlug
		}
	}

	if req.Description != "" {
		project.Description = strings.TrimSpace(req.Description)
	}

	if req.Status != "" {
		if req.Status != "active" && req.Status != "archived" {
			return nil, domain.ErrInvalidInput
		}
		project.Status = req.Status
	}

	// 保存更新
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete 删除项目
func (s *ProjectService) Delete(ctx context.Context, id uint) error {
	// 检查项目是否存在
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 删除项目
	return s.projectRepo.Delete(ctx, id)
}

// GetAccessibleProjects 获取用户可访问的项目列表
func (s *ProjectService) GetAccessibleProjects(ctx context.Context, userID uint, limit, offset int, keyword string) ([]*domain.Project, int64, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	// 如果是管理员，返回所有项目
	if user.Role == "admin" {
		return s.GetAll(ctx, limit, offset, keyword)
	}

	// 普通用户：获取用户参与的项目
	members, err := s.projectMemberRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	// 如果用户没有参与任何项目
	if len(members) == 0 {
		return []*domain.Project{}, 0, nil
	}

	// 提取项目ID列表
	projectIDs := make([]uint, len(members))
	for i, member := range members {
		projectIDs[i] = member.ProjectID
	}

	// 获取用户参与的项目详情
	var projects []*domain.Project
	var filteredProjects []*domain.Project

	for _, projectID := range projectIDs {
		project, err := s.projectRepo.GetByID(ctx, projectID)
		if err != nil {
			continue // 跳过不存在的项目
		}
		projects = append(projects, project)
	}

	// 应用关键词过滤
	if keyword != "" {
		keyword = strings.ToLower(keyword)
		for _, project := range projects {
			if strings.Contains(strings.ToLower(project.Name), keyword) ||
				strings.Contains(strings.ToLower(project.Description), keyword) {
				filteredProjects = append(filteredProjects, project)
			}
		}
	} else {
		filteredProjects = projects
	}

	total := int64(len(filteredProjects))

	// 应用分页
	start := offset
	end := offset + limit
	if start > len(filteredProjects) {
		start = len(filteredProjects)
	}
	if end > len(filteredProjects) {
		end = len(filteredProjects)
	}

	paginatedProjects := filteredProjects[start:end]
	return paginatedProjects, total, nil
}
