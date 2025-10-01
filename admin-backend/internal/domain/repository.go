package domain

import "context"

// UserRepository 用户数据访问接口
type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context, limit, offset int, keyword string) ([]*User, int64, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

// ProjectRepository 项目数据访问接口
type ProjectRepository interface {
	GetByID(ctx context.Context, id uint) (*Project, error)
	GetBySlug(ctx context.Context, slug string) (*Project, error)
	GetAll(ctx context.Context, limit, offset int, keyword string) ([]*Project, int64, error)
	Create(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uint) error
}

// LanguageRepository 语言数据访问接口
type LanguageRepository interface {
	GetByID(ctx context.Context, id uint) (*Language, error)
	GetByCode(ctx context.Context, code string) (*Language, error)
	GetAll(ctx context.Context) ([]*Language, error)
	Create(ctx context.Context, language *Language) error
	Update(ctx context.Context, language *Language) error
	Delete(ctx context.Context, id uint) error
	GetDefault(ctx context.Context) (*Language, error)
}

// TranslationRepository 翻译数据访问接口
type TranslationRepository interface {
	GetByID(ctx context.Context, id uint) (*Translation, error)
	GetByProjectID(ctx context.Context, projectID uint, limit, offset int) ([]*Translation, int64, error)
	GetByProjectAndLanguage(ctx context.Context, projectID, languageID uint) ([]*Translation, error)
	GetByProjectKeyLanguage(ctx context.Context, projectID uint, keyName string, languageID uint) (*Translation, error)
	GetMatrix(ctx context.Context, projectID uint, limit, offset int, keyword string) (map[string]map[string]string, int64, error)
	Create(ctx context.Context, translation *Translation) error
	CreateBatch(ctx context.Context, translations []*Translation) error
	UpsertBatch(ctx context.Context, translations []*Translation) error
	Update(ctx context.Context, translation *Translation) error
	Delete(ctx context.Context, id uint) error
	DeleteBatch(ctx context.Context, ids []uint) error
}

// ProjectMemberRepository 项目成员数据访问接口
type ProjectMemberRepository interface {
	GetByProjectAndUser(ctx context.Context, projectID, userID uint) (*ProjectMember, error)
	GetByProjectID(ctx context.Context, projectID uint) ([]*ProjectMember, error)
	GetByUserID(ctx context.Context, userID uint) ([]*ProjectMember, error)
	Create(ctx context.Context, member *ProjectMember) error
	Update(ctx context.Context, member *ProjectMember) error
	Delete(ctx context.Context, projectID, userID uint) error
}
