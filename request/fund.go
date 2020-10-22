package request

//代理转账
type ReqAgencyTransferForm struct {
	// 转入用户名
	ToUserName string `form:"to_user_name" json:"to_user_name" binding:"required"`
	// 转账金额
	Amount int64 `form:"amount" json:"amount" binding:"required,min=1,max=10000000000,amountValidated"`
	// 交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
	// 摘要
	AppendInfo string `form:"append_info" json:"append_info"`
}

//系统转账
type ReqSysTransferForm struct {
	// 代理 不传默认当前用户的同一代理
	ToAgency string `form:"agency" json:"agency" binding:"required"`
	// 转入用户名
	ToUserName string `form:"to_user_name" json:"to_user_name" binding:"required"`
	// 转账金额
	Amount int64 `form:"amount" json:"amount" binding:"required,min=1,amountValidated"`
	// 交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
	// 摘要
	AppendInfo string `form:"append_info" json:"append_info"`
}
