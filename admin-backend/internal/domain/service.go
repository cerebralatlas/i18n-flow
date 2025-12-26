package domain

import (
	"context"
)

// UserService 用户服务接口
type UserService interface {
	Login(ctx context.Context, params LoginParams) (*LoginResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error)
	GetUserInfo(ctx context.Context, userID uint64) (*User, error)

	// 用户管理
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	GetAllUsers(ctx context.Context, limit, offset int, keyword string) ([]*User, int64, error)
	GetUserByID(ctx context.Context, id uint64) (*User, error)
	UpdateUser(ctx context.Context, id uint64, params UpdateUserParams) (*User, error)
	ChangePassword(ctx context.Context, userID uint64, params ChangePasswordParams) error
	ResetPassword(ctx context.Context, userID uint64, newPassword string) error
	DeleteUser(ctx context.Context, id uint64) error
}

// ProjectService 项目服务接口
type ProjectService interface {
	Create(ctx context.Context, params CreateProjectParams, userID uint64) (*Project, error)
	GetByID(ctx context.Context, id uint64) (*Project, error)
	GetAll(ctx context.Context, limit, offset int, keyword string) ([]*Project, int64, error)
	GetAccessibleProjects(ctx context.Context, userID uint64, limit, offset int, keyword string) ([]*Project, int64, error)
	Update(ctx context.Context, id uint64, params UpdateProjectParams, userID uint64) (*Project, error)
	Delete(ctx context.Context, id uint64) error
}

// LanguageService 语言服务接口
type LanguageService interface {
	Create(ctx context.Context, params CreateLanguageParams, userID uint64) (*Language, error)
	GetByID(ctx context.Context, id uint64) (*Language, error)
	GetAll(ctx context.Context) ([]*Language, error)
	Update(ctx context.Context, id uint64, params CreateLanguageParams, userID uint64) (*Language, error)
	Delete(ctx context.Context, id uint64) error
}

// TranslationService 翻译服务接口
type TranslationService interface {
	Create(ctx context.Context, input TranslationInput, userID uint64) (*Translation, error)
	CreateBatch(ctx context.Context, inputs []TranslationInput) error
	CreateBatchFromRequest(ctx context.Context, params BatchTranslationParams) error
	UpsertBatch(ctx context.Context, inputs []TranslationInput) error
	GetByID(ctx context.Context, id uint64) (*Translation, error)
	GetByProjectID(ctx context.Context, projectID uint64, limit, offset int) ([]*Translation, int64, error)
	GetMatrix(ctx context.Context, projectID uint64, limit, offset int, keyword string) (map[string]map[string]string, int64, error)
	Update(ctx context.Context, id uint64, input TranslationInput, userID uint64) (*Translation, error)
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
	GenerateToken(ctx context.Context, user *User) (string, error)
	GenerateRefreshToken(ctx context.Context, user *User) (string, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
	ValidateRefreshToken(ctx context.Context, token string) (*User, error)
}

// ProjectMemberService 项目成员服务接口
type ProjectMemberService interface {
	AddMember(ctx context.Context, projectID uint64, params AddMemberParams, createdBy uint64) (*ProjectMember, error)
	GetProjectMembers(ctx context.Context, projectID uint64) ([]*ProjectMemberInfo, error)
	GetUserProjects(ctx context.Context, userID uint64) ([]*Project, error)
	UpdateMemberRole(ctx context.Context, projectID, userID uint64, params UpdateMemberRoleParams) (*ProjectMember, error)
	RemoveMember(ctx context.Context, projectID, userID uint64) error
	CheckPermission(ctx context.Context, userID, projectID uint64, requiredRole string) (bool, error)
	GetMemberRole(ctx context.Context, userID, projectID uint64) (string, error)
}
