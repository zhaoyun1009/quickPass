package response

import "QuickPass/models"

//后台管理登录
type RespUserLogin struct {
	//凭证
	Token string `json:"token"`
	//角色 （1:代理，2：承兑人、3：商家）
	Role int `json:"role"`
	//规则
	Rules string `json:"rules"`
}

//获取用户信息列表
type RespUserInfoList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.User `json:"list"`
}

//是否存在交易密码
type RespExistTradeKey struct {
	// true存在 false不存在
	Exist bool `json:"exist"`
}
