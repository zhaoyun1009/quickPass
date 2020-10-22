package models

import (
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type Bill struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 账单号
	BillNo string `json:"bill_no" gorm:"column:bill_no"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 角色（1:代理，2：承兑人、3：商家）
	OwnRole int `json:"own_role" gorm:"column:own_role"`
	// 已方交易账号
	OwnUserName string `json:"own_user_name" gorm:"column:own_user_name"`
	// 对方角色（1:代理，2：承兑人、3：商家）
	OppositeRole int `json:"opposite_role" gorm:"column:opposite_role"`
	// 对方交易账号
	OppositeUserName string `json:"opposite_user_name" gorm:"column:opposite_user_name"`
	// 涉及金额
	Amount int64 `json:"amount" gorm:"column:amount"`
	// 订单编号
	OrderNo string `json:"order_no" gorm:"order_no"`
	// 商家订单编号
	MerchantOrderNo string `json:"merchant_order_no" gorm:"merchant_order_no"`
	// 当前可用金额
	UsableAmount int64 `json:"usable_amount" gorm:"column:usable_amount"`
	// 当前冻结金额
	FrozenAmount int64 `json:"frozen_amount" gorm:"column:frozen_amount"`
	// 收支类型（1：支出，2：收入）
	IncomeExpensesType int `json:"income_expenses_type" gorm:"column:income_expenses_type"`
	// 会计科目（1：转账，2：转账手续费，3：收款手续费，4：OTC买入放行，5：OTC买入手续费 6：OTC提现手续费，7：OTC卖出提现放行
	BusinessType int `json:"business_type" gorm:"column:business_type"`
	// 摘要信息
	AppendInfo string `json:"append_info" gorm:"column:append_info"`
	// 当前汇率
	CurrentRate string `json:"current_rate" gorm:"column:current_rate"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

type BillExt struct {
	Bill
	// 数量,统计用
	Count int64 `json:"-" gorm:"column:count"`
}

// 设置Bill的表名为`bill`
func (Bill) TableName() string {
	return "bill"
}

func NewBillModel(session *Session) *Bill {
	return &Bill{Session: session}
}

// ExistBillByID checks if an bill exists based on ID
func (a *Bill) ExistBillByID(id int64) (bool, error) {
	bill := new(Bill)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(bill).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if bill.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetBillTotal gets the total number of bills based on the constraints
func (a *Bill) GetBillTotal(params *BillListQueryParams, maps interface{}) (int64, error) {
	var count int64
	query := getBillListParams(a, params)
	err := query.Where(maps).Count(&count).Error

	return count, err
}

// GetBills gets a list of bills based on paging constraints
func (a *Bill) GetBills(params *BillListQueryParams, pageNum int, pageSize int, maps interface{}) ([]*Bill, error) {
	var bills []*Bill
	query := getBillListParams(a, params)

	err := query.Where(maps).Offset(pageNum).Limit(pageSize).Find(&bills).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return bills, nil
}

func (a *Bill) RolesBillStatistics(params *BillListQueryParams, agency string, ownRole, incomeExpensesType, businessType int) (*Bill, error) {
	bill := new(Bill)
	query := getBillListParams(a, params)

	err := query.
		Select("? as agency, IFNULL(sum( amount ),0) as amount, ? as income_expenses_type", agency, incomeExpensesType).
		Where("agency = ? and own_role = ? and income_expenses_type = ? and business_type = ?", agency, ownRole, incomeExpensesType, businessType).
		First(bill).Error
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (a *Bill) UserBillStatistics(params *BillListQueryParams, agency, username string, incomeExpensesType int) (*Bill, error) {
	bill := new(Bill)
	query := getBillListParams(a, params)

	err := query.
		Select("? as agency, ? as own_user_name, IFNULL(sum( amount ),0) as amount, ? as income_expenses_type", agency, username, incomeExpensesType).
		Where("agency = ? and own_user_name = ? and income_expenses_type = ?", agency, username, incomeExpensesType).
		First(bill).Error
	if err != nil {
		return nil, err
	}

	return bill, nil
}

// 定时任务统计
func (a *Bill) CronBillStatistics(params *BillListQueryParams, agency, username string, incomeExpensesType, businessType int) (*BillExt, error) {
	bill := new(BillExt)
	query := getBillListParams(a, params)

	err := query.
		Table(a.TableName()).
		Select("? as agency, ? as own_user_name, IFNULL(sum( amount ),0) as amount, ? as income_expenses_type, ? as business_type, count(*) as count",
			agency,
			username,
			incomeExpensesType,
			businessType).
		Where("agency = ? and own_user_name = ? and income_expenses_type = ? and business_type = ? ", agency, username, incomeExpensesType, businessType).
		First(bill).Error
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func getBillListParams(a *Bill, params *BillListQueryParams) *gorm.DB {
	query := a.Session.db.Model(&Bill{})
	if params.StartTime != "" {
		query = query.Where("create_time >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		query = query.Where("create_time < ?", params.EndTime)
	}
	if params.MinAmount != 0 {
		query = query.Where("amount >= ?", params.MinAmount)
	}
	if params.MaxAmount != 0 {
		query = query.Where("amount <= ?", params.MaxAmount)
	}
	query = query.Order("id desc")
	return query
}

// GetBill Get a single bill based on ID
func (a *Bill) GetBill(id int64) (*Bill, error) {
	bill := new(Bill)
	err := a.Session.db.Where("id = ?", id).First(bill).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return bill, nil
}

// AddBill add a single bill
func (a *Bill) AddBill(bill *Bill) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(bill).Error
}

// DeleteBill delete a single bill
func (a *Bill) DeleteBill(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(Bill{}).Error
}

type BillListQueryParams struct {
	// 开始时间
	StartTime string `form:"start_time" json:"start_time"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time"`
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
}
