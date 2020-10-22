package response

import "QuickPass/models"

//获取代理的通道信息列表
type RespChannelInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Channel `json:"list"`
}
