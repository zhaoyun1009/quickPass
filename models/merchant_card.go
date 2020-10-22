package models

import (
	"QuickPass/pkg/errors"
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type MerchantCard struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 商家账号
	Merchant string `json:"merchant" gorm:"column:merchant"`
	// 卡类型={“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"}
	CardType string `json:"card_type" gorm:"column:card_type"`
	// 卡号
	CardNo string `json:"card_no" gorm:"column:card_no"`
	// 卡账户名
	CardAccount string `json:"card_account" gorm:"column:card_account"`
	// 银行名称
	CardBank string `json:"card_bank" gorm:"column:card_bank"`
	// 银行支行
	CardSubBank string `json:"card_sub_bank" gorm:"column:card_sub_bank"`
	// 图片地址
	CardImg string `json:"card_img" gorm:"column:card_img"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置MerchantCard的表名为`merchant_card`
func (MerchantCard) TableName() string {
	return "merchant_card"
}

func NewMerchantCardModel(session *Session) *MerchantCard {
	return &MerchantCard{Session: session}
}

// ExistMerchantCardByID checks if an merchantCard exists based on ID
func (a *MerchantCard) ExistMerchantCard(agency, merchant, cardNo, cardType string) (bool, error) {
	merchantCard := new(MerchantCard)
	err := a.Session.db.Select("id").
		Where("agency = ? and merchant = ? and card_no = ? and card_type = ?", agency, merchant, cardNo, cardType).
		First(merchantCard).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if merchantCard.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetMerchantCardTotal gets the total number of merchantCards based on the constraints
func (a *MerchantCard) GetMerchantCardTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&MerchantCard{}).Where(maps).Count(&count).Error
	return count, err
}

// GetMerchantCards gets a list of merchantCards based on paging constraints
func (a *MerchantCard) GetMerchantCards(pageNum int, pageSize int, maps interface{}) ([]*MerchantCard, error) {
	var merchantCards []*MerchantCard
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&merchantCards).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return merchantCards, nil
}

// GetMerchantCard Get a single merchantCard based on ID
func (a *MerchantCard) GetMerchantCard(id int64) (*MerchantCard, error) {
	merchantCard := new(MerchantCard)
	err := a.Session.db.Where("id = ?", id).First(merchantCard).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return merchantCard, nil
}

// EditMerchantCard modify a single merchantCard
func (a *MerchantCard) EditMerchantCard(id int64, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&MerchantCard{}).Where("id = ?", id).Updates(data).Error
}

// AddMerchantCard add a single merchantCard
func (a *MerchantCard) AddMerchantCard(merchantCard *MerchantCard) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(merchantCard).Error
}

// DeleteMerchantCard delete a single merchantCard
func (a *MerchantCard) DeleteAllMerchantCardByAgencyAndName(agency, name string) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("agency = ? and merchant = ?", agency, name).Delete(MerchantCard{}).Error
}

// DeleteMerchantCard delete a single merchantCard
func (a *MerchantCard) DeleteMerchantCardByAgencyAndName(id int64, agency, name string) error {
	tx := GetSessionTx(a.Session)
	tx = tx.Where("id = ? and agency = ? and merchant = ?", id, agency, name).Delete(MerchantCard{})
	if !(tx.RowsAffected > 0) {
		return errors.New("rows affected 0")
	}

	return tx.Error
}
