package tests

import (
	"os"
	"testing"

	"i18n-flow/utils"
)

func TestMain(m *testing.M) {
	// 设置测试环境
	setup()

	// 运行测试
	exitCode := m.Run()

	// 清理测试环境
	teardown()

	// 退出测试
	os.Exit(exitCode)
}

func setup() {
	// 初始化测试环境的日志配置
	logConfig := utils.GetDefaultMultiLogConfig("test")
	utils.InitMultiLogger(logConfig)

	utils.AppInfo("测试环境初始化完成")
}

func teardown() {
	// 同步日志，确保所有日志都已写入
	utils.SyncAll()

	utils.AppInfo("测试环境清理完成")
}
