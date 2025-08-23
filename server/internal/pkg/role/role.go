package role

type Role string

const (
	Admin          Role = "admin"           // 系统管理员
	ProjectManager Role = "project_manager" // 项目管理员
	Translator     Role = "translator"      // 翻译人员
	Viewer         Role = "viewer"          // 只读用户
)
