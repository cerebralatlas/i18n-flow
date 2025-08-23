package validator

import (
	"i18n-flow/internal/pkg/code"
	"i18n-flow/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

// BindAndValidate 绑定 JSON 并验证结构体字段
// 返回 true 表示出错并已经响应客户端，API 可以直接 return
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	// 绑定 JSON
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Error(c, code.ErrInvalidParams, code.GetMsg(code.ErrInvalidParams))
		return true
	}

	// 验证字段
	if err := validate.Struct(obj); err != nil {
		response.Error(c, code.ErrInvalidParams, err.Error())
		return true
	}

	return false
}
