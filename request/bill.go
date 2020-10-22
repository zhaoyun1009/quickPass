package request

//获取用户账单
type ReqGetUserBillListForm struct {
	//账单号
	BillNo string `form:"bill_no" json:"bill_no"`
	// 对方交易账户
	OppositeUserName string `form:"opposite_user_name" json:"opposite_user_name"`
	//收支类型（1：支出，2：收入）
	IncomeExpensesType int `form:"income_expenses_type" json:"income_expenses_type" binding:"omitempty,oneof=0 1 2"`

	// 订单编号
	OrderNo string `form:"order_no" json:"order_no"`
	// 商家订单编号
	MerchantOrderNo string `form:"merchant_order_no" json:"merchant_order_no"`

	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount" binding:"omitempty,number"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount" binding:"omitempty,number"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`

	StartPage int `form:"start_page" json:"start_page"`
	PageSize  int `form:"page_size" json:"page_size"`
}

//收入统计
type ReqInComeStatisticsForm struct {
	// 收支类型（1：支出，2：收入）
	IncomeExpensesType int `form:"income_expenses_type" json:"income_expenses_type" binding:"required,oneof=1 2"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`
}

// 代理获取承兑人账单
type ReqGetAcceptorBillForm struct {
	StartPage int `form:"start_page" json:"start_page"`
	PageSize  int `form:"page_size" json:"page_size"`
}
