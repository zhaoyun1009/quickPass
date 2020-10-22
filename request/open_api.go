package request

// 商家买单签名数据
type ReqBuySignData struct {
	// 代理
	Agency string `form:"agency" json:"agency" binding:"required"`
	// 会员
	Member string `form:"member" json:"member" binding:"required"`
	// 商家
	Merchant string `form:"merchant" json:"merchant" binding:"required"`
	// 支付类型(BANK_CARD,ALIPAY,WECHAT)
	PayType string `form:"pay_type" json:"pay_type" binding:"required,oneof=BANK_CARD ALIPAY WECHAT"`
	// 支付金额
	Amount int64 `form:"amount" json:"amount" binding:"required,amountValidated"`
	// 商户订单号，需要保证不重复
	MerchantOrderNo string `json:"merchant_order_no" binding:"required"`
	// returnUrl
	ReturnUrl string `json:"return_url"`
	// notifyUrl
	NotifyUrl string `json:"notify_url"`
	// 接口模式[page, api]
	Model string `json:"model" binding:"required,oneof=page api"`
	// 摘要
	AppendInfo string `json:"append_info"`
}

//接口买入
type ReqOpenApiBuyForm struct {
	// 签名的数据
	SignData ReqBuySignData `json:"sign_data" binding:"required"`
	// 商户请求参数的签名串
	Sign string `json:"sign" binding:"required"`
}

// 商家订单状态查询
type ReqGetStatusByMerchantOrderNoForm struct {
	// 代理
	Agency string `form:"agency" json:"agency" binding:"required"`
	// 商家用户名
	MerchantUsername string `form:"merchant_username" json:"merchant_username" binding:"required"`
	// 商家订单号
	MerchantOrderNo string `form:"merchant_order_no" json:"merchant_order_no" binding:"required"`
}
