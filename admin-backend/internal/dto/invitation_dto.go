package dto

// CreateInvitationRequest 创建邀请请求
type CreateInvitationRequest struct {
	Role           string `json:"role" binding:"omitempty,oneof=admin member viewer"`
	ExpiresInDays  int    `json:"expires_in_days"`
	Description    string `json:"description"`
}

// CreateInvitationResponse 创建邀请响应
type CreateInvitationResponse struct {
	Code          string `json:"code"`
	InvitationURL string `json:"invitation_url"`
	Role          string `json:"role"`
	ExpiresAt     string `json:"expires_at"`
	Description   string `json:"description,omitempty"`
}

// InvitationInviter 邀请人信息
type InvitationInviter struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// InvitationResponse 邀请详情响应
type InvitationResponse struct {
	ID          uint64             `json:"id"`
	Code        string             `json:"code"`
	InviterID   uint64             `json:"inviter_id"`
	Inviter     *InvitationInviter `json:"inviter,omitempty"`
	Role        string             `json:"role"`
	Status      string             `json:"status"`
	ExpiresAt   string             `json:"expires_at"`
	UsedAt      *string            `json:"used_at,omitempty"`
	UsedBy      *uint64            `json:"used_by,omitempty"`
	Description string             `json:"description,omitempty"`
	CreatedAt   string             `json:"created_at"`
}

// InvitationListResponse 邀请列表响应
type InvitationListResponse struct {
	Invitations []*InvitationResponse `json:"invitations"`
	Total       int64                 `json:"total"`
}

// ValidateInvitationResponse 验证邀请响应
type ValidateInvitationResponse struct {
	Valid     bool               `json:"valid"`
	Inviter   *InvitationInviter `json:"inviter,omitempty"`
	Role      string             `json:"role"`
	ExpiresAt string             `json:"expires_at"`
	Message   string             `json:"message,omitempty"`
}

// RegisterWithInvitationRequest 使用邀请码注册请求
type RegisterWithInvitationRequest struct {
	Code     string `json:"code" binding:"required"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
