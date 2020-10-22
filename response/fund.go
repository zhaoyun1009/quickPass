package response

import "QuickPass/models"

//账单
type RespUserBillList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.Bill `json:"list"`
}

// 获取资金账户信息
type RespGetFundInfo struct {
	// 代理
	Agency string `json:"agency"`
	// 用户名
	UserName string `json:"user_name"`
	// 账户类型(1：系统资金账户  2：普通资金账户)
	Type int `json:"type"`
	// 可用资金
	AvailableAmount int64 `json:"available_amount"`
	// 冻结资金
	FrozenAmount int64 `json:"frozen_amount"`
}

// 统计收入
type RespInComeStatistics struct {
	// 代理
	Agency string `json:"agency"`
	// 已方交易账号
	Username string `json:"username"`
	// 涉及金额
	Amount int64 `json:"amount"`
	// 收支类型（1：支出，2：收入）
	IncomeExpensesType int `json:"income_expenses_type"`
}
