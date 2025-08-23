package response

import (
	"net/http"

	"i18n-flow/internal/pkg/code"

	"github.com/gin-gonic/gin"
)

// BaseResp 统一返回结构
type BaseResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功返回
func Success(c *gin.Context, data interface{}, message ...string) {
	msg := code.GetMsg(code.Success)
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(200, BaseResp{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	})
}

// 失败返回
func Error(c *gin.Context, code int, message string, data ...interface{}) {
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	c.JSON(200, BaseResp{
		Code:    code,
		Message: message,
		Data:    d,
	})
}

// 分页返回
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func Page(c *gin.Context, list interface{}, total int64, message ...string) {
	msg := code.GetMsg(code.Success)
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(200, BaseResp{
		Code:    http.StatusOK,
		Message: msg,
		Data: PageData{
			List:  list,
			Total: total,
		},
	})
}
