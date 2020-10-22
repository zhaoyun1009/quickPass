package request

//商家下单接口 创建订单
type ReqCreateOrderForm struct {
	// 商家账户名
	MerchantName string `form:"merchant_name" json:"merchant_name" binding:"required"`
	// 商家平台订单号
	MerchantOrderId string `form:"merchant_orderId" json:"merchant_orderId" binding:"required"`
	// 会员号
	Member string `form:"member" json:"member" binding:"required"`
	// 支付通道
	Channel string `form:"channel" json:"channel" binding:"required,oneof=BANK_CARD ALIPAY WECHAT"`
	// 金额
	Amount int64 `form:"amount" json:"amount" binding:"required"`
	// 回调接口
	CallbackUrl string `form:"callback_url" json:"callback_url" binding:"required"`
	// 支付成功后返回的URL
	ReturnUrl string `form:"return_url" json:"return_url" binding:"required"`
	// 签名
	Sign string `form:"sign" json:"sign" binding:"required"`
}

//获取下单信息
type ReqGetOrderingInfoForm struct {
	// 代理
	Agency string `form:"agency" json:"agency" binding:"required"`
	// 商家账户名
	MerchantName string `form:"merchant_name" json:"merchant_name" binding:"required"`
	// 订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

//获取商家订单信息列表
type ReqGetMerchantOrderInfosForm struct {
	// 订单号
	OrderNo int `form:"order_no" json:"order_no" `
	// 商家平台订单号
	MerchantOrderId string `form:"merchant_orderId" json:"merchant_orderId"`
	// 会员账号
	Member string `form:"member" json:"member"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 下单起始日期
	StartDate string `form:"max_amount" json:"max_amount"`
	// 下单截止日期
	EndDate string `form:"max_amount" json:"max_amount"`
}
