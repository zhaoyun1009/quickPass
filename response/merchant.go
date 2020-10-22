package response

import "QuickPass/models"

//获取商家信息列表
type RespMerchantInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Merchant `json:"list"`
}

//商家统计管理
type RespMerchantStatistics struct {
	//今日承兑数额
	TodayAcceptAmount int64 `json:"today_accept_amount"`
}

//获取商家卡信息列表
type RespMerchantCardInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.MerchantCard `json:"list"`
}

// 获取私钥
type RespGetKey struct {
	// 公钥
	PublicKey string `json:"public_key"`
	// 私钥
	PrivateKey string `json:"private_key"`
}

type RespApiCallUrl struct {
	// 商家接口返回url
	ReturnUrl string `json:"return_url"`
	// 商家接口回调url
	NotifyUrl string `json:"notify_url"`
}
