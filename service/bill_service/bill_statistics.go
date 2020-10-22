package bill_service

import (
	"QuickPass/models"
	"QuickPass/pkg/util"
)

type BillStatistics struct {
	//
	Id int64
	// 代理
	Agency string
	// 用户名
	Username string
	// 角色
	Role int
	// 统计日期
	StatisticsDate string
	// 剩余金额
	LeftAmount int64
	// 承兑金额
	AcceptAmount int64
	// 承兑次数
	AcceptCount int64
	// 提现金额
	WithdrawalAmount int64
	// 提现次数
	WithdrawalCount int64
	// 充值金额
	RechargeAmount int64
	// 充值次数
	RechargeCount int64
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int

	//-------------------
	// 开始时间
	StartDate string
	// 结束时间
	EndDate string
}

func (a *BillStatistics) Add() error {
	billStatistics := &models.BillStatistics{
		Agency:           a.Agency,
		Username:         a.Username,
		Role:             a.Role,
		StatisticsDate:   a.StatisticsDate,
		LeftAmount:       a.LeftAmount,
		AcceptAmount:     a.AcceptAmount,
		AcceptCount:      a.AcceptCount,
		WithdrawalAmount: a.WithdrawalAmount,
		WithdrawalCount:  a.WithdrawalCount,
		RechargeAmount:   a.RechargeAmount,
		RechargeCount:    a.RechargeCount,
	}

	session := models.NewSession()
	if err := models.NewBillStatisticsModel(session).AddBillStatistics(billStatistics); err != nil {
		return err
	}

	return nil
}

func (a *BillStatistics) Edit() error {
	session := models.NewSession()
	return models.NewBillStatisticsModel(session).EditBillStatistics(a.Id, map[string]interface{}{
		"agency":            a.Agency,
		"username":          a.Username,
		"role":              a.Role,
		"statistics_date":   a.StatisticsDate,
		"left_amount":       a.LeftAmount,
		"accept_amount":     a.AcceptAmount,
		"accept_count":      a.AcceptCount,
		"withdrawal_amount": a.WithdrawalAmount,
		"withdrawal_count":  a.WithdrawalCount,
		"recharge_amount":   a.RechargeAmount,
		"recharge_count":    a.RechargeCount,
	})
}

func (a *BillStatistics) Get() (*models.BillStatistics, error) {
	session := models.NewSession()
	return models.NewBillStatisticsModel(session).GetBillStatistics(a.Id)
}

func (a *BillStatistics) GetAll() ([]*models.BillStatistics, error) {
	session := models.NewSession()
	return models.NewBillStatisticsModel(session).GetBillStatisticsList(&models.BillStatisticsQuery{
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
	}, a.PageNum, a.PageSize, a.getMaps())
}

func (a *BillStatistics) Delete() error {
	session := models.NewSession()
	return models.NewBillStatisticsModel(session).DeleteBillStatistics(a.Id)
}

func (a *BillStatistics) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewBillStatisticsModel(session).ExistBillStatisticsByID(a.Id)
}

func (a *BillStatistics) Count() (int64, error) {
	session := models.NewSession()
	return models.NewBillStatisticsModel(session).GetBillStatisticsTotal(&models.BillStatisticsQuery{
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
	}, a.getMaps())
}

func (a *BillStatistics) GetUserByRoleBillStatistics(agency, username, startTime, endTime string, role int) (*models.BillStatistics, error) {
	session := models.NewSession()
	return models.NewBillStatisticsModel(session).GetUserByRoleBillStatistics(agency, username, startTime, endTime, role)
}

func (a *BillStatistics) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.Username != "" {
		maps["username"] = a.Username
	}
	if a.StatisticsDate != "" {
		maps["statistics_date"] = a.StatisticsDate
	}

	return maps
}
