package controller

import (
	"i18n-flow/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProjectController 项目控制器
type ProjectController struct {
	projectService service.ProjectService
}

// NewProjectController 创建项目控制器
func NewProjectController() *ProjectController {
	return &ProjectController{
		projectService: service.ProjectService{},
	}
}

// CreateProject 创建项目
// @Summary      创建新项目
// @Description  创建一个新的翻译项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        project  body      service.ProjectRequest  true  "项目信息"
// @Success      201      {object}  service.ProjectResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects [post]
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req service.ProjectRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := c.projectService.CreateProject(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, project)
}

// GetProjects 获取项目列表
// @Summary      获取项目列表
// @Description  获取所有项目列表，支持分页和搜索
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "页码，默认1"
// @Param        page_size  query     int     false  "每页数量，默认10"
// @Param        keyword  query     string  false  "搜索关键词"
// @Success      200      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects [get]
func (c *ProjectController) GetProjects(ctx *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	keyword := ctx.Query("keyword")

	projects, total, err := c.projectService.GetProjects(page, pageSize, keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  projects,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetProjectByID 根据ID获取项目
// @Summary      根据ID获取项目
// @Description  通过项目ID获取项目详细信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "项目ID"
// @Success      200  {object}  service.ProjectResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects/detail/{id} [get]
func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	project, err := c.projectService.GetProjectByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, project)
}

// UpdateProject 更新项目
// @Summary      更新项目
// @Description  通过项目ID更新项目信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id       path      int                     true  "项目ID"
// @Param        project  body      service.ProjectRequest  true  "项目信息"
// @Success      200      {object}  service.ProjectResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects/update/{id} [put]
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	var req service.ProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := c.projectService.UpdateProject(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, project)
}

// DeleteProject 删除项目
// @Summary      删除项目
// @Description  通过项目ID删除项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "项目ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /projects/delete/{id} [delete]
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	if err := c.projectService.DeleteProject(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "项目删除成功"})
}
