package response

import "QuickPass/models"

//获取承兑人信息列表
type RespAcceptorInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Acceptor `json:"list"`
}

//获取承兑人卡列表
type RespAcceptorCardInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.AcceptorCard `json:"list"`
}

//获取承兑人卡列表
type RespAcceptorAllCardInfo struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.AcceptorCard `json:"list"`
}

//承兑人信息统计
type RespAcceptorStatistics struct {
	//已有卡商
	TotalCount int64 `json:"total_count"`
	//当前在线
	OnLineCount int64 `json:"on_line_count"`
	//自动承兑开启
	AutoAcceptCount int64 `json:"auto_accept_count"`
	//今日承兑数额
	TodayAcceptAmount int64 `json:"today_accept_amount"`
}
