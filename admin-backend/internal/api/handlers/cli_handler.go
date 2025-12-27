package handlers

import (
	"i18n-flow/internal/api/response"
	"i18n-flow/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CLIHandler CLI处理器
type CLIHandler struct {
	translationService domain.TranslationService
	projectService     domain.ProjectService
	languageService    domain.LanguageService
}

// NewCLIHandler 创建CLI处理器
func NewCLIHandler(
	translationService domain.TranslationService,
	projectService domain.ProjectService,
	languageService domain.LanguageService,
) *CLIHandler {
	return &CLIHandler{
		translationService: translationService,
		projectService:     projectService,
		languageService:    languageService,
	}
}

// Auth CLI身份验证
// @Summary      CLI身份验证
// @Description  验证CLI API Key
// @Tags         CLI
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.APIResponse
// @Failure      401  {object}  response.APIResponse
// @Security     ApiKeyAuth
// @Router       /cli/auth [get]
func (h *CLIHandler) Auth(ctx *gin.Context) {
	// API Key认证由中间件处理，能到这里说明认证成功
	response.Success(ctx, gin.H{
		"status":  "ok",
		"message": "CLI authentication successful",
	})
}

// GetTranslations 获取翻译数据
// @Summary      获取翻译数据
// @Description  获取项目翻译数据供CLI使用
// @Tags         CLI
// @Accept       json
// @Produce      json
// @Param        project_id  query     string  false  "项目ID"
// @Param        locale      query     string  false  "语言代码"
// @Success      200         {object}  response.APIResponse
// @Failure      400         {object}  response.APIResponse
// @Failure      404         {object}  response.APIResponse
// @Security     ApiKeyAuth
// @Router       /cli/translations [get]
func (h *CLIHandler) GetTranslations(ctx *gin.Context) {
	projectIDStr := ctx.Query("project_id")
	locale := ctx.Query("locale")

	// 如果没有指定项目ID，返回错误
	if projectIDStr == "" {
		response.BadRequest(ctx, "project_id is required")
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "invalid project_id")
		return
	}

	// 验证项目是否存在
	_, err = h.projectService.GetByID(ctx.Request.Context(), projectID)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			response.NotFound(ctx, err.Error())
		default:
			response.InternalServerError(ctx, "获取项目失败")
		}
		return
	}

	// 获取翻译矩阵数据（不分页，获取所有数据）
	matrix, _, err := h.translationService.GetMatrix(ctx.Request.Context(), projectID, -1, 0, "")
	if err != nil {
		response.InternalServerError(ctx, "获取翻译数据失败")
		return
	}

	// 转换为简单格式 (key -> language -> value)
	simpleMatrix := make(map[string]map[string]string)
	for key, langs := range matrix {
		simpleMatrix[key] = make(map[string]string)
		for lang, cell := range langs {
			simpleMatrix[key][lang] = cell.Value
		}
	}

	// 如果指定了locale，只返回该语言的数据
	if locale != "" {
		filteredMatrix := make(map[string]map[string]string)
		for key, translations := range simpleMatrix {
			if value, exists := translations[locale]; exists {
				filteredMatrix[key] = map[string]string{locale: value}
			}
		}
		response.Success(ctx, filteredMatrix)
		return
	}

	// 返回完整的翻译矩阵
	response.Success(ctx, simpleMatrix)
}

// PushKeysRequest 推送键请求
type PushKeysRequest struct {
	ProjectID    string                       `json:"project_id" binding:"required"`
	Keys         []string                     `json:"keys" binding:"required"`
	Defaults     map[string]string            `json:"defaults"`     // 保持向后兼容（已废弃）
	Translations map[string]map[string]string `json:"translations"` // 新增：语言代码 -> 键值对映射
}

// PushKeysResponse 推送键响应
type PushKeysResponse struct {
	Added   []string `json:"added"`
	Existed []string `json:"existed"`
	Failed  []string `json:"failed"`
}

// PushKeys 推送翻译键
// @Summary      推送翻译键
// @Description  从CLI推送新的翻译键
// @Tags         CLI
// @Accept       json
// @Produce      json
// @Param        request  body      PushKeysRequest  true  "推送键请求"
// @Success      200      {object}  response.APIResponse
// @Failure      400      {object}  response.APIResponse
// @Failure      404      {object}  response.APIResponse
// @Security     ApiKeyAuth
// @Router       /cli/keys [post]
func (h *CLIHandler) PushKeys(ctx *gin.Context) {
	var req PushKeysRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err.Error())
		return
	}

	projectID, err := strconv.ParseUint(req.ProjectID, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "invalid project_id")
		return
	}

	// 验证项目是否存在
	_, err = h.projectService.GetByID(ctx.Request.Context(), projectID)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			response.NotFound(ctx, err.Error())
		default:
			response.InternalServerError(ctx, "获取项目失败")
		}
		return
	}

	// 获取所有语言
	languages, err := h.languageService.GetAll(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "获取语言列表失败")
		return
	}

	// 找到默认语言（通常是第一个或标记为默认的语言）
	var defaultLanguage *domain.Language
	for _, lang := range languages {
		if lang.IsDefault {
			defaultLanguage = lang
			break
		}
	}
	if defaultLanguage == nil && len(languages) > 0 {
		defaultLanguage = languages[0] // 如果没有默认语言，使用第一个语言
	}

	if defaultLanguage == nil {
		response.BadRequest(ctx, "no languages available in project")
		return
	}

	// 获取现有的翻译键
	matrix, _, err := h.translationService.GetMatrix(ctx.Request.Context(), projectID, -1, 0, "")
	if err != nil {
		response.InternalServerError(ctx, "获取现有翻译失败")
		return
	}

	var added []string
	var existed []string
	var failed []string

	// 处理每个键
	for _, key := range req.Keys {
		if _, exists := matrix[key]; exists {
			existed = append(existed, key)
			continue
		}

		// 为所有语言创建新的翻译记录
		keyAdded := false
		keyFailed := false

		for _, language := range languages {
			// 确定翻译值
			var value string

			// 优先使用新的多语言数据结构
			if req.Translations != nil {
				if langTranslations, exists := req.Translations[language.Code]; exists {
					value = langTranslations[key]
				}
			} else {
				// 向后兼容：使用旧的 Defaults 字段
				if language.Code == defaultLanguage.Code {
					value = req.Defaults[key]
				}
			}

			// DTO -> Domain params
			input := domain.TranslationInput{
				ProjectID:  projectID,
				KeyName:    key,
				LanguageID: language.ID,
				Value:      value,
			}

			_, err := h.translationService.Create(ctx.Request.Context(), input, 1) // 使用系统管理员ID
			if err != nil {
				keyFailed = true
			} else if !keyAdded {
				keyAdded = true
			}
		}

		// 记录结果
		if keyFailed && !keyAdded {
			failed = append(failed, key)
		} else if keyAdded {
			added = append(added, key)
		}
	}

	result := PushKeysResponse{
		Added:   added,
		Existed: existed,
		Failed:  failed,
	}

	response.Success(ctx, result)
}
