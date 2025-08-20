package handlers

import (
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TranslationHandler 翻译处理器
type TranslationHandler struct {
	translationService domain.TranslationService
}

// NewTranslationHandler 创建翻译处理器
func NewTranslationHandler(translationService domain.TranslationService) *TranslationHandler {
	return &TranslationHandler{
		translationService: translationService,
	}
}

// Create 创建翻译
// @Summary      创建翻译
// @Description  创建新的翻译
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        translation  body      domain.CreateTranslationRequest  true  "翻译信息"
// @Success      201          {object}  domain.Translation
// @Failure      400          {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations [post]
func (h *TranslationHandler) Create(ctx *gin.Context) {
	var req domain.CreateTranslationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translation, err := h.translationService.Create(ctx.Request.Context(), req)
	if err != nil {
		if err == domain.ErrProjectNotFound || err == domain.ErrLanguageNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建翻译失败"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, translation)
}

// CreateBatch 批量创建翻译
// @Summary      批量创建翻译
// @Description  批量创建多个翻译，支持两种格式：数组格式和前端对象格式
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        translations  body      domain.BatchTranslationRequest  true  "批量翻译请求"
// @Success      201           {object}  response.APIResponse
// @Failure      400           {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /translations/batch [post]
func (h *TranslationHandler) CreateBatch(ctx *gin.Context) {
	// 先尝试解析为前端格式（带有translations字段的对象格式）
	var batchReq domain.BatchTranslationRequest
	if err := ctx.ShouldBindJSON(&batchReq); err == nil && batchReq.Translations != nil {
		// 使用前端格式处理
		err := h.translationService.CreateBatchFromRequest(ctx.Request.Context(), batchReq)
		if err != nil {
			if err == domain.ErrProjectNotFound || err == domain.ErrLanguageNotFound {
				response.BadRequest(ctx, err.Error())
			} else {
				response.InternalServerError(ctx, "批量创建翻译失败")
			}
			return
		}
		response.Success(ctx, gin.H{"message": "批量创建成功"})
		return
	}

	// 如果前端格式解析失败，尝试数组格式
	var requests []domain.CreateTranslationRequest
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	err := h.translationService.CreateBatch(ctx.Request.Context(), requests)
	if err != nil {
		if err == domain.ErrProjectNotFound || err == domain.ErrLanguageNotFound {
			response.BadRequest(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "批量创建翻译失败")
		}
		return
	}

	response.Success(ctx, gin.H{"message": "批量创建成功"})
}

// GetByProjectID 根据项目ID获取翻译
// @Summary      获取项目翻译
// @Description  根据项目ID获取翻译列表
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int  true   "项目ID"
// @Param        page        query     int  false  "页码"  default(1)
// @Param        page_size   query     int  false  "每页数量"  default(10)
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/by-project/{project_id} [get]
func (h *TranslationHandler) GetByProjectID(ctx *gin.Context) {
	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	translations, total, err := h.translationService.GetByProjectID(ctx.Request.Context(), uint(projectID), pageSize, offset)
	if err != nil {
		if err == domain.ErrProjectNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取翻译列表失败"})
		}
		return
	}

	meta := &response.Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}

	response.SuccessWithMeta(ctx, translations, meta)
}

// GetMatrix 获取翻译矩阵
// @Summary      获取翻译矩阵
// @Description  获取项目的翻译矩阵（键-语言映射），支持分页
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int     true   "项目ID"
// @Param        page        query     int     false  "页码"  default(1)
// @Param        page_size   query     int     false  "每页数量"  default(10)
// @Param        keyword     query     string  false  "搜索关键词"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/matrix/by-project/{project_id} [get]
func (h *TranslationHandler) GetMatrix(ctx *gin.Context) {
	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的项目ID")
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

	matrix, total, err := h.translationService.GetMatrix(ctx.Request.Context(), uint(projectID), pageSize, offset, keyword)
	if err != nil {
		if err == domain.ErrProjectNotFound {
			response.NotFound(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "获取翻译矩阵失败")
		}
		return
	}

	meta := &response.Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}

	response.SuccessWithMeta(ctx, matrix, meta)
}

