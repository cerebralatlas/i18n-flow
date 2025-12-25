package dto

// AddProjectMemberRequest 添加项目成员请求
type AddProjectMemberRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=owner editor viewer"`
}

// UpdateProjectMemberRequest 更新项目成员请求
type UpdateProjectMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=owner editor viewer"`
}

// ProjectMemberInfo 项目成员信息
type ProjectMemberInfo struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
