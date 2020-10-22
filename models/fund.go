package models

import (
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Fund struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 用户名
	UserName string `json:"user_name" gorm:"column:user_name"`
	// 账户类型(1：系统资金账户  2：普通资金账户)
	Type int `json:"type" gorm:"column:type"`
	// 可用资金
	AvailableAmount int64 `json:"available_amount" gorm:"column:available_amount"`
	// 冻结资金
	FrozenAmount int64 `json:"frozen_amount" gorm:"column:frozen_amount"`
	// 资金版本号(乐观锁)
	Version int64 `json:"-" gorm:"column:version"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

type FundExt struct {
	Fund
	Role int `json:"-" gorm:"column:role"`
}

// 设置Fund的表名为`fund`
func (Fund) TableName() string {
	return "fund"
}

func NewFundModel(session *Session) *Fund {
	return &Fund{Session: session}
}

// ExistFundByID checks if an fund exists based on ID
func (a *Fund) ExistFundByID(id int64) (bool, error) {
	fund := new(Fund)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(fund).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if fund.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetFundTotal gets the total number of funds based on the constraints
func (a *Fund) GetFundTotal(maps interface{}) (int, error) {
	var count int
	err := a.Session.db.Model(&Fund{}).Where(maps).Count(&count).Error
	return count, err
}

// GetFunds gets a list of funds based on paging constraints
func (a *Fund) GetFunds(pageNum int, pageSize int, maps interface{}) ([]*Fund, error) {
	var funds []*Fund
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&funds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return funds, nil
}

// GetFunds gets a list of funds based on paging constraints
func (a *Fund) GetMerchantAndAcceptor() ([]*FundExt, error) {
	var funds []*FundExt
	err := a.Session.db.Table(a.TableName()).
		Select("fund.*, `user`.role as role").
		Joins("LEFT JOIN `user` ON `user`.agency = fund.agency and `user`.user_name = fund.user_name").
		Where("`user`.role in (?, ?)", constant.ACCEPTOR, constant.MERCHANT).
		Find(&funds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return funds, nil
}

// GetFund Get a single fund based on ID
func (a *Fund) GetFund(agency string, userName string) (*Fund, error) {
	fund := new(Fund)
	err := a.Session.db.Where("agency = ? and user_name = ?", agency, userName).First(fund).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return fund, nil
}

// GetFund Get a single fund based on ID
func (a *Fund) GetSuperAdminFund() (*Fund, error) {
	fund := new(Fund)
	err := a.Session.db.Where("type = ? and user_name = ?", constant.SystemAccount, constant.SuperAdministrator).First(fund).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return fund, nil
}

// AddFund add a single fund
func (a *Fund) AddFund(model *Fund) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(model).Error
}

// 冻结资金
func (a *Fund) Frozen(agency string, username string, amount int64) (fund *Fund, err error) {
	tx := GetSessionTx(a.Session)

	fund, err = a.GetFund(agency, username)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"available_amount": gorm.Expr("available_amount - ?", amount),
		"frozen_amount":    gorm.Expr("frozen_amount + ?", amount),
		"version":          gorm.Expr("version + 1"),
	}

	tx = tx.Model(&Fund{}).Where("agency = ?", agency).
		Where("user_name = ?", username).
		Where("available_amount >= ?", amount).
		Where("available_amount > 0").
		Where("version = ?", fund.Version).
		Updates(data)

	if tx.Error != nil {
		return fund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return fund, errors.New("Frozen [rows affected 0]")
	}

	fund.AvailableAmount -= amount
	fund.FrozenAmount += amount
	return fund, nil
}

// 解冻资金
func (a *Fund) UnFrozen(agency, username string, amount int64) (fund *Fund, err error) {
	tx := GetSessionTx(a.Session)

	fund, err = a.GetFund(agency, username)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"available_amount": gorm.Expr("available_amount + ?", amount),
		"frozen_amount":    gorm.Expr("frozen_amount - ?", amount),
		"version":          gorm.Expr("version + 1"),
	}

	tx = tx.Model(&Fund{}).Where("agency = ? and user_name = ? and frozen_amount >= ? and frozen_amount > 0 and version = ?", agency, username, amount, fund.Version).Updates(data)
	if tx.Error != nil {
		return fund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return fund, errors.New("UnFrozen [rows affected 0]")
	}

	fund.AvailableAmount += amount
	fund.FrozenAmount -= amount
	return fund, nil
}

// 买单放行(从冻结资金扣除amount, 转到toUser账户)
func (a *Fund) DischargeBuy(agency, formUsername, toUsername string, orderAmount, feeAmount int64) (agencyFund, fromFund, toFund *Fund, err error) {
	tx := GetSessionTx(a.Session)

	// 1.承兑人冻结资金中扣除
	fromFund, err = a.GetFund(agency, formUsername)
	if err != nil {
		return nil, nil, nil, err
	}
	if fromFund == nil {
		return nil, nil, nil, errors.New(fmt.Sprintf("agency[%s] Fund[%s] account not found", agency, formUsername))
	}

	// 2.商家获取资金和手续费的差价
	toFund, err = a.GetFund(agency, toUsername)
	if err != nil {
		return nil, fromFund, nil, err
	}
	if toFund == nil {
		return nil, fromFund, nil, errors.New(fmt.Sprintf("agency[%s] Fund[%s] account not found", agency, toUsername))
	}

	// 3.代理收取手续费
	agencyFund, err = a.GetFund(agency, agency)
	if err != nil {
		return nil, fromFund, toFund, err
	}
	if agencyFund == nil {
		return nil, fromFund, toFund, errors.New(fmt.Sprintf("agency[%s] Fund[%s] account not found", agency, agency))
	}

	now := util.JSONTimeNow()

	tx = tx.Exec("UPDATE `fund` SET `frozen_amount` = frozen_amount - ?, `update_time` = ?, `version` = version + 1 "+
		"WHERE (agency = ? and user_name = ? and frozen_amount >= ? and frozen_amount > 0 and version = ?)", orderAmount, now, agency, formUsername, orderAmount, fromFund.Version)
	if tx.Error != nil {
		return agencyFund, fromFund, toFund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return agencyFund, fromFund, toFund, errors.New("rows affected 0")
	}

	tx = tx.Exec("UPDATE `fund` SET `available_amount` = available_amount + ?, `update_time` = ?, `version` = version + 1 "+
		"WHERE (agency = ? and user_name = ? and version = ?)", orderAmount-feeAmount, now, agency, toUsername, toFund.Version)
	if tx.Error != nil {
		return agencyFund, fromFund, toFund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return agencyFund, fromFund, toFund, errors.New("rows affected 0")
	}

	if feeAmount > 0 {
		tx = tx.Exec("UPDATE `fund` SET `available_amount` = available_amount + ?, `update_time` = ?, `version` = version + 1 "+
			"WHERE (agency = ? and user_name = ? and version = ?)", feeAmount, now, agency, agency, agencyFund.Version)
		if tx.Error != nil {
			return agencyFund, fromFund, toFund, tx.Error
		}
		if !(tx.RowsAffected > 0) {
			return agencyFund, fromFund, toFund, errors.New("rows affected 0")
		}
	}

	fromFund.FrozenAmount -= orderAmount
	toFund.AvailableAmount += orderAmount - feeAmount
	agencyFund.AvailableAmount += feeAmount
	return agencyFund, fromFund, toFund, nil
}

// 卖单放行
func (a *Fund) DischargeSell(fromAgency, formUsername string, orderAmount int64) (fromFund, adminFund *Fund, err error) {
	tx := GetSessionTx(a.Session)

	// 1.商家冻结资金中扣除
	fromFund, err = a.GetFund(fromAgency, formUsername)
	if err != nil {
		return nil, nil, err
	}
	if fromFund == nil {
		return nil, nil, errors.New(fmt.Sprintf("agency[%s] Fund[%s] account not found", fromAgency, formUsername))
	}

	if fromFund.FrozenAmount <= 0 || fromFund.FrozenAmount < orderAmount {
		return nil, nil, errors.New("冻结余额不足")
	}

	// 2.系统回收资金
	// 2.1 获取系统资金账户
	adminFund, err = a.GetSuperAdminFund()
	if err != nil {
		return fromFund, nil, err
	}
	if adminFund == nil {
		return fromFund, nil, errors.New("adminFund not found")
	}

	now := util.JSONTimeNow()

	tx = tx.Exec("UPDATE `fund` SET `frozen_amount` = frozen_amount - ?, `update_time` = ?, `version` = version + 1 "+
		"WHERE (agency = ? and user_name = ? and frozen_amount >= ? and frozen_amount > 0 and version = ?)", orderAmount, now, fromAgency, formUsername, orderAmount, fromFund.Version)
	if tx.Error != nil {
		return fromFund, adminFund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return fromFund, adminFund, errors.New("rows affected 0")
	}

	tx = tx.Exec("UPDATE `fund` SET `available_amount` = available_amount + ?, `update_time` = ?, `version` = version + 1 "+
		"WHERE (agency = ? and user_name = ? and version = ?)", orderAmount, now, "", adminFund.UserName, adminFund.Version)
	if tx.Error != nil {
		return fromFund, adminFund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return fromFund, adminFund, errors.New("rows affected 0")
	}

	fromFund.FrozenAmount -= orderAmount
	adminFund.AvailableAmount += orderAmount
	return fromFund, adminFund, nil
}

// 转账交易
func (a *Fund) Transaction(fromAgency, toAgency, formUsername, toUsername string, amount int64) (fromFund, toFund *Fund, err error) {
	tx := GetSessionTx(a.Session)

	// 1.获取转账方资金账户
	fromFund, err = a.GetFund(fromAgency, formUsername)
	if err != nil {
		return nil, nil, err
	}
	if fromFund == nil {
		return nil, nil, errors.New("fund account not found")
	}
	// 2.判断转账方可用余额
	if fromFund.AvailableAmount <= 0 || fromFund.AvailableAmount < amount {
		return fromFund, nil, errors.New("可用余额不足")
	}

	// 3.获取收款方资金账户
	toFund, err = a.GetFund(toAgency, toUsername)
	if err != nil {
		return fromFund, nil, err
	}
	if toFund == nil {
		return fromFund, nil, errors.New("toFund account not found")
	}

	now := util.JSONTimeNow()
	// 4.交易
	//更新转出账户资金
	tx = tx.Exec("UPDATE `fund` SET `available_amount` = available_amount - ?, `update_time` = ?, `version` = version + 1 "+
		" WHERE (agency = ? and user_name = ? and available_amount >= ? and available_amount > 0 and version = ?)", amount, now, fromAgency, formUsername, amount, fromFund.Version)
	if tx.Error != nil {
		return fromFund, toFund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		// 有内鬼，终止交易
		return fromFund, toFund, errors.New("rows affected 0")
	}

	//更新转入账户资金
	tx = tx.Exec("UPDATE `fund` SET `available_amount` = available_amount + ?, `update_time` = ?, `version` = version + 1  "+
		"WHERE (agency = ? and user_name = ? and version = ?)", amount, now, toAgency, toUsername, toFund.Version)
	if tx.Error != nil {
		return fromFund, toFund, tx.Error
	}
	if !(tx.RowsAffected > 0) {
		// 有内鬼，终止交易
		return fromFund, toFund, errors.New("rows affected 0")
	}

	fromFund.AvailableAmount -= amount
	toFund.AvailableAmount += amount
	return fromFund, toFund, nil
}

//删除资金账户
func (a *Fund) Remove(agency, username string) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("agency = ? and user_name = ?", agency, username).Delete(Fund{}).Error
}
