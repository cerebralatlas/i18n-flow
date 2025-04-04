package controller

import (
	"i18n-flow/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LanguageController 语言控制器
type LanguageController struct {
	languageService service.LanguageService
}

// NewLanguageController 创建语言控制器
func NewLanguageController() *LanguageController {
	return &LanguageController{
		languageService: service.LanguageService{},
	}
}

// GetLanguages 获取所有语言
// @Summary      获取所有语言
// @Description  获取系统中所有可用的语言列表
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Success      200  {array}   service.LanguageResponse
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages [get]
func (c *LanguageController) GetLanguages(ctx *gin.Context) {
	languages, err := c.languageService.GetLanguages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, languages)
}

// CreateLanguage 创建语言
// @Summary      创建新语言
// @Description  在系统中添加一种新的语言
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Param        language  body      service.LanguageRequest  true  "语言信息"
// @Success      201       {object}  service.LanguageResponse
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages [post]
func (c *LanguageController) CreateLanguage(ctx *gin.Context) {
	var req service.LanguageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	language, err := c.languageService.CreateLanguage(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, language)
}

// UpdateLanguage 更新语言
// @Summary      更新语言
// @Description  更新系统中现有语言的信息
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Param        id        path      int                      true  "语言ID"
// @Param        language  body      service.LanguageRequest  true  "语言信息"
// @Success      200       {object}  service.LanguageResponse
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages/{id} [put]
func (c *LanguageController) UpdateLanguage(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的语言ID"})
		return
	}

	var req service.LanguageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	language, err := c.languageService.UpdateLanguage(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, language)
}

// DeleteLanguage 删除语言
// @Summary      删除语言
// @Description  从系统中删除指定的语言
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "语言ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /languages/{id} [delete]
func (c *LanguageController) DeleteLanguage(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的语言ID"})
		return
	}

	if err := c.languageService.DeleteLanguage(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "语言删除成功"})
}
