package request

//承兑人订单查询
type ReqBuyOrderInfoList struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no"`
	// 订单状态
	OrderStatus int64 `form:"order_status" json:"order_status" binding:"omitempty,oneof=0 2 3 4 6"`
	// 交易通道(BANK_CARD,ALIPAY,WECHAT)
	Channel string `form:"channel" json:"channel" binding:"omitempty,oneof=BANK_CARD ALIPAY WECHAT"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`

	// 完成开始时间
	FinishStartTime string `form:"finish_start_time" json:"finish_start_time" binding:"timeValidated"`
	// 完成结束时间
	FinishEndTime string `form:"finish_end_time" json:"finish_end_time" binding:"timeValidated"`

	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//商家买入订单查询
type ReqMerchantBuyOrderInfoList struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no"`
	// 商家订单号
	MerchantOrderNo string `form:"merchant_order_no" json:"merchant_order_no"`
	// 订单状态
	//OrderStatus int64 `form:"order_status" json:"order_status"`
	// 交易通道(BANK_CARD,ALIPAY,WECHAT)
	Channel string `form:"channel" json:"channel" binding:"omitempty,oneof=BANK_CARD ALIPAY WECHAT"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`

	// 完成开始时间
	FinishStartTime string `form:"finish_start_time" json:"finish_start_time" binding:"timeValidated"`
	// 完成结束时间
	FinishEndTime string `form:"finish_end_time" json:"finish_end_time" binding:"timeValidated"`

	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//商家卖出订单查询
type ReqMerchantSellOrderInfoList struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no"`
	// 订单状态
	OrderStatus int64 `form:"order_status" json:"order_status"`
	// 交易通道(BANK_CARD,ALIPAY,WECHAT)
	Channel string `form:"channel" json:"channel" binding:"omitempty,oneof=BANK_CARD ALIPAY WECHAT"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`

	// 完成开始时间
	FinishStartTime string `form:"finish_start_time" json:"finish_start_time" binding:"timeValidated"`
	// 完成结束时间
	FinishEndTime string `form:"finish_end_time" json:"finish_end_time" binding:"timeValidated"`

	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//代理查询商家卖出订单查询
type ReqAgencySellOrderInfoList struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no"`
	// 订单状态
	OrderStatus int64 `form:"order_status" json:"order_status"`
	// 交易通道(BANK_CARD,ALIPAY,WECHAT)
	Channel string `form:"channel" json:"channel" binding:"omitempty,oneof=BANK_CARD ALIPAY WECHAT"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`

	// 完成开始时间
	FinishStartTime string `form:"finish_start_time" json:"finish_start_time" binding:"timeValidated"`
	// 完成结束时间
	FinishEndTime string `form:"finish_end_time" json:"finish_end_time" binding:"timeValidated"`

	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//代理订单信息查询
type ReqGetAgencyOrderInfoList struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no"`
	// 商家订单号
	MerchantOrderNo string `form:"merchant_order_no" json:"merchant_order_no"`
	// 订单状态
	OrderStatus int64 `form:"order_status" json:"order_status"`
	// 交易通道(BANK_CARD,ALIPAY,WECHAT)
	Channel string `form:"channel" json:"channel" binding:"omitempty,oneof=BANK_CARD ALIPAY WECHAT"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`
	// 所属承兑人
	AcceptorName string `form:"acceptor_name" json:"acceptor_name"`

	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//代理放行
type ReqAgencyDischargeForm struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

//代理放行超时取消的订单
type ReqAgencyDischargeCancelOrderForm struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

//代理再次回调商家
type ReqAgencyAgainCallback struct {
	// 订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
	// 回调地址
	CallbackUrl string `form:"callback_url" json:"callback_url" binding:"required"`
}

//商家订单查询
type ReqSellOrderInfoList struct {
	// 订单号
	OrderId int64 `form:"order_id" json:"order_id"`
	// 订单状态
	OrderStatus int64 `form:"order_status" json:"order_status"`
	// 交易通道
	Channel string `form:"channel" json:"channel" binding:"omitempty,oneof=BANK_CARD ALIPAY WECHAT"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`

	// 完成开始时间
	FinishStartTime string `form:"finish_start_time" json:"finish_start_time" binding:"timeValidated"`
	// 完成结束时间
	FinishEndTime string `form:"finish_end_time" json:"finish_end_time" binding:"timeValidated"`

	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//买入
type ReqBuyForm struct {
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
}

// 订单状态查询
type ReqGetStatusByOrderNoForm struct {
	// 代理
	Agency string `form:"agency" json:"agency" binding:"required"`
	// 订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

//取消买入订单
type ReqCancelBuyForm struct {
	// 代理
	Agency string `form:"agency" json:"agency" binding:"required"`
	//订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

//取消卖出订单
type ReqCancelSellForm struct {
	//订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

//确认支付
type ReqConfirmPayForm struct {
	// 代理
	Agency string `form:"agency" json:"agency" binding:"required"`
	//订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

type ReqConfirmSellPayForm struct {
	//订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
}

//商家放行卖单
type ReqMerchantDischargeSellForm struct {
	//订单号
	OrderNo string `form:"order_no" json:"order_no" binding:"required"`
	//交易密码
	TradeKey string `form:"trade_key" json:"trade_key" binding:"required,len=6"`
}
