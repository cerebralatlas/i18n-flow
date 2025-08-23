package code

// ==========================
// 通用状态码
// ==========================
const (
	Success          = 0    // 成功
	ErrInternal      = 1000 // 内部错误
	ErrInvalidParams = 1001 // 参数错误
	ErrUnauthorized  = 1002 // 未认证或 token 无效
	ErrForbidden     = 1003 // 无权限操作
	ErrNotFound      = 1004 // 资源未找到
)

// ==========================
// 用户模块 (2000~2999)
// ==========================
const (
	ErrUserNotFound      = 2001 // 用户不存在
	ErrUserExists        = 2002 // 用户已存在
	ErrPasswordWrong     = 2003 // 密码错误
	ErrLoginFailed       = 2004 // 登录失败
	ErrTokenExpired      = 2005 // token 过期
	ErrUsernameDuplicate = 2006 // 用户名重复
)

// ==========================
// 项目模块 (3000~3999)
// ==========================
const (
	ErrProjectNotFound = 3001 // 项目不存在
	ErrProjectExists   = 3002 // 项目已存在
)

// ==========================
// 翻译模块 (4000~4999)
// ==========================
const (
	ErrTranslationNotFound = 4001 // 翻译不存在
	ErrTranslationExists   = 4002 // 翻译已存在
)

// ==========================
// 错误码对应提示信息
// ==========================
var Msg = map[int]string{
	Success:          "success",
	ErrInternal:      "internal server error",
	ErrInvalidParams: "invalid parameters",
	ErrUnauthorized:  "unauthorized",
	ErrForbidden:     "forbidden",
	ErrNotFound:      "resource not found",

	// 用户模块
	ErrUserNotFound:      "user not found",
	ErrUserExists:        "user already exists",
	ErrPasswordWrong:     "password incorrect",
	ErrLoginFailed:       "login failed",
	ErrTokenExpired:      "token expired",
	ErrUsernameDuplicate: "username already exists",

	// 项目模块
	ErrProjectNotFound: "project not found",
	ErrProjectExists:   "project already exists",

	// 翻译模块
	ErrTranslationNotFound: "translation not found",
	ErrTranslationExists:   "translation already exists",
}

// GetMsg 获取错误码对应提示
func GetMsg(code int) string {
	if msg, ok := Msg[code]; ok {
		return msg
	}
	return "unknown error"
}
