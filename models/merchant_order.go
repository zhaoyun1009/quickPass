package models

import (
	"QuickPass/pkg/constant"
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type MerchantOrder struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 商家账户名
	Username string `json:"username" gorm:"column:username"`
	// 系统订单号
	OrderNo string `json:"order_no" gorm:"column:order_no"`
	// 商家平台传过来的订单号
	MerchantOrderNo string `json:"merchant_order_no" gorm:"column:merchant_order_no"`
	// 下单方式（1、接口，2、非接口）
	SubmitType int `json:"submit_type" gorm:"column:submit_type"`
	// 支付成功的返回URL
	ReturnUrl string `json:"return_url" gorm:"column:return_url"`
	// 回调接口地址
	CallbackUrl string `json:"callback_url" gorm:"column:callback_url"`
	// 回调结果（1、未处理，2、回调成功，3、回调失败）
	CallbackStatus int `json:"callback_status" gorm:"column:callback_status"`
	// 回调次数
	CallbackTimes int `json:"callback_times" gorm:"column:callback_times"`
	// 回调失败的信息
	CallbackInfo string `json:"callback_info" gorm:"column:callback_info"`
	// 回调请求状态（1、未请求，2、请求中）
	CallbackLock int `json:"callback_lock" gorm:"column:callback_lock"`
	// 备注信息
	AppendInfo string `json:"append_info" gorm:"column:append_info"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置MerchantOrder的表名为`merchant_order`
func (MerchantOrder) TableName() string {
	return "merchant_order"
}

func NewMerchantOrderModel(session *Session) *MerchantOrder {
	return &MerchantOrder{Session: session}
}

// ExistMerchantOrderByID checks if an merchantOrder exists based on ID
func (a *MerchantOrder) ExistMerchantOrderByID(id int64) (bool, error) {
	merchantOrder := new(MerchantOrder)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(merchantOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if merchantOrder.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetMerchantOrderTotal gets the total number of merchantOrders based on the constraints
func (a *MerchantOrder) GetMerchantOrderTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&MerchantOrder{}).Where(maps).Count(&count).Error
	return count, err
}

// GetMerchantOrders gets a list of merchantOrders based on paging constraints
func (a *MerchantOrder) GetMerchantOrders(pageNum int, pageSize int, maps interface{}) ([]*MerchantOrder, error) {
	var merchantOrders []*MerchantOrder
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&merchantOrders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return merchantOrders, nil
}

// GetMerchantOrder Get a single merchantOrder based on ID
func (a *MerchantOrder) GetMerchantOrder(agency, orderNo string) (*MerchantOrder, error) {
	merchantOrder := new(MerchantOrder)
	err := a.Session.db.Where("agency = ? and order_no = ?", agency, orderNo).First(merchantOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return merchantOrder, nil
}

// GetMerchantOrder Get a mul merchantOrder based on ID
func (a *MerchantOrder) GetMerchantOrderByMerchantOrderNo(agency, merchantOrderNo string) ([]*MerchantOrder, error) {
	var merchantOrders []*MerchantOrder
	err := a.Session.db.Where("agency = ? and merchant_order_no = ?", agency, merchantOrderNo).Find(&merchantOrders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return merchantOrders, nil
}

// GetMerchantOrder Get a single merchantOrder based on ID
func (a *MerchantOrder) GetOneMerchantOrderByMerchantOrderNo(agency, merchantName, merchantOrderNo string) (*MerchantOrder, error) {
	merchantOrder := new(MerchantOrder)
	err := a.Session.db.Where("agency = ? and username = ? and merchant_order_no = ?", agency, merchantName, merchantOrderNo).First(merchantOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return merchantOrder, nil
}

// EditMerchantOrder modify a single merchantOrder
func (a *MerchantOrder) EditMerchantOrder(id int64, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&MerchantOrder{}).Where("id = ?", id).Updates(data).Error
}

// AddMerchantOrder add a single merchantOrder
func (a *MerchantOrder) AddMerchantOrder(merchantOrder *MerchantOrder) error {
	tx := GetSessionTx(a.Session)
	merchantOrder.CallbackLock = constant.CallbackUnLock
	return tx.Create(merchantOrder).Error
}

// DeleteMerchantOrder delete a single merchantOrder
func (a *MerchantOrder) DeleteMerchantOrder(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(MerchantOrder{}).Error
}
