package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 运行测试
	exitCode := m.Run()

	// 退出测试
	os.Exit(exitCode)
}
