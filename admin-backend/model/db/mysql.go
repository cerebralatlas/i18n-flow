package db

import (
	"fmt"
	"i18n-flow/config"
	"i18n-flow/model"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Language{},
		&model.Translation{},
	)
	if err != nil {
		log.Fatalf("自动迁移表结构失败: %v", err)
	}

	initSeedData()
}

// 初始化数据
func initSeedData() {
	// 创建管理员用户
	createAdminUser()
	// 创建常见语言
	createDefaultLanguages()
}

// createAdminUser 创建默认管理员用户
func createAdminUser() {
	var count int64
	DB.Model(&model.User{}).Count(&count)

	if count == 0 {
		password, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("生成密码哈希失败: %v", err)
			return
		}

		admin := model.User{
			Username: os.Getenv("ADMIN_USERNAME"),
			Password: string(password),
		}

		result := DB.Create(&admin)
		if result.Error != nil {
			log.Printf("创建管理员用户失败: %v", result.Error)
		} else {
			log.Println("已创建默认管理员用户")
		}
	} else {
		log.Println("管理员用户已存在，无需创建")
	}
}

// createDefaultLanguages 创建默认语言
func createDefaultLanguages() {
	var count int64
	DB.Model(&model.Language{}).Count(&count)

	if count == 0 {
		// 定义常见语言列表
		languages := []model.Language{
			{Code: "en", Name: "English", IsDefault: true},
			{Code: "zh-CN", Name: "简体中文", IsDefault: false},
			{Code: "zh-TW", Name: "繁體中文", IsDefault: false},
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

		result := DB.CreateInBatches(languages, len(languages))
		if result.Error != nil {
			log.Printf("创建默认语言失败: %v", result.Error)
		} else {
			log.Println("已创建默认语言列表")
		}
	} else {
		log.Println("语言列表已存在，无需创建")
	}
}
