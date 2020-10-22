package request

//获取用户信息列表
type ReqGetUserInfoListForm struct {
	// 代理 所属系统 默认admin
	Agency string `form:"agency" json:"agency" binding:"required"`
	// 起始页
	StartPage int `form:"start_age" json:"start_age" binding:"required"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size" binding:"required"`
}

//添加代理用户
type ReqAddAgencyUserForm struct {
	// 代理账户名
	UserName string `form:"user_name" json:"user_name" binding:"required,len=6,number"`
	// 代理昵称
	FullName string `form:"full_name" json:"full_name" binding:"required"`
	// 代理密码
	Password string `form:"password" json:"password" binding:"required"`
	// 电话号码
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required,max=11,number"`
	// 地址
	Address string `form:"address" json:"address" binding:"required"`
}

//用户登录请求信息
type ReqUserLoginForm struct {
	// 代理
	Agency string `form:"agency" json:"agency" binding:"required"`
	// 用户名
	UserName string `form:"user_name" json:"user_name" binding:"required,max=255"`
	// 密码
	Password string `form:"password" json:"password" binding:"required,max=255"`
}

//更新用户信息
type ReqUserUpdateForm struct {
	// 电话号码
	PhoneNumber string `form:"phone_number" json:"phone_number" `
	// 姓名
	FullName string `form:"full_name" json:"full_name"`
	// 地址
	Address string `form:"address" json:"address" `
}

//修改登录密码
type ReqPasswordModifyForm struct {
	// 原始密码
	OriginalPassword string `form:"original_password" json:"original_password" binding:"required,max=255"`
	// 更新密码
	UpdatePassword string `form:"update_password" json:"update_password" binding:"required,max=255"`
}

//修改交易密码
type ReqTradeKeyModifyForm struct {
	// 原始交易密码
	OriginalTradePassword string `form:"original_trade_password" json:"original_trade_password" binding:"required,len=6"`
	// 更新交易密码
	UpdateTradePassword string `form:"update_trade_password" json:"update_trade_password" binding:"required,len=6"`
}

//设置交易密码
type ReqSettingTradeKeyForm struct {
	// 交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
}

//重置登录密码
type ReqAgencyResetUserPasswordForm struct {
	// 用户名称
	Username string `form:"username" json:"username" binding:"required,len=6"`
}

//重置交易密码
type ReqAgencyResetUserTradeKeyForm struct {
	// 用户名称
	Username string `form:"username" json:"username" binding:"required,len=6"`
}
