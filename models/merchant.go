package models

import (
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type Merchant struct {
	// 主键id
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 商家账号
	MerchantName string `json:"merchant_name" gorm:"column:merchant_name"`
	// 所属代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 商家接口返回url
	ReturnUrl string `json:"-" gorm:"return_url"`
	// 商家接口回调url
	NotifyUrl string `json:"-" gorm:"notify_url"`
	// 商家公钥
	MerchantPublicKey string `json:"-" gorm:"column:merchant_public_key"`
	// 商家私钥
	MerchantPrivateKey string `json:"-" gorm:"column:merchant_private_key"`
	// 系统公钥
	SystemPublicKey string `json:"-" gorm:"column:system_public_key"`
	// 系统私钥
	SystemPrivateKey string `json:"-" gorm:"column:system_private_key"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`

	//--------------------------------------
	// 总资金
	TotalAmount int64 `json:"total_amount" gorm:"-"`
	// 冻结资金
	FrozenAmount int64 `json:"frozen_amount" gorm:"-"`
	// 今日承兑金额
	TodayAcceptAmount int64 `json:"today_accept_amount" gorm:"-"`
	// 昵称
	Nickname string `json:"nickname" gorm:"-"`
	// 商家汇率列表
	Rate []*MerchantRate `json:"rate" gorm:"-"`
}

// 设置Merchant的表名为`merchant`
func (Merchant) TableName() string {
	return "merchant"
}

func NewMerchantModel(session *Session) *Merchant {
	return &Merchant{Session: session}
}

// ExistMerchantByID checks if an merchant exists based on agency and merchant_name
func (a *Merchant) ExistMerchantByAgencyAndUsername(agency, username string) (bool, error) {
	merchant := new(Merchant)
	err := a.Session.db.Select("id").Where("agency = ? and merchant_name = ?", agency, username).First(merchant).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if merchant.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetMerchantTotal gets the total number of merchants based on the constraints
func (a *Merchant) GetMerchantTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&Merchant{}).Where(maps).Count(&count).Error
	return count, err
}

// GetMerchants gets a list of merchants based on paging constraints
func (a *Merchant) GetMerchants(pageNum int, pageSize int, maps interface{}) ([]*Merchant, error) {
	var merchants []*Merchant
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&merchants).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return merchants, nil
}

// GetMerchant Get a single merchant based on ID
func (a *Merchant) GetMerchant(agency, username string) (*Merchant, error) {
	merchant := new(Merchant)
	err := a.Session.db.Where("agency = ? and merchant_name = ?", agency, username).First(merchant).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return merchant, nil
}

// EditMerchant modify a single merchant
func (a *Merchant) EditMerchant(agency, username string, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&Merchant{}).Where("agency = ? and merchant_name = ?", agency, username).Updates(data).Error
}

// AddMerchant add a single merchant
func (a *Merchant) AddMerchant(merchant *Merchant) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(merchant).Error
}

// DeleteMerchant delete a single merchant
func (a *Merchant) DeleteMerchant(agency string, username string) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("agency = ? and merchant_name = ?", agency, username).Delete(Merchant{}).Error
}
