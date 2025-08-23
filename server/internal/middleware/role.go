package middleware

import (
	"net/http"

	"i18n-flow/internal/model"
	"i18n-flow/internal/pkg/response"
	"i18n-flow/internal/pkg/role"

	"github.com/gin-gonic/gin"
)

// RequireRoles 返回一个中间件，要求用户必须拥有其中一个角色才能访问
func RequireRoles(allowedRoles ...role.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userObj, exists := c.Get("user") // JWT 中间件必须先写入 "user"
		if !exists {
			response.Error(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		user, ok := userObj.(*model.User)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "用户信息异常")
			c.Abort()
			return
		}

		if !HasRole(user, allowedRoles...) {
			response.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// HasRole 判断用户是否拥有任意一个允许的角色
func HasRole(user *model.User, allowedRoles ...role.Role) bool {
	roleMap := make(map[role.Role]bool, len(user.Roles))
	for _, r := range user.Roles {
		roleMap[r] = true
	}

	for _, r := range allowedRoles {
		if roleMap[r] {
			return true
		}
	}

	return false
}
