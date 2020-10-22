package request

type ReqGetAcceptorAbnormalOrderInfoList struct {
	// 异常单号
	AbnormalOrderNo string `form:"abnormal_order_no" json:"abnormal_order_no"`
	// 异常类型（1: 未知, 2：超时取消，3：未生成订单， 4：订单金额不符）
	AbnormalOrderType int64 `form:"abnormal_order_type" json:"abnormal_order_type" binding:"omitempty,oneof=0 1 2 3 4"`
	// 订单状态（1：未处理  2：已处理  3：处理中）
	AbnormalOrderStatus int64 `form:"abnormal_order_status" json:"abnormal_order_status" binding:"omitempty,oneof=0 1 2 3"`

	// 最大金额
	MaxAmount string `form:"max_amount" json:"max_amount" binding:"omitempty,number"`
	// 最小金额
	MinAmount string `form:"min_amount" json:"min_amount" binding:"omitempty,number"`
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

type ReqGetAgencyAbnormalOrderInfoList struct {
	// 异常单号
	AbnormalOrderNo string `form:"abnormal_order_no" json:"abnormal_order_no"`
	// 异常类型（1: 未知, 2：超时取消，3：未生成订单， 4：订单金额不符）
	AbnormalOrderType int64 `form:"abnormal_order_type" json:"abnormal_order_type" binding:"omitempty,oneof=0 1 2 3 4"`
	// 订单状态（1：未处理  2：已处理  3：处理中）
	AbnormalOrderStatus int64 `form:"abnormal_order_status" json:"abnormal_order_status" binding:"omitempty,oneof=0 1 2 3"`
	// 所属承兑人
	AcceptorName string `form:"acceptor_name" json:"acceptor_name"`

	// 最大金额
	MaxAmount string `form:"max_amount" json:"max_amount" binding:"omitempty,number"`
	// 最小金额
	MinAmount string `form:"min_amount" json:"min_amount" binding:"omitempty,number"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`

	// 起始页
	StartPage int `form:"start_page" json:"start_page"`
	// 页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

// 新建异常单
type ReqAddAcceptorAbnormalOrderFrom struct {
	// 收款通道({“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"})
	Channel string `form:"channel" json:"channel" binding:"required,oneof=BANK_CARD ALIPAY WECHAT"`
	// 金额
	Amount int64 `form:"amount" json:"amount" binding:"required,min=1,max=500000000,amountValidated"`
	// 收款账户
	AcceptCardAccount string `form:"accept_card_account" json:"accept_card_account" binding:"required"`
	// 收款卡号
	AcceptCardNo string `form:"accept_card_no" json:"accept_card_no" binding:"required"`
	// 收款卡银行
	AcceptCardBank string `form:"accept_card_bank" json:"accept_card_bank"`
	// 收款时间
	PaymentDate string `form:"payment_date" json:"payment_date" binding:"required,timeValidated"`
}

// 更新异常单生成订单
type ReqUpdateAbnormalOrderFrom struct {
	// 异常单号
	AbnormalOrderNo string `form:"abnormal_order_no" json:"abnormal_order_no" binding:"required"`
	// 商家账户
	Merchant string `form:"merchant" json:"merchant" binding:"required"`
	// 会员名称
	Member string `form:"member" json:"member"`
	// 实际金额
	Amount int64 `form:"amount" json:"amount" binding:"required"`
	// 异常类型（1: 未知, 2：超时取消，3：未生成订单， 4：订单金额不符）
	AbnormalOrderType int64 `form:"abnormal_order_type" json:"abnormal_order_type" binding:"required,oneof=1 2 3 4"`
	// 摘要
	AppendInfo string `form:"append_info" json:"append_info"`
}
