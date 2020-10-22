package request

import "QuickPass/pkg/util"

type CallbackParams struct {
	//商家订单号
	MerchantOrderNo string `json:"merchant_order_no"`
	// 订单状态（1：创建，2：待支付，3：待放行，4：已取消  5：已失败  6：已完成）
	OrderStatus int64 `json:"order_status"`
	// 会员号
	Member string `json:"member"`
	// 金额
	Amount int64 `json:"amount"`
	// 通道类型= (BANK_CARD,ALIPAY,WECHAT)
	ChannelType string `json:"channel_type"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time"`
	// 订单完成时间
	FinishTime *util.JSONTime `json:"finish_time"`
}

// 签名的回调参数
type SignCallbackParams struct {
	Sign string         `json:"sign"`
	Data CallbackParams `json:"data"`
}
