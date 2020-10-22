package response

// 查询列表返回值
type RespQueryList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List interface{} `json:"list"`
}

//刷新token
type RespRefreshToken struct {
	Token string `json:"token"`
}

// 获取承兑人信息
type RespAcceptorInfo struct {
	//总余额
	TotalAmount int64 `json:"total_amount"`
	// 可用余额
	AvailableAmount int64 `json:"available_amount"`
	// 冻结余额
	FrozenAmount int64 `json:"frozen_amount"`
	// 承兑开关(1：关闭，2：开启)
	AcceptSwitch int `json:"accept_switch"`
	// 是否自动承兑（1：不自动，2：自动）
	IfAutoAccept int `json:"if_auto_accept"`
}
