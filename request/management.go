package request

//后台管理登录
type ReqManagementLoginForm struct {
	//用户名
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	//密码
	Password string `form:"password" json:"password" binding:"required"`
}
