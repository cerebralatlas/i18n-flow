package domain

import "context"

// DTOs - 数据传输对象
type (
	LoginRequest struct {
		Username string `json:"username" binding:"required" example:"admin"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	LoginResponse struct {
		Token        string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		User         User   `json:"user"`
	}

	RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	}

	CreateProjectRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	UpdateProjectRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	CreateLanguageRequest struct {
		Code      string `json:"code" binding:"required"`
		Name      string `json:"name" binding:"required"`
		IsDefault bool   `json:"is_default"`
	}

	CreateTranslationRequest struct {
		ProjectID  uint   `json:"project_id" binding:"required"`
		KeyName    string `json:"key_name" binding:"required"`
		Context    string `json:"context"`
		LanguageID uint   `json:"language_id" binding:"required"`
		Value      string `json:"value" binding:"required"`
	}

	DashboardStats struct {
		TotalProjects     int `json:"total_projects"`
		TotalLanguages    int `json:"total_languages"`
		TotalTranslations int `json:"total_translations"`
		TotalKeys         int `json:"total_keys"`
	}
)

// UserService 用户服务接口
type UserService interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req RefreshRequest) (*LoginResponse, error)
	GetUserInfo(ctx context.Context, userID uint) (*User, error)
}

// ProjectService 项目服务接口
type ProjectService interface {
	Create(ctx context.Context, req CreateProjectRequest) (*Project, error)
	GetByID(ctx context.Context, id uint) (*Project, error)
	GetAll(ctx context.Context, limit, offset int) ([]*Project, int64, error)
	Update(ctx context.Context, id uint, req UpdateProjectRequest) (*Project, error)
	Delete(ctx context.Context, id uint) error
}

// LanguageService 语言服务接口
type LanguageService interface {
	Create(ctx context.Context, req CreateLanguageRequest) (*Language, error)
	GetByID(ctx context.Context, id uint) (*Language, error)
	GetAll(ctx context.Context) ([]*Language, error)
	Update(ctx context.Context, id uint, req CreateLanguageRequest) (*Language, error)
	Delete(ctx context.Context, id uint) error
}

// TranslationService 翻译服务接口
type TranslationService interface {
	Create(ctx context.Context, req CreateTranslationRequest) (*Translation, error)
	CreateBatch(ctx context.Context, translations []CreateTranslationRequest) error
	GetByID(ctx context.Context, id uint) (*Translation, error)
	GetByProjectID(ctx context.Context, projectID uint, limit, offset int) ([]*Translation, int64, error)
	GetMatrix(ctx context.Context, projectID uint, limit, offset int, keyword string) (map[string]map[string]string, int64, error)
	Update(ctx context.Context, id uint, req CreateTranslationRequest) (*Translation, error)
	Delete(ctx context.Context, id uint) error
	DeleteBatch(ctx context.Context, ids []uint) error
	Export(ctx context.Context, projectID uint, format string) ([]byte, error)
	Import(ctx context.Context, projectID uint, data []byte, format string) error
}

// DashboardService 仪表板服务接口
type DashboardService interface {
	GetStats(ctx context.Context) (*DashboardStats, error)
}

// AuthService 认证服务接口
type AuthService interface {
	GenerateToken(user *User) (string, error)
	GenerateRefreshToken(user *User) (string, error)
	ValidateToken(token string) (*User, error)
	ValidateRefreshToken(token string) (*User, error)
}
