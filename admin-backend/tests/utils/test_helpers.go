package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"i18n-flow/internal/domain"
	"i18n-flow/utils"
)

// CreateTempDir 创建临时目录并在测试结束后删除
func CreateTempDir(t *testing.T, prefix string) string {
	tempDir, err := os.MkdirTemp("", prefix)
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	return tempDir
}

// SetupTestLogger 设置测试日志器并返回观察者
func SetupTestLogger(t *testing.T) (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	// 保存原始日志器
	originalLogger := utils.Logger
	originalSugar := utils.SugaredLogger

	// 设置测试日志器
	utils.Logger = logger
	utils.SugaredLogger = logger.Sugar()

	// 测试结束后恢复原始日志器
	t.Cleanup(func() {
		utils.Logger = originalLogger
		utils.SugaredLogger = originalSugar
	})

	return logger, logs
}

// SetupTestDB 创建内存数据库用于测试
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移模型
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Project{},
		&domain.Language{},
		&domain.Translation{},
	)
	require.NoError(t, err)

	return db
}

// CreateTestFile 创建测试文件并返回路径
func CreateTestFile(t *testing.T, dir, name, content string) string {
	if dir == "" {
		dir = t.TempDir()
	}

	filePath := filepath.Join(dir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	return filePath
}

// CreateTestUser 创建测试用户
func CreateTestUser(t *testing.T, db *gorm.DB, username, password string) *domain.User {
	user := &domain.User{
		Username: username,
		Password: password,
	}

	err := db.Create(user).Error
	require.NoError(t, err)

	return user
}

// CreateTestProject 创建测试项目
func CreateTestProject(t *testing.T, db *gorm.DB, name, description, slug string) *domain.Project {
	project := &domain.Project{
		Name:        name,
		Description: description,
		Slug:        slug,
		Status:      "active",
	}

	err := db.Create(project).Error
	require.NoError(t, err)

	return project
}

// CreateTestLanguage 创建测试语言
func CreateTestLanguage(t *testing.T, db *gorm.DB, code, name string, isDefault bool) *domain.Language {
	language := &domain.Language{
		Code:      code,
		Name:      name,
		IsDefault: isDefault,
		Status:    "active",
	}

	err := db.Create(language).Error
	require.NoError(t, err)

	return language
}

// CreateTestTranslation 创建测试翻译
func CreateTestTranslation(t *testing.T, db *gorm.DB, projectID uint, keyName, context string, languageID uint, value string) *domain.Translation {
	translation := &domain.Translation{
		ProjectID:  projectID,
		KeyName:    keyName,
		Context:    context,
		LanguageID: languageID,
		Value:      value,
		Status:     "active",
	}

	err := db.Create(translation).Error
	require.NoError(t, err)

	return translation
}

// AssertFileExists 断言文件存在
func AssertFileExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.NoError(t, err)
}

// AssertFileNotExists 断言文件不存在
func AssertFileNotExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.True(t, os.IsNotExist(err))
}

// ReadTestFile 读取测试文件内容
func ReadTestFile(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(content)
}
