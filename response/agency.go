package response

import "QuickPass/models"

// 代理信息统计
type RespAgencyStatistics struct {
	//已有商家
	MerchantCount int64 `json:"merchant_count"`
	//已有承兑
	AcceptorCount int64 `json:"acceptor_count"`
	//承兑总数额
	TotalAmount int64 `json:"total_amount"`
	//手续费收入
	FeeAmount int64 `json:"fee_amount"`
}

//代理获取用户的账单统计
type RespUserBillStatisticsList struct {
	// 总数
	Total int64 `json:"total"`
	// 数据列表
	List []*models.BillStatistics `json:"list"`
}

type RespUserByRoleBillStatistics struct {
	// 承兑金额
	AcceptAmount int64 `json:"accept_amount"`
	// 承兑次数
	AcceptCount int64 `json:"accept_count"`
	// 提现金额
	WithdrawalAmount int64 `json:"withdrawal_amount"`
	// 提现次数
	WithdrawalCount int64 `json:"withdrawal_count"`
	// 充值金额
	RechargeAmount int64 `json:"recharge_amount"`
	// 充值次数
	RechargeCount int64 `json:"recharge_count"`
}
