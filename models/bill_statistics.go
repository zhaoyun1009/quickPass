package models

import (
	"QuickPass/pkg/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

type BillStatistics struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 用户名
	Username string `json:"username" gorm:"column:username"`
	// 角色
	Role int `json:"role" gorm:"column:role"`
	// 统计日期
	StatisticsDate string `json:"statistics_date" gorm:"column:statistics_date"`
	// 剩余金额
	LeftAmount int64 `json:"left_amount" gorm:"column:left_amount"`
	// 承兑金额
	AcceptAmount int64 `json:"accept_amount" gorm:"column:accept_amount"`
	// 承兑次数
	AcceptCount int64 `json:"accept_count" gorm:"column:accept_count"`
	// 提现金额
	WithdrawalAmount int64 `json:"withdrawal_amount" gorm:"column:withdrawal_amount"`
	// 提现次数
	WithdrawalCount int64 `json:"withdrawal_count" gorm:"column:withdrawal_count"`
	// 充值金额
	RechargeAmount int64 `json:"recharge_amount" gorm:"column:recharge_amount"`
	// 充值次数
	RechargeCount int64 `json:"recharge_count" gorm:"column:recharge_count"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置BillStatistics的表名为`bill_statistics`
func (BillStatistics) TableName() string {
	return "bill_statistics"
}

func NewBillStatisticsModel(session *Session) *BillStatistics {
	return &BillStatistics{Session: session}
}

// ExistBillStatisticsByID checks if an billStatistics exists based on ID
func (a *BillStatistics) ExistBillStatisticsByID(id int64) (bool, error) {
	billStatistics := new(BillStatistics)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(billStatistics).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if billStatistics.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetBillStatisticsTotal gets the total number of billStatisticss based on the constraints
func (a *BillStatistics) GetBillStatisticsTotal(params *BillStatisticsQuery, maps interface{}) (int64, error) {
	var count int64
	query := getBillStatisticsListParams(a, params)
	err := query.Where(maps).Count(&count).Error
	return count, err
}

// GetBillStatisticss gets a list of billStatisticss based on paging constraints
func (a *BillStatistics) GetBillStatisticsList(params *BillStatisticsQuery, pageNum, pageSize int, maps interface{}) ([]*BillStatistics, error) {
	var billStatisticsList []*BillStatistics
	query := getBillStatisticsListParams(a, params)

	err := query.Where(maps).Offset(pageNum).Limit(pageSize).Find(&billStatisticsList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return billStatisticsList, nil
}

// GetBillStatistics Get a single billStatistics based on ID
func (a *BillStatistics) GetBillStatistics(id int64) (*BillStatistics, error) {
	billStatistics := new(BillStatistics)
	err := a.Session.db.Where("id = ?", id).First(billStatistics).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return billStatistics, nil
}

// EditBillStatistics modify a single billStatistics
func (a *BillStatistics) EditBillStatistics(id int64, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&BillStatistics{}).Where("id = ?", id).Updates(data).Error
}

// AddBillStatistics add a single billStatistics
func (a *BillStatistics) AddBillStatistics(billStatistics *BillStatistics) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(billStatistics).Error
}

// DeleteBillStatistics delete a single billStatistics
func (a *BillStatistics) DeleteBillStatistics(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(BillStatistics{}).Error
}

func (a *BillStatistics) GetUserByRoleBillStatistics(agency, username, startTime, endTime string, role int) (*BillStatistics, error) {
	statistics := new(BillStatistics)
	query := a.Session.db.Model(&BillStatistics{})
	if startTime != "" {
		query = query.Where("statistics_date >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("statistics_date < ?", endTime)
	}
	if username != "" {
		query = query.Where("username = ?", username)
	}
	selectSql := fmt.Sprintf("%s,%s,%s,%s,%s,%s",
		"IFNULL(SUM( accept_amount ),0) as accept_amount",
		"IFNULL(SUM( accept_count ),0) as accept_count",
		"IFNULL(SUM( withdrawal_amount ),0) as withdrawal_amount",
		"IFNULL(SUM( withdrawal_count ),0) as withdrawal_count",
		"IFNULL(SUM( recharge_amount ),0) as recharge_amount",
		"IFNULL(SUM( recharge_count ),0) as recharge_count")
	err := query.Select(selectSql).
		Where("agency = ? and role = ?", agency, role).
		First(statistics).Error
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func getBillStatisticsListParams(a *BillStatistics, params *BillStatisticsQuery) *gorm.DB {
	query := a.Session.db.Model(&BillStatistics{})
	if params.StartDate != "" {
		query = query.Where("statistics_date >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("statistics_date < ?", params.EndDate)
	}
	query = query.Order("id desc")
	return query
}

type BillStatisticsQuery struct {
	// 开始时间
	StartDate string
	// 结束时间
	EndDate string
}
