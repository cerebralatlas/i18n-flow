package repository

import (
	"fmt"
	"i18n-flow/internal/config"
	"i18n-flow/internal/domain"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// 优化DSN配置，添加连接参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s&interpolateParams=true",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName)

	// GORM配置优化
	gormConfig := &gorm.Config{
		// 禁用默认事务以提高性能
		SkipDefaultTransaction: true,
		// 批量插入优化
		CreateBatchSize: 1000,
		// 准备语句缓存
		PrepareStmt: true,
	}

	// 在生产环境中禁用详细日志
	if os.Getenv("GO_ENV") == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 获取底层的sql.DB对象进行连接池优化
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 连接池优化配置
	sqlDB.SetMaxIdleConns(10)                   // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                  // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)         // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 连接最大空闲时间

	// 自动迁移表结构
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Project{},
		&domain.Language{},
		&domain.Translation{},
	)
	if err != nil {
		return nil, fmt.Errorf("自动迁移表结构失败: %w", err)
	}

	// 创建额外的性能优化索引
	if err := createOptimizationIndexes(db); err != nil {
		log.Printf("创建优化索引时出现警告: %v", err)
	}

	// 初始化种子数据
	if err := initSeedData(db); err != nil {
		return nil, fmt.Errorf("初始化种子数据失败: %w", err)
	}

	return db, nil
}

// initSeedData 初始化种子数据
func initSeedData(db *gorm.DB) error {
	// 创建管理员用户
	if err := createAdminUser(db); err != nil {
		return err
	}

	// 创建常见语言
	if err := createDefaultLanguages(db); err != nil {
		return err
	}

	return nil
}

// createAdminUser 创建默认管理员用户
func createAdminUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			adminPassword = "admin123" // 默认密码
		}

		password, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("生成密码哈希失败: %w", err)
		}

		adminUsername := os.Getenv("ADMIN_USERNAME")
		if adminUsername == "" {
			adminUsername = "admin"
		}

		admin := domain.User{
			Username: adminUsername,
			Password: string(password),
		}

		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("创建管理员用户失败: %w", err)
		}

		log.Printf("已创建默认管理员用户: %s", adminUsername)
	} else {
		log.Println("管理员用户已存在，无需创建")
	}

	return nil
}

// createDefaultLanguages 创建默认语言
func createDefaultLanguages(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.Language{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		// 定义常见语言列表
		languages := []domain.Language{
			{Code: "en", Name: "English", IsDefault: true},
			{Code: "zh_CN", Name: "简体中文", IsDefault: false},
			{Code: "zh_TW", Name: "繁體中文", IsDefault: false},
			{Code: "ja", Name: "日本語", IsDefault: false},
			{Code: "ko", Name: "한국어", IsDefault: false},
			{Code: "fr", Name: "Français", IsDefault: false},
			{Code: "de", Name: "Deutsch", IsDefault: false},
			{Code: "es", Name: "Español", IsDefault: false},
			{Code: "it", Name: "Italiano", IsDefault: false},
			{Code: "ru", Name: "Русский", IsDefault: false},
			{Code: "pt", Name: "Português", IsDefault: false},
			{Code: "nl", Name: "Nederlands", IsDefault: false},
			{Code: "ar", Name: "العربية", IsDefault: false},
			{Code: "hi", Name: "हिन्दी", IsDefault: false},
			{Code: "tr", Name: "Türkçe", IsDefault: false},
			{Code: "pl", Name: "Polski", IsDefault: false},
			{Code: "vi", Name: "Tiếng Việt", IsDefault: false},
			{Code: "th", Name: "ไทย", IsDefault: false},
			{Code: "id", Name: "Bahasa Indonesia", IsDefault: false},
			{Code: "sv", Name: "Svenska", IsDefault: false},
		}

		if err := db.CreateInBatches(languages, len(languages)).Error; err != nil {
			return fmt.Errorf("创建默认语言失败: %w", err)
		}

		log.Println("已创建默认语言列表")
	} else {
		log.Println("语言列表已存在，无需创建")
	}

	return nil
}

// createOptimizationIndexes 创建额外的性能优化索引
func createOptimizationIndexes(db *gorm.DB) error {
	// 这些索引将在模型标签中自动创建，但我们可以添加一些复合索引
	indexes := []string{
		// 为常用查询组合创建复合索引
		"CREATE INDEX IF NOT EXISTS idx_translations_project_status ON translations(project_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_translations_search_key ON translations(project_id, key_name, status)",
		// 对于TEXT字段，指定索引长度
		"CREATE INDEX IF NOT EXISTS idx_translations_search_value ON translations(project_id, value(191), status)",
		"CREATE INDEX IF NOT EXISTS idx_projects_status_name ON projects(status, name)",
		// 为翻译值创建前缀索引，用于搜索
		"CREATE INDEX IF NOT EXISTS idx_translations_value_prefix ON translations(value(191))",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			// 记录警告但不中断启动过程，因为索引可能已存在
			log.Printf("创建索引警告: %v, SQL: %s", err, indexSQL)
		}
	}

	return nil
}
