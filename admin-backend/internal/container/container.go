package container

import (
	"context"

	"go.uber.org/fx"

	"i18n-flow/internal/api/routes"
	"i18n-flow/internal/config"
	"i18n-flow/internal/di"
	"i18n-flow/utils"
)

// Container 依赖注入容器（基于 uber-go/fx）
type Container struct {
	*fx.App
	Config *config.Config
}

// NewContainer 创建容器
func NewContainer(cfg *config.Config, onStart ...func(*routes.Router)) *Container {
	var router *routes.Router

	app := fx.New(
		fx.NopLogger, // 禁用默认日志，使用应用自己的日志
		di.AppModule,
		fx.Invoke(func(lifecycle fx.Lifecycle, r *routes.Router) {
			router = r
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					utils.AppInfo("Router initialized via fx lifecycle")
					return nil
				},
			})
		}),
	)

	// 调用可选的启动回调
	for _, fn := range onStart {
		fn(router)
	}

	return &Container{
		App:   app,
		Config: cfg,
	}
}

// MustNewContainer 创建容器（失败 panic）
func MustNewContainer(cfg *config.Config, onStart ...func(*routes.Router)) *Container {
	c := NewContainer(cfg, onStart...)
	if err := c.Err(); err != nil {
		panic(err)
	}
	return c
}

// Start 启动容器
func (c *Container) Start(ctx context.Context) error {
	return c.App.Start(ctx)
}

// Stop 停止容器
func (c *Container) Stop(ctx context.Context) error {
	return c.App.Stop(ctx)
}

// GetConfig 获取配置
func (c *Container) GetConfig() *config.Config {
	return c.Config
}

// AppModule 导出 fx 模块供外部使用
var AppModule = di.AppModule
