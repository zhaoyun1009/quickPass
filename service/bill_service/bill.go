package bill_service

import (
	"QuickPass/models"
	"QuickPass/pkg/util"
)

type Bill struct {
	//
	Id int64
	// 账单号
	BillNo string
	// 代理
	Agency string
	// 已方交易账号
	OwnUserName string
	// 角色（1:代理，2：承兑人、3：商家）
	OwnRole int
	// 对方交易账号
	OppositeUserName string
	// 对方角色（1:代理，2：承兑人、3：商家）
	OppositeRole int
	// 涉及金额
	Amount int64
	// 订单编号
	OrderNo string
	// 商家订单编号
	MerchantOrderNo string
	// 当前可用金额
	UsableAmount int64
	// 当前冻结金额
	FrozenAmount int64
	// 收支类型（1：支出，2：收入）
	IncomeExpensesType int
	// 会计科目（1：转账，2：转账手续费，3：收款手续费，4：OTC买入放行，5：OTC买入手续费 6：OTC提现手续费，7：OTC卖出提现放行）
	BusinessType int64
	// 摘要信息
	AppendInfo string
	// 当前汇率
	CurrentRate string
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int

	//------------------------
	StartTime string
	EndTime   string
	// 最大金额
	MaxAmount int64
	// 最小金额
	MinAmount int64
}

func (a *Bill) Get() (*models.Bill, error) {
	session := models.NewSession()
	bill, err := models.NewBillModel(session).GetBill(a.Id)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (a *Bill) GetAll() ([]*models.Bill, error) {
	session := models.NewSession()
	return models.NewBillModel(session).GetBills(&models.BillListQueryParams{
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		MaxAmount: a.MaxAmount,
		MinAmount: a.MinAmount,
	}, a.PageNum, a.PageSize, a.getMaps())
}

func (a *Bill) Delete() error {
	session := models.NewSession()
	return models.NewBillModel(session).DeleteBill(a.Id)
}

func (a *Bill) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewBillModel(session).ExistBillByID(a.Id)
}

func (a *Bill) Count() (int64, error) {
	session := models.NewSession()
	return models.NewBillModel(session).GetBillTotal(&models.BillListQueryParams{
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		MaxAmount: a.MaxAmount,
		MinAmount: a.MinAmount,
	}, a.getMaps())
}

// 统计单个代理(承兑人/商家)收入支出信息
func (a *Bill) RolesBillStatistics(agency, startTime, endTime string, ownRole, incomeExpensesType, businessType int) (*models.Bill, error) {
	session := models.NewSession()
	billModel := models.NewBillModel(session)
	return billModel.RolesBillStatistics(&models.BillListQueryParams{
		StartTime: startTime,
		EndTime:   endTime,
	}, agency, ownRole, incomeExpensesType, businessType)
}

// 统计单个用户收入支出信息
func (a *Bill) UserBillStatistics(agency, username, startTime, endTime string, incomeExpensesType int) (*models.Bill, error) {
	session := models.NewSession()
	billModel := models.NewBillModel(session)
	return billModel.UserBillStatistics(&models.BillListQueryParams{
		StartTime: startTime,
		EndTime:   endTime,
	}, agency, username, incomeExpensesType)
}

// 定时任务统计
func (a *Bill) CronBillStatistics(agency, username, startTime, endTime string, incomeExpensesType, businessType int) (*models.BillExt, error) {
	session := models.NewSession()
	billModel := models.NewBillModel(session)
	return billModel.CronBillStatistics(&models.BillListQueryParams{
		StartTime: startTime,
		EndTime:   endTime,
	}, agency, username, incomeExpensesType, businessType)
}

func (a *Bill) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.BillNo != "" {
		maps["bill_no"] = a.BillNo
	}
	if a.OwnUserName != "" {
		maps["own_user_name"] = a.OwnUserName
	}
	if a.OppositeUserName != "" {
		maps["opposite_user_name"] = a.OppositeUserName
	}
	if a.OrderNo != "" {
		maps["order_no"] = a.OrderNo
	}
	if a.MerchantOrderNo != "" {
		maps["merchant_order_no"] = a.MerchantOrderNo
	}
	if a.IncomeExpensesType != 0 {
		maps["income_expenses_type"] = a.IncomeExpensesType
	}
	if a.BusinessType != 0 {
		maps["business_type"] = a.BusinessType
	}

	return maps
}
