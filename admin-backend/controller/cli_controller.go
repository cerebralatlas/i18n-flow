package controller

import (
	"i18n-flow/config"
	"i18n-flow/errors"
	"i18n-flow/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CLIController 处理CLI工具相关的请求
type CLIController struct {
	translationService service.TranslationService
	projectService     service.ProjectService
}

// NewCLIController 创建一个新的CLI控制器
func NewCLIController() *CLIController {
	return &CLIController{
		translationService: service.NewTranslationService(),
		projectService:     service.NewProjectService(),
	}
}

// GetAllTranslations 获取所有翻译，支持按项目和语言筛选
// @Summary 获取所有翻译
// @Description 获取所有翻译，供CLI工具使用
// @Tags CLI
// @Accept json
// @Produce json
// @Param project_id query int false "项目ID"
// @Param locale query string false "语言代码"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Router /cli/translations [get]
func (cc *CLIController) GetAllTranslations(c *gin.Context) {
	projectID := c.Query("project_id")
	locale := c.Query("locale")

	// 调用 service 获取数据
	translations, err := cc.translationService.GetTranslationsForCLI(projectID, locale)
	if err != nil {
		errors.ErrorResponse(c, errors.NewDatabaseError(err).WithDetails("获取翻译数据失败"))
		return
	}

	errors.SuccessResponse(c, translations)
}

// PushKeys 推送新的翻译键
// @Summary 推送新的翻译键
// @Description 推送新的翻译键到服务器
// @Tags CLI
// @Accept json
// @Produce json
// @Param request body service.KeysPushRequest true "键值请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Router /cli/keys [post]
func (cc *CLIController) PushKeys(c *gin.Context) {
	var request service.KeysPushRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errors.ErrorResponse(c, errors.NewInvalidParamsError("无效的请求数据: "+err.Error()))
		return
	}

	result, err := cc.translationService.PushKeysFromCLI(request)
	if err != nil {
		errors.ErrorResponse(c, errors.NewDatabaseError(err).WithDetails("推送键值失败"))
		return
	}

	errors.SuccessResponse(c, result)
}

// CheckAPIKey 检查API Key是否有效
// @Summary 检查API Key是否有效
// @Description 用于CLI工具测试连接
// @Tags CLI
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} errors.Response
// @Router /cli/auth [get]
func (cc *CLIController) CheckAPIKey(c *gin.Context) {
	// 从请求头中获取API Key
	apiKey := c.GetHeader("X-API-Key")

	// 验证API Key是否有效
	if apiKey == "" || apiKey != config.GetConfig().CLI.APIKey {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Invalid or missing API key",
			"status": "error",
		})
		return
	}

	// API Key有效，返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"message": "API Key 验证成功",
		"status":  "ok",
	})
}
