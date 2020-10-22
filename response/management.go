package response

import "QuickPass/models"

//后台管理登录
type RespManagementLogin struct {
	//凭证
	Token string `json:"token"`
	//角色
	Role int8 `json:"role"`
	//规则
	Rules string `json:"rules"`
}

//获取后台管理员信息列表
type RespManagementInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Management `json:"list"`
}

// 代理列表
type ReqAgencyList struct {
	List []string `json:"list"`
}