// GetByID 根据ID获取翻译
// @Summary      获取翻译详情
// @Description  根据翻译ID获取翻译详细信息
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "翻译ID"
// @Success      200  {object}  domain.Translation
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/{id} [get]
func (h *TranslationHandler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的翻译ID"})
		return
	}

	translation, err := h.translationService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrTranslationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取翻译失败"})
		}
		return
	}

	response.Success(ctx, translation)
}

// Update 更新翻译
// @Summary      更新翻译
// @Description  更新翻译信息
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        id           path      int                               true  "翻译ID"
// @Param        translation  body      domain.CreateTranslationRequest  true  "翻译信息"
// @Success      200          {object}  domain.Translation
// @Failure      400          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/{id} [put]
func (h *TranslationHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的翻译ID"})
		return
	}

	var req domain.CreateTranslationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translation, err := h.translationService.Update(ctx.Request.Context(), uint(id), req)
	if err != nil {
		if err == domain.ErrTranslationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err == domain.ErrProjectNotFound || err == domain.ErrLanguageNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新翻译失败"})
		}
		return
	}

	response.Success(ctx, translation)
}

// Delete 删除翻译
// @Summary      删除翻译
// @Description  删除指定的翻译
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "翻译ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/{id} [delete]
func (h *TranslationHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的翻译ID"})
		return
	}

	err = h.translationService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrTranslationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除翻译失败"})
		}
		return
	}

	response.NoContent(ctx)
}

// DeleteBatch 批量删除翻译
// @Summary      批量删除翻译
// @Description  批量删除多个翻译
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        ids  body      []uint  true  "翻译ID列表"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/batch-delete [post]
func (h *TranslationHandler) DeleteBatch(ctx *gin.Context) {
	var ids []uint

	if err := ctx.ShouldBindJSON(&ids); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.translationService.DeleteBatch(ctx.Request.Context(), ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "批量删除翻译失败"})
		return
	}

	response.NoContent(ctx)
}

// Export 导出翻译
// @Summary      导出翻译
// @Description  导出项目翻译为指定格式
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int     true   "项目ID"
// @Param        format      query     string  false  "导出格式"  default("json")
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Security     BearerAuth
// @Router       /exports/project/{project_id} [get]
func (h *TranslationHandler) Export(ctx *gin.Context) {
	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	format := ctx.DefaultQuery("format", "json")

	data, err := h.translationService.Export(ctx.Request.Context(), uint(projectID), format)
	if err != nil {
		if err == domain.ErrProjectNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "导出翻译失败"})
		}
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Disposition", "attachment; filename=translations.json")
	ctx.Data(http.StatusOK, "application/json", data)
}

// Import 导入翻译
// @Summary      导入翻译
// @Description  导入项目翻译数据
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int                                       true  "项目ID"
// @Param        data        body      map[string]map[string]string             true  "翻译数据，格式为 {\"key1\": {\"en\": \"value1\", \"zh\": \"值1\"}}"
// @Param        format      query     string                                   false "导入格式" default("json")
// @Success      200         {object}  response.APIResponse
// @Failure      400         {object}  response.APIResponse
// @Failure      404         {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /imports/project/{project_id} [post]
func (h *TranslationHandler) Import(ctx *gin.Context) {
	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的项目ID")
		return
	}

	format := ctx.DefaultQuery("format", "json")

	// 读取请求体
	data, err := ctx.GetRawData()
	if err != nil {
		response.BadRequest(ctx, "读取请求数据失败")
		return
	}

	err = h.translationService.Import(ctx.Request.Context(), uint(projectID), data, format)
	if err != nil {
		if err == domain.ErrProjectNotFound {
			response.NotFound(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "导入翻译失败: "+err.Error())
		}
		return
	}

	response.Success(ctx, gin.H{"message": "导入翻译成功"})
}
