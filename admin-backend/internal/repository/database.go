package repository

import (
	"fmt"
	"i18n-flow/internal/config"
	"i18n-flow/internal/domain"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

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
