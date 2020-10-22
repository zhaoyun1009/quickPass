package request

//增加承兑人 收款方式
type ReqAddAcceptorCardForm struct {
	//卡类型={“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"}
	CardType string `form:"card_type" json:"card_type" binding:"required,oneof=BANK_CARD ALIPAY WECHAT"`
	//卡号
	CardNo string `form:"card_no" json:"card_no" binding:"required"`
	//卡账户名
	CardAccount string `form:"card_account" json:"card_account" binding:"required"`
	// 银行名称
	CardBank string `form:"card_bank" json:"card_bank"`
	// 支行
	CardSubBank string `form:"card_sub_bank" json:"card_sub_bank" `
	// 卡图片
	CardImg string `form:"card_img" json:"card_img"`
}

//更新承兑人卡信息
type ReqUpdateAcceptorCardForm struct {
	// 卡id
	CardId int64 `form:"card_id" json:"card_id" binding:"required"`
	//卡账户名
	CardAccount string `form:"card_account" json:"card_account" binding:"required"`
	// 银行名称
	CardBank string `form:"card_bank" json:"card_bank"`
	// 支行
	CardSubBank string `form:"card_sub_bank" json:"card_sub_bank" `
	// 卡图片
	CardImg string `form:"card_img" json:"card_img"`
}

//更新承兑人卡状态
type ReqUpdateCardStatusForm struct {
	// 卡id
	CardId int64 `form:"card_id" json:"card_id" binding:"required"`
	// 是否开启在线={1:关闭，2：开启}
	IfOpen int `form:"if_open" json:"if_open" binding:"required,oneof=1 2"`
}

//删除承兑人卡
type ReqDeleteAcceptorCardForm struct {
	// 卡id
	CardId int64 `form:"card_id" json:"card_id" binding:"required"`
}

//获取承兑人卡列表
type ReqGetAcceptorCardInfoListForm struct {
	// 卡类型={“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"}
	CardType string `form:"card_type" json:"card_type"`
	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//更新承兑权限
type ReqUpdateAcceptSwitchFrom struct {
	//承兑开关(1：关闭，2：开启)
	AcceptSwitch int `form:"accept_switch" json:"accept_switch" binding:"required,oneof=1 2"`
}

//代理更新承兑权限
type ReqAgencyUpdateAcceptSwitchFrom struct {
	//承兑开关(1：关闭，2：开启)
	AcceptSwitch int `form:"accept_switch" json:"accept_switch" binding:"required,oneof=1 2"`
	//承兑人
	Username string `form:"username" json:"username" binding:"required"`
}

//代理更新承兑账户状态
type ReqAgencyUpdateAcceptStatusFrom struct {
	//承兑账户状态开关(1：关闭，2：开启)
	AcceptStatus int `form:"accept_status" json:"accept_status" binding:"required,oneof=1 2"`
	//承兑人
	Username string `form:"username" json:"username" binding:"required"`
}

//更新承兑权限
type ReqUpdateIfAutoAcceptFrom struct {
	//自动承兑(1：关闭，2：开启)
	IfAutoAccept int `form:"if_auto_accept" json:"if_auto_accept" binding:"required,oneof=1 2"`
}

//代理更新承兑权限
type ReqAgencyUpdateIfAutoAcceptFrom struct {
	//自动承兑(1：关闭，2：开启)
	IfAutoAccept int `form:"if_auto_accept" json:"if_auto_accept" binding:"required,oneof=1 2"`
	//承兑人
	Username string `form:"username" json:"username" binding:"required"`
}

//承兑人放行
type ReqAcceptorDischargeForm struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
	// 交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
}

//增加承兑人请求参数
type ReqAddAcceptorForm struct {
	//承兑人账户名
	UserName string `form:"user_name" json:"user_name" binding:"required,len=6,number"`
	//承兑人姓名
	FullName string `form:"full_name" json:"full_name" binding:"required"`
	//密码
	Password string `form:"password" json:"password" binding:"required"`
	//电话号码
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required,len=11"`
	//地址
	Address string `form:"address" json:"address" binding:"required"`
}

//获取承兑人信息列表
type ReqGetAcceptorInfoListForm struct {
	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}
