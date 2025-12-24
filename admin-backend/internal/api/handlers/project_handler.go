package handlers

import (
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectService domain.ProjectService
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler(projectService domain.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// Create 创建项目
// @Summary      创建项目
// @Description  创建新的翻译项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        project  body      domain.CreateProjectRequest  true  "项目信息"
// @Success      201      {object}  domain.Project
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects [post]
func (h *ProjectHandler) Create(ctx *gin.Context) {
	var req domain.CreateProjectRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "未找到用户信息")
		return
	}

	project, err := h.projectService.Create(ctx.Request.Context(), req, userID.(uint64))
	if err != nil {
		switch err {
		case domain.ErrProjectExists:
			response.Conflict(ctx, err.Error())
		case domain.ErrInvalidSlug:
			response.BadRequest(ctx, err.Error())
		default:
			response.InternalServerError(ctx, "创建项目失败")
		}
		return
	}

	response.Created(ctx, project)
}

// GetByID 根据ID获取项目
// @Summary      获取项目详情
// @Description  根据项目ID获取项目详细信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "项目ID"
// @Success      200  {object}  domain.Project
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects/detail/{id} [get]
func (h *ProjectHandler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的项目ID")
		return
	}

	project, err := h.projectService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			response.NotFound(ctx, err.Error())
		default:
			response.InternalServerError(ctx, "获取项目失败")
		}
		return
	}

	response.Success(ctx, project)
}

// GetAll 获取项目列表
// @Summary      获取项目列表
// @Description  获取项目列表，支持分页和关键词搜索
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "页码"  default(1)
// @Param        page_size query     int     false  "每页数量"  default(10)
// @Param        keyword   query     string  false  "搜索关键词"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects [get]
func (h *ProjectHandler) GetAll(ctx *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	keyword := ctx.DefaultQuery("keyword", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	projects, total, err := h.projectService.GetAll(ctx.Request.Context(), pageSize, offset, keyword)
	if err != nil {
		response.InternalServerError(ctx, "获取项目列表失败")
		return
	}

	meta := &response.Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}

	response.SuccessWithMeta(ctx, projects, meta)
}

// GetAccessibleProjects 获取用户可访问的项目列表
// @Summary      获取用户可访问的项目列表
// @Description  根据用户权限返回可访问的项目列表，管理员返回所有项目，普通用户返回参与的项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "页码"        default(1)
// @Param        page_size query     int     false  "每页数量"     default(10)
// @Param        keyword   query     string  false  "搜索关键词"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects/accessible [get]
func (h *ProjectHandler) GetAccessibleProjects(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	keyword := ctx.DefaultQuery("keyword", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	projects, total, err := h.projectService.GetAccessibleProjects(ctx.Request.Context(), userID.(uint64), pageSize, offset, keyword)
	if err != nil {
		response.InternalServerError(ctx, "获取可访问项目列表失败")
		return
	}

	meta := &response.Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}

	response.SuccessWithMeta(ctx, projects, meta)
}

// Update 更新项目
// @Summary      更新项目
// @Description  更新项目信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "项目ID"
// @Param        project  body      domain.UpdateProjectRequest  true  "项目信息"
// @Success      200      {object}  domain.Project
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects/update/{id} [put]
func (h *ProjectHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的项目ID")
		return
	}

	var req domain.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "未找到用户信息")
		return
	}

	project, err := h.projectService.Update(ctx.Request.Context(), id, req, userID.(uint64))
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			response.NotFound(ctx, err.Error())
		case domain.ErrProjectExists, domain.ErrInvalidInput:
			response.BadRequest(ctx, err.Error())
		default:
			response.InternalServerError(ctx, "更新项目失败")
		}
		return
	}

	response.Success(ctx, project)
}

// Delete 删除项目
// @Summary      删除项目
// @Description  删除指定的项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "项目ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects/delete/{id} [delete]
func (h *ProjectHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的项目ID")
		return
	}

	err = h.projectService.Delete(ctx.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			response.NotFound(ctx, err.Error())
		default:
			response.InternalServerError(ctx, "删除项目失败")
		}
		return
	}

	response.NoContent(ctx)
}
