package request

//获取代理的通道列表
type ReqGetChannelsForm struct {
	//代理名称
	Agency string `form:"agency" json:"agency"`
}

//获取代理的通道信息列表
type ReqGetChannelInfoListForm struct {
	//起始页
	StartPage int `form:"start_page" json:"start_page"`
	//页面大小
	PageSize int `form:"page_size" json:"page_size"`
}

//更新商家通道汇率
type ReqUpdateMerchantSingleRate struct {
	//商家名称
	MerchantName string `form:"merchant_name" json:"merchant_name" binding:"required"`
	// 通道名称
	Channel string `form:"channel" json:"channel" binding:"omitempty,oneof=BANK_CARD ALIPAY WECHAT"`
	//通道汇率
	Rate int64 `form:"rate" json:"rate" binding:"min=0,max=10000"`
}

//统一更新商家通道汇率
type ReqUpdateMerchantAllChannelRate struct {
	//商家名称
	MerchantName string `form:"merchant_name" json:"merchant_name" binding:"required"`
	//通道汇率
	Rate int64 `form:"rate" json:"rate" binding:"min=0,max=10000"`
}

//更新每日承兑上限
type ReqUpdateLimitAmountForm struct {
	//通道ID
	Id int64 `form:"id" json:"id" binding:"required"`
	//每日承兑上限金额
	LimitAmount int64 `form:"limit_amount" json:"limit_amount" binding:"required,min=1"`
}

//统一更新每日承兑上限
type ReqUpdateAllLimitAmountForm struct {
	//每日承兑上限金额
	LimitAmount int64 `form:"limit_amount" json:"limit_amount" binding:"required,min=1"`
}

//更新买入最大金额
type ReqUpdateBuyMaxAmountForm struct {
	//通道ID
	Id     int64 `form:"id" json:"id" binding:"required"`
	Amount int64 `form:"amount" json:"amount" binding:"required,min=1,max=500000000,amountValidated"`
}

//更新买入最小金额
type ReqUpdateBuyMinAmountForm struct {
	//通道ID
	Id     int64 `form:"id" json:"id" binding:"required"`
	Amount int64 `form:"amount" json:"amount" binding:"required,min=1,max=500000000,amountValidated"`
}

//通道开关
type ReqChannelSwitchForm struct {
	//通道ID
	ChannelId int64 `form:"channel_id" json:"channel_id" binding:"required"`
	//通道开关(1: 关闭  2:开启)
	IfOpen int `form:"if_open" json:"if_open" binding:"required,oneof=1 2"`
}
