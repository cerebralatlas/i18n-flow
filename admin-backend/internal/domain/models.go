package domain

import (
	"time"

	"gorm.io/gorm"
)

// User 用户领域模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"unique;size:50;not null" json:"username"`
	Email     string    `gorm:"unique;size:100" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	Role      string    `gorm:"size:20;default:member;index:idx_user_role" json:"role"`     // admin, member, viewer
	Status    string    `gorm:"size:20;default:active;index:idx_user_status" json:"status"` // active, disabled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Project 项目领域模型
type Project struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null;unique;index:idx_project_search" json:"name"` // 项目名称
	Description  string         `gorm:"size:500;index:idx_project_search" json:"description"`          // 项目描述
	Slug         string         `gorm:"size:100;not null;unique;index" json:"slug"`                    // 项目标识，用于URL
	Status       string         `gorm:"size:20;default:active;index:idx_project_status" json:"status"` // 项目状态：active, archived
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Translations []Translation  `gorm:"foreignKey:ProjectID" json:"-"` // 关联的翻译
}

// Language 语言领域模型
type Language struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:10;not null;unique" json:"code"`  // 语言代码，如 en, zh-CN
	Name      string         `gorm:"size:50;not null" json:"name"`         // 语言名称，如 English, 简体中文
	IsDefault bool           `gorm:"default:false" json:"is_default"`      // 是否为默认语言
	Status    string         `gorm:"size:20;default:active" json:"status"` // 状态：active, inactive
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Translation 翻译领域模型
type Translation struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ProjectID  uint           `gorm:"not null;index:idx_translation_project;uniqueIndex:idx_translation_unique,priority:1" json:"project_id"`    // 关联的项目ID
	KeyName    string         `gorm:"size:255;not null;index:idx_translation_key;uniqueIndex:idx_translation_unique,priority:2" json:"key_name"` // 翻译键名
	Context    string         `gorm:"size:500" json:"context"`                                                                                   // 上下文说明
	LanguageID uint           `gorm:"not null;index:idx_translation_language;uniqueIndex:idx_translation_unique,priority:3" json:"language_id"`  // 语言ID
	Value      string         `gorm:"type:text" json:"value"`                                                                                    // 翻译值
	Status     string         `gorm:"size:20;default:active;index:idx_translation_status" json:"status"`                                         // 状态：active, deprecated
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Project  Project  `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`   // 关联的项目
	Language Language `gorm:"foreignKey:LanguageID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"` // 关联的语言
}

// ProjectMember 项目成员关联模型
type ProjectMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProjectID uint           `gorm:"not null;index:idx_project_member;uniqueIndex:idx_project_member_unique,priority:1" json:"project_id"`
	UserID    uint           `gorm:"not null;index:idx_project_member;uniqueIndex:idx_project_member_unique,priority:2" json:"user_id"`
	Role      string         `gorm:"size:20;default:viewer;index:idx_project_member_role" json:"role"` // owner, editor, viewer
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Project Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
