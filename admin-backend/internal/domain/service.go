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
		ProjectID  uint64 `json:"project_id" binding:"required"`
		KeyName    string `json:"key_name" binding:"required"`
		Context    string `json:"context"`
		LanguageID uint64 `json:"language_id" binding:"required"`
		Value      string `json:"value" binding:"required"`
	}

	// BatchTranslationRequest 批量翻译请求（前端格式）
	BatchTranslationRequest struct {
		ProjectID    uint64            `json:"project_id" binding:"required"`
		KeyName      string            `json:"key_name" binding:"required"`
		Context      string            `json:"context"`
		Translations map[string]string `json:"translations" binding:"required"`
	}

	DashboardStats struct {
		TotalProjects     int `json:"total_projects"`
		TotalLanguages    int `json:"total_languages"`
		TotalTranslations int `json:"total_translations"`
		TotalKeys         int `json:"total_keys"`
	}

	// 用户管理相关DTOs
	CreateUserRequest struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=admin member viewer"`
	}

	UpdateUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email" binding:"omitempty,email"`
		Role     string `json:"role" binding:"omitempty,oneof=admin member viewer"`
		Status   string `json:"status" binding:"omitempty,oneof=active disabled"`
	}

	ChangePasswordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	ResetPasswordRequest struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	// 项目成员管理相关DTOs
	AddProjectMemberRequest struct {
		UserID uint64 `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required,oneof=owner editor viewer"`
	}

	UpdateProjectMemberRequest struct {
		Role string `json:"role" binding:"required,oneof=owner editor viewer"`
	}

	ProjectMemberInfo struct {
		ID       uint64 `json:"id"`
		UserID   uint64 `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
)

// UserService 用户服务接口
type UserService interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req RefreshRequest) (*LoginResponse, error)
	GetUserInfo(ctx context.Context, userID uint64) (*User, error)

	// 用户管理
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	GetAllUsers(ctx context.Context, limit, offset int, keyword string) ([]*User, int64, error)
	GetUserByID(ctx context.Context, id uint64) (*User, error)
	UpdateUser(ctx context.Context, id uint64, req UpdateUserRequest) (*User, error)
	ChangePassword(ctx context.Context, userID uint64, req ChangePasswordRequest) error
	ResetPassword(ctx context.Context, userID uint64, req ResetPasswordRequest) error
	DeleteUser(ctx context.Context, id uint64) error
}

// ProjectService 项目服务接口
type ProjectService interface {
	Create(ctx context.Context, req CreateProjectRequest, userID uint64) (*Project, error)
	GetByID(ctx context.Context, id uint64) (*Project, error)
	GetAll(ctx context.Context, limit, offset int, keyword string) ([]*Project, int64, error)
	GetAccessibleProjects(ctx context.Context, userID uint64, limit, offset int, keyword string) ([]*Project, int64, error)
	Update(ctx context.Context, id uint64, req UpdateProjectRequest, userID uint64) (*Project, error)
	Delete(ctx context.Context, id uint64) error
}

// LanguageService 语言服务接口
type LanguageService interface {
	Create(ctx context.Context, req CreateLanguageRequest, userID uint64) (*Language, error)
	GetByID(ctx context.Context, id uint64) (*Language, error)
	GetAll(ctx context.Context) ([]*Language, error)
	Update(ctx context.Context, id uint64, req CreateLanguageRequest, userID uint64) (*Language, error)
	Delete(ctx context.Context, id uint64) error
}

// TranslationService 翻译服务接口
type TranslationService interface {
	Create(ctx context.Context, req CreateTranslationRequest, userID uint64) (*Translation, error)
	CreateBatch(ctx context.Context, translations []CreateTranslationRequest) error
	CreateBatchFromRequest(ctx context.Context, req BatchTranslationRequest) error
	UpsertBatch(ctx context.Context, translations []CreateTranslationRequest) error
	GetByID(ctx context.Context, id uint64) (*Translation, error)
	GetByProjectID(ctx context.Context, projectID uint64, limit, offset int) ([]*Translation, int64, error)
	GetMatrix(ctx context.Context, projectID uint64, limit, offset int, keyword string) (map[string]map[string]string, int64, error)
	Update(ctx context.Context, id uint64, req CreateTranslationRequest, userID uint64) (*Translation, error)
	Delete(ctx context.Context, id uint64) error
	DeleteBatch(ctx context.Context, ids []uint64) error
	Export(ctx context.Context, projectID uint64, format string) ([]byte, error)
	Import(ctx context.Context, projectID uint64, data []byte, format string) error
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

// ProjectMemberService 项目成员服务接口
type ProjectMemberService interface {
	AddMember(ctx context.Context, projectID uint64, req AddProjectMemberRequest, userID uint64) (*ProjectMember, error)
	GetProjectMembers(ctx context.Context, projectID uint64) ([]*ProjectMemberInfo, error)
	GetUserProjects(ctx context.Context, userID uint64) ([]*Project, error)
	UpdateMemberRole(ctx context.Context, projectID, userID uint64, req UpdateProjectMemberRequest) (*ProjectMember, error)
	RemoveMember(ctx context.Context, projectID, userID uint64) error
	CheckPermission(ctx context.Context, userID, projectID uint64, requiredRole string) (bool, error)
	GetMemberRole(ctx context.Context, userID, projectID uint64) (string, error)
}
