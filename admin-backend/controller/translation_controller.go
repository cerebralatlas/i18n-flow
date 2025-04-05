package controller

import (
	"i18n-flow/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TranslationController 翻译控制器
type TranslationController struct {
	translationService service.TranslationService
}

// NewTranslationController 创建翻译控制器
func NewTranslationController() *TranslationController {
	return &TranslationController{
		translationService: service.TranslationService{},
	}
}

// CreateTranslation 创建翻译
// @Summary      创建翻译条目
// @Description  为特定项目和语言创建新的翻译条目
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        translation  body      service.TranslationRequest  true  "翻译信息"
// @Success      201          {object}  service.TranslationResponse
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations [post]
func (c *TranslationController) CreateTranslation(ctx *gin.Context) {
	var req service.TranslationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translation, err := c.translationService.CreateTranslation(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, translation)
}

// BatchCreateTranslations 批量创建翻译
// @Summary      批量创建翻译
// @Description  为特定项目批量创建多种语言的翻译
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        translations  body      service.BatchTranslationRequest  true  "批量翻译信息，translations 字段为语言代码到翻译值的映射，如 {'zh-CN':'中文值','en':'English value'}"
// @Success      201           {array}   service.TranslationResponse
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/batch [post]
func (c *TranslationController) BatchCreateTranslations(ctx *gin.Context) {
	var req service.BatchTranslationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translations, err := c.translationService.BatchCreateTranslations(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, translations)
}

// GetTranslationsByProject 获取项目的所有翻译
// @Summary      获取项目翻译
// @Description  获取指定项目的所有翻译条目，支持分页和搜索
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int     true   "项目ID"
// @Param        page        query     int     false  "页码，默认1"
// @Param        page_size   query     int     false  "每页数量，默认10"
// @Param        keyword     query     string  false  "搜索关键词"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/by-project/{project_id} [get]
func (c *TranslationController) GetTranslationsByProject(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("project_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	keyword := ctx.Query("keyword")

	translations, total, err := c.translationService.GetTranslationsByProject(uint(projectID), page, pageSize, keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  translations,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetTranslationByID 根据ID获取翻译
// @Summary      获取翻译详情
// @Description  通过ID获取特定翻译条目的详细信息
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "翻译ID"
// @Success      200  {object}  service.TranslationResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/{id} [get]
func (c *TranslationController) GetTranslationByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的翻译ID"})
		return
	}

	translation, err := c.translationService.GetTranslationByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, translation)
}

// UpdateTranslation 更新翻译
// @Summary      更新翻译
// @Description  更新特定翻译条目的内容
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        id           path      int                        true  "翻译ID"
// @Param        translation  body      service.TranslationRequest  true  "翻译信息"
// @Success      200          {object}  service.TranslationResponse
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/{id} [put]
func (c *TranslationController) UpdateTranslation(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的翻译ID"})
		return
	}

	var req service.TranslationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translation, err := c.translationService.UpdateTranslation(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, translation)
}

// DeleteTranslation 删除翻译
// @Summary      删除翻译
// @Description  删除特定的翻译条目
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "翻译ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/{id} [delete]
func (c *TranslationController) DeleteTranslation(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的翻译ID"})
		return
	}

	if err := c.translationService.DeleteTranslation(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "翻译删除成功"})
}

// ExportTranslations 导出项目的翻译
// @Summary      导出项目翻译
// @Description  导出指定项目的所有翻译，支持不同格式
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int     true   "项目ID"
// @Param        format      query     string  false  "导出格式，默认json"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /exports/project/{project_id} [get]
func (c *TranslationController) ExportTranslations(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("project_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	format := ctx.DefaultQuery("format", "json")

	data, err := c.translationService.ExportTranslations(uint(projectID), format)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// ImportTranslations 导入项目的翻译
// @Summary      导入项目翻译
// @Description  为指定项目导入翻译数据
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int                           true  "项目ID"
// @Param        data        body      map[string]map[string]string  true  "翻译数据"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /imports/project/{project_id} [post]
func (c *TranslationController) ImportTranslations(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("project_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	var data map[string]map[string]string
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := c.translationService.ImportTranslations(uint(projectID), data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "翻译导入成功",
		"count":   count,
	})
}

// BatchDeleteTranslations 批量删除翻译
// @Summary      批量删除翻译
// @Description  批量删除多个翻译条目
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        ids  body      []uint  true  "翻译ID数组"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/batch-delete [post]
func (c *TranslationController) BatchDeleteTranslations(ctx *gin.Context) {
	var ids []uint
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据，请提供ID数组"})
		return
	}

	count, err := c.translationService.BatchDeleteTranslations(ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "批量删除翻译成功",
		"count":   count,
	})
}

// GetTranslationMatrix 获取项目的翻译矩阵
// @Summary      获取项目翻译矩阵
// @Description  获取指定项目的翻译矩阵，按键名分组，每组包含各语言的翻译，支持分页和搜索
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int     true   "项目ID"
// @Param        page        query     int     false  "页码，默认1"
// @Param        page_size   query     int     false  "每页数量，默认10"
// @Param        keyword     query     string  false  "搜索关键词"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/matrix/by-project/{project_id} [get]
func (c *TranslationController) GetTranslationMatrix(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("project_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	keyword := ctx.Query("keyword")

	matrix, total, err := c.translationService.GetTranslationMatrix(uint(projectID), page, pageSize, keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  matrix,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
