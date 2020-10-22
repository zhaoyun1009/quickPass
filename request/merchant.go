package request

//获取商家信息列表
type ReqGetMerchantInfoListForm struct {
	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//获取商家卡信息列表
type ReqGetMerchantCardInfoListForm struct {
	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//添加商家
type ReqAddMerchantForm struct {
	// 商家账户名
	UserName string `form:"user_name" json:"user_name" binding:"required,len=6,number"`
	// 商家名称
	FullName string `form:"full_name" json:"full_name" binding:"required"`
	//密码
	Password string `form:"password" json:"password" binding:"required"`
	// 电话号码
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required,len=11"`
	// 地址
	Address string `form:"address" json:"address" binding:"required"`
	// 商家汇率(所有通道统一设置,默认0)
	Rate int64 `form:"rate" json:"rate"`
}

//删除商家
type ReqRemoveMerchantForm struct {
	// 商家账户名
	UserName string `form:"user_name" json:"user_name" binding:"required"`
}

//商家卖单
type ReqSellForm struct {
	// 卡信息
	// 收款人
	CardAccount string `form:"card_account" json:"card_account" binding:"required"`
	// 卡类型
	CardType string `form:"card_type" json:"card_type" binding:"required,oneof=BANK_CARD"`
	// 开户行
	CardBank string `form:"card_bank" json:"card_bank" `
	// 银行支行
	CardSubBank string `form:"card_sub_bank" json:"card_sub_bank"`
	// 卡号
	CardNo string `form:"card_no" json:"card_no"`

	// 收款信息
	// 金额
	Amount int64 `form:"amount" json:"amount" binding:"required,min=1,max=10000000000,sellAmountValidated"`
	// 备注
	Remark string `form:"remark" json:"remark"`
	// 交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
}

// 获取商家私钥
type ReqGetKeyForm struct {
	// 交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
}

// 代理获取商家私钥
type ReqAgencyGetMerchantKeyForm struct {
	// 商家用户名称
	Username string `form:"username" json:"username" binding:"required"`
	// 交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
}

// 添加商家卡
type ReqAddMerchantCard struct {
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

// 删除商家卡
type ReqRemoveMerchantCard struct {
	// 卡id
	CardId int64 `form:"card_id" json:"card_id" binding:"required"`
}

// 修改返回地址
type ReqUpdateReturnUrl struct {
	// 返回地址
	ReturnUrl string `form:"return_url" json:"return_url" binding:"required"`
}

// 修改回调地址
type ReqUpdateNotifyUrl struct {
	// 回调地址
	NotifyUrl string `form:"notify_url" json:"notify_url" binding:"required"`
}
