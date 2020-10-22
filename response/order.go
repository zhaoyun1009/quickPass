package response

import (
	"QuickPass/models"
	"QuickPass/pkg/util"
)

//获取订单号
type RespOrderId struct {
	//订单ID
	OrderId int64 `json:"order_id"`
}

//买单
type RespOrderBuy struct {
	//订单号
	OrderNo string `json:"order_no"`
	//卡号
	CardNo string `json:"card_no"`
	//卡账号
	CardAccount string `json:"card_account"`
	//银行
	CardBank string `json:"card_bank"`
	//卡图片
	CardImg string `json:"card_img"`
	//银行支行
	CardSubBank string `json:"card_sub_bank"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time"`
	// 服务器当前时间
	CurrentTime *util.JSONTime `json:"current_time"`
	//过期时间
	ExpirationTime util.JSONTime `json:"expiration_time"`
}

// 订单状态查询
type RespGetOrderStatus struct {
	//代理
	Agency string `json:"agency"`
	//订单号
	OrderNo string `json:"order_no"`
	// 订单类型（1：转账，2：买入，3：卖出）
	OrderType int64 `json:"order_type"`
	// 订单状态（1：创建，2：待支付，3：待放行，4：已取消  5：已失败  6：已完成）
	OrderStatus int64 `json:"order_status"`
	// 金额
	Amount int64 `json:"amount"`
	// 通道类型= (BANK_CARD,ALIPAY,WECHAT)
	ChannelType string `json:"channel_type"`
	//卡号
	CardNo string `json:"card_no"`
	//卡账号
	CardAccount string `json:"card_account"`
	//银行
	CardBank string `json:"card_bank"`
	//卡图片
	CardImg string `json:"card_img"`
	//银行支行
	CardSubBank string `json:"card_sub_bank"`
	// 返回地址
	ReturnUrl string `json:"returnUrl"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time"`
	// 服务器当前时间
	CurrentTime *util.JSONTime `json:"current_time"`
	//过期时间
	ExpirationTime util.JSONTime `json:"expiration_time"`
	// 订单完成时间
	FinishTime *util.JSONTime `json:"finish_time"`
}

// 买单确认
type RespConfirmPay struct {
	//过期时间
	ExpirationTime util.JSONTime `json:"expiration_time"`
}

//获取买单列表
type RespBuyOrderInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Order `json:"list"`
}

//代理获取买单列表
type RespAgencyBuyOrderInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Order `json:"list"`
}

//商家获取买单列表
type RespMerchantBuyOrderInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Order `json:"list"`
}

//获取卖单列表
type RespSellOrderInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Order `json:"list"`
}

//获取承兑人名称列表
type RespGetAcceptorGroup struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []string `json:"list"`
}

//代理卖单未处理数
type RespGetSellOrderUnprocessed struct {
	// 未处理数
	UnprocessedCount int64 `json:"unprocessed_count"`
}
