package handlers

import "i18n-flow/internal/container"

// HandlerFactory 处理器工厂
type HandlerFactory struct {
	container *container.Container
}

// NewHandlerFactory 创建处理器工厂
func NewHandlerFactory(container *container.Container) *HandlerFactory {
	return &HandlerFactory{
		container: container,
	}
}

// UserHandler 获取用户处理器
func (f *HandlerFactory) UserHandler() *UserHandler {
	return NewUserHandler(f.container.UserService())
}

// ProjectHandler 获取项目处理器
func (f *HandlerFactory) ProjectHandler() *ProjectHandler {
	return NewProjectHandler(f.container.ProjectService())
}

// LanguageHandler 获取语言处理器
func (f *HandlerFactory) LanguageHandler() *LanguageHandler {
	return NewLanguageHandler(f.container.LanguageService())
}

// TranslationHandler 获取翻译处理器
func (f *HandlerFactory) TranslationHandler() *TranslationHandler {
	return NewTranslationHandler(f.container.TranslationService())
}

// DashboardHandler 获取仪表板处理器
func (f *HandlerFactory) DashboardHandler() *DashboardHandler {
	return NewDashboardHandler(f.container.DashboardService())
}
