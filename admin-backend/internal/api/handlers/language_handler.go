package handlers

import (
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LanguageHandler 语言处理器
type LanguageHandler struct {
	languageService domain.LanguageService
}

// NewLanguageHandler 创建语言处理器
func NewLanguageHandler(languageService domain.LanguageService) *LanguageHandler {
	return &LanguageHandler{
		languageService: languageService,
	}
}

// Create 创建语言
// @Summary      创建语言
// @Description  创建新的语言
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Param        language  body      domain.CreateLanguageRequest  true  "语言信息"
// @Success      201       {object}  domain.Language
// @Failure      400       {object}  map[string]string
// @Failure      409       {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages [post]
func (h *LanguageHandler) Create(ctx *gin.Context) {
	var req domain.CreateLanguageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	language, err := h.languageService.Create(ctx.Request.Context(), req)
	if err != nil {
		if err == domain.ErrLanguageExists {
			response.Conflict(ctx, err.Error())
		} else if err == domain.ErrInvalidLanguage {
			response.ValidationError(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "创建语言失败")
		}
		return
	}

	response.Created(ctx, language)
}

// GetAll 获取所有语言
// @Summary      获取语言列表
// @Description  获取所有语言列表
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Language
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages [get]
func (h *LanguageHandler) GetAll(ctx *gin.Context) {
	languages, err := h.languageService.GetAll(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "获取语言列表失败")
		return
	}

	response.Success(ctx, languages)
}

// Update 更新语言
// @Summary      更新语言
// @Description  更新语言信息
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Param        id        path      int                            true  "语言ID"
// @Param        language  body      domain.CreateLanguageRequest  true  "语言信息"
// @Success      200       {object}  domain.Language
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages/{id} [put]
func (h *LanguageHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的语言ID")
		return
	}

	var req domain.CreateLanguageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	language, err := h.languageService.Update(ctx.Request.Context(), uint(id), req)
	if err != nil {
		if err == domain.ErrLanguageNotFound {
			response.NotFound(ctx, err.Error())
		} else if err == domain.ErrLanguageExists || err == domain.ErrInvalidInput {
			response.ValidationError(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "更新语言失败")
		}
		return
	}

	response.Success(ctx, language)
}

// Delete 删除语言
// @Summary      删除语言
// @Description  删除指定的语言
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "语言ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages/{id} [delete]
func (h *LanguageHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的语言ID")
		return
	}

	err = h.languageService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrLanguageNotFound {
			response.NotFound(ctx, err.Error())
		} else if err == domain.ErrInvalidInput {
			response.ValidationError(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, "删除语言失败")
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
