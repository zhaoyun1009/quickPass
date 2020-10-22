package response

import "QuickPass/models"

type RespAbnormalOrderInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.AbnormalOrder `json:"list"`
}
