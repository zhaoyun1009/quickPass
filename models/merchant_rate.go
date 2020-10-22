package models

import (
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type MerchantRate struct {
	// 主键id
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 商家账号
	MerchantName string `json:"merchant_name" gorm:"column:merchant_name"`
	// 所属代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 通道名称
	Channel string `json:"channel" gorm:"column:channel"`
	// 通道汇率
	Rate int64 `json:"rate" gorm:"column:rate"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置Merchant的表名为`merchant`
func (MerchantRate) TableName() string {
	return "merchant_rate"
}

func NewMerchantRateModel(session *Session) *MerchantRate {
	return &MerchantRate{Session: session}
}

// ExistMerchantByID checks if an merchant exists based on agency and merchant_name
func (a *MerchantRate) ExistMerchantRateByAgencyAndUsername(agency, username string) (bool, error) {
	merchantRate := new(MerchantRate)
	err := a.Session.db.Select("id").Where("agency = ? and merchant_name = ?", agency, username).First(merchantRate).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if merchantRate.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetMerchantRateTotal gets the total number of merchantRates based on the constraints
func (a *MerchantRate) GetMerchantRateTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&MerchantRate{}).Where(maps).Count(&count).Error
	return count, err
}

// GetMerchantRates gets a list of merchants based on paging constraints
func (a *MerchantRate) GetMerchantRates(agency, username string) ([]*MerchantRate, error) {
	var merchants []*MerchantRate
	err := a.Session.db.Where("agency = ? and merchant_name = ?", agency, username).Find(&merchants).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return merchants, nil
}

// GetMerchantRate Get a single merchantRate based on ID
func (a *MerchantRate) GetMerchantRate(agency, username, channel string) *MerchantRate {
	merchantRate := new(MerchantRate)
	err := a.Session.db.Where("agency = ? and merchant_name = ? and channel = ?", agency, username, channel).First(merchantRate).Error
	if err != nil {
		return nil
	}

	return merchantRate
}

// EditAllMerchantRate modify a single merchant
func (a *MerchantRate) EditMerchantAllChannelRate(agency, username string, rate int64) error {
	tx := GetSessionTx(a.Session)
	data := map[string]interface{}{
		"rate": rate,
	}
	return tx.Model(&MerchantRate{}).
		Where("agency = ? and merchant_name = ?", agency, username).
		Updates(data).
		Error
}

// EditAllMerchantRate modify a single merchant
func (a *MerchantRate) EditMerchantSingleRate(agency, username, channel string, rate int64) error {
	tx := GetSessionTx(a.Session)
	data := map[string]interface{}{
		"rate": rate,
	}
	return tx.Model(&MerchantRate{}).
		Where("agency = ? and merchant_name = ? and channel = ?", agency, username, channel).
		Updates(data).
		Error
}

// AddMerchantRate add a single merchantRate
func (a *MerchantRate) AddMerchantRate(merchantRate *MerchantRate) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(merchantRate).Error
}

// DeleteMerchantRate delete a single merchantRate
func (a *MerchantRate) DeleteMerchantRate(agency string, username string) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("agency = ? and merchant_name = ?", agency, username).Delete(MerchantRate{}).Error
}
