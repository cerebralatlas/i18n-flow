package v1

import (
	"github.com/gin-gonic/gin"

	"i18n-flow/internal/dto/request"
	dto "i18n-flow/internal/dto/response"
	"i18n-flow/internal/pkg/code"
	"i18n-flow/internal/pkg/response"
	"i18n-flow/internal/pkg/validator"
	"i18n-flow/internal/service"
)

func Register(c *gin.Context) {
	var req request.RegisterReq
	if validator.BindAndValidate(c, &req) {
		return
	}

	if err := service.RegisterUser(req.Username, req.Password); err != nil {
		if err.Error() == code.GetMsg(code.ErrUsernameDuplicate) {
			response.Error(c, code.ErrUsernameDuplicate, code.GetMsg(code.ErrUsernameDuplicate))
		} else {
			response.Error(c, code.ErrInternal, code.GetMsg(code.ErrInternal))
		}
		return
	}

	response.Success(c, nil)
}

func Login(c *gin.Context) {
	var req request.LoginReq
	if validator.BindAndValidate(c, &req) {
		return
	}

	user, accessToken, refreshToken, err := service.LoginUser(req.Username, req.Password)
	if err != nil {
		response.Error(c, code.ErrUserNotFound, code.GetMsg(code.ErrUserNotFound))
		return
	}

	response.Success(c, dto.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserInfoResp{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}
