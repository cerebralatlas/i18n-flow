package handlers

import (
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
// @Description  批量创建多个翻译
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        translations  body      []domain.CreateTranslationRequest  true  "翻译信息列表"
// @Success      201           {object}  map[string]string
// @Failure      400           {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/batch [post]
func (h *TranslationHandler) CreateBatch(ctx *gin.Context) {
	var requests []domain.CreateTranslationRequest

	if err := ctx.ShouldBindJSON(&requests); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.translationService.CreateBatch(ctx.Request.Context(), requests)
	if err != nil {
		if err == domain.ErrProjectNotFound || err == domain.ErrLanguageNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建翻译失败"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "批量创建成功"})
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

	response := gin.H{
		"data": translations,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total_count": total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMatrix 获取翻译矩阵
// @Summary      获取翻译矩阵
// @Description  获取项目的翻译矩阵（键-语言映射）
// @Tags         翻译管理
// @Accept       json
// @Produce      json
// @Param        project_id  path      int  true  "项目ID"
// @Success      200         {object}  map[string]map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Security     BearerAuth
// @Router       /translations/matrix/by-project/{project_id} [get]
func (h *TranslationHandler) GetMatrix(ctx *gin.Context) {
	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	matrix, err := h.translationService.GetMatrix(ctx.Request.Context(), uint(projectID))
	if err != nil {
		if err == domain.ErrProjectNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取翻译矩阵失败"})
		}
		return
	}

	ctx.JSON(http.StatusOK, matrix)
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

	ctx.JSON(http.StatusOK, translation)
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

	ctx.JSON(http.StatusOK, translation)
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

	ctx.Status(http.StatusNoContent)
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

	ctx.Status(http.StatusNoContent)
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
