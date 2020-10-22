package models

import (
	"QuickPass/pkg/constant"
	"QuickPass/pkg/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

type AbnormalOrder struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 用户名
	Username string `json:"username" gorm:"column:username"`
	// 昵称
	Nickname string `json:"nickname" gorm:"column:nickname"`
	// 收款通道({“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"})
	Channel string `json:"channel" gorm:"column:channel"`
	// 原始订单号
	OriginalOrderNo string `json:"original_order_no" gorm:"column:original_order_no"`
	// 异常单号
	AbnormalOrderNo string `json:"abnormal_order_no" gorm:"column:abnormal_order_no"`
	// 异常类型（1: 未知, 2：超时取消，3：未生成订单， 4：订单金额不符）
	AbnormalOrderType int64 `json:"abnormal_order_type" gorm:"column:abnormal_order_type"`
	// 订单状态（1：未处理  2：已处理  3：处理中）
	AbnormalOrderStatus int64 `json:"abnormal_order_status" gorm:"column:abnormal_order_status"`
	// 收款账户
	AcceptCardAccount string `json:"accept_card_account" gorm:"column:accept_card_account"`
	// 收款卡号
	AcceptCardNo string `json:"accept_card_no" gorm:"column:accept_card_no"`
	// 银行名称
	AcceptCardBank string `json:"accept_card_bank" gorm:"column:accept_card_bank"`
	// 金额
	Amount int64 `json:"amount" gorm:"column:amount"`
	// 摘要
	AppendInfo string `json:"append_info" gorm:"column:append_info"`
	// 付款日期
	PaymentDate util.JSONTime `json:"payment_date" gorm:"column:payment_date"`
	// 完成时间
	FinishTime *util.JSONTime `json:"finish_time" gorm:"column:finish_time"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置AbnormalOrder的表名为`abnormal_order`
func (AbnormalOrder) TableName() string {
	return "abnormal_order"
}

func NewAbnormalOrderModel(session *Session) *AbnormalOrder {
	return &AbnormalOrder{Session: session}
}

// ExistAbnormalOrderByID checks if an abnormalOrder exists based on ID
func (a *AbnormalOrder) ExistAbnormalOrderByID(id int64) (bool, error) {
	abnormalOrder := new(AbnormalOrder)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(abnormalOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if abnormalOrder.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetAbnormalOrderTotal gets the total number of abnormalOrders based on the constraints
func (a *AbnormalOrder) GetAbnormalOrderTotal(params *AbnormalOrderListQueryParams, maps interface{}) (int64, error) {
	var count int64
	query := getAbnormalOrderListParams(a, params)
	err := query.Where(maps).Count(&count).Error
	return count, err
}

// GetAbnormalOrders gets a list of abnormalOrders based on paging constraints
func (a *AbnormalOrder) GetAbnormalOrders(params *AbnormalOrderListQueryParams, pageNum int, pageSize int, maps interface{}) ([]*AbnormalOrder, error) {
	var abnormalOrders []*AbnormalOrder
	query := getAbnormalOrderListParams(a, params)
	err := query.Where(maps).Offset(pageNum).Limit(pageSize).Find(&abnormalOrders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return abnormalOrders, nil
}

func getAbnormalOrderListParams(a *AbnormalOrder, params *AbnormalOrderListQueryParams) *gorm.DB {
	query := a.Session.db.Model(&AbnormalOrder{})
	if params.StartTime != "" {
		query = query.Where("create_time >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		query = query.Where("create_time <= ?", params.EndTime)
	}
	if params.FinishStartTime != "" {
		query = query.Where("finish_time >= ?", params.FinishStartTime)
	}
	if params.FinishEndTime != "" {
		query = query.Where("finish_time <= ?", params.FinishEndTime)
	}
	if params.MinAmount != "" {
		query = query.Where("amount >= ?", params.MinAmount)
	}
	if params.MaxAmount != "" {
		query = query.Where("amount <= ?", params.MaxAmount)
	}
	orderBy := fmt.Sprintf("FIELD(abnormal_order_status, %d, %d, %d), create_time desc", constant.Unprocessed, constant.Processing, constant.Processed)
	query = query.Order(orderBy)
	return query
}

// GetAbnormalOrder Get a single abnormalOrder based on ID
func (a *AbnormalOrder) GetAbnormalOrder(agency, abnormalOrderNo string) (*AbnormalOrder, error) {
	abnormalOrder := new(AbnormalOrder)
	err := a.Session.db.Where("agency = ? and abnormal_order_no = ?", agency, abnormalOrderNo).First(abnormalOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return abnormalOrder, nil
}

// 通过原订单查询异常单
func (a *AbnormalOrder) GetAbnormalOrderByOriginalOrderNo(agency, originalOrderNo string) (*AbnormalOrder, error) {
	abnormalOrder := new(AbnormalOrder)
	err := a.Session.db.Where("agency = ? and original_order_no = ?", agency, originalOrderNo).First(abnormalOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return abnormalOrder, nil
}

// EditAbnormalOrder modify a single abnormalOrder
func (a *AbnormalOrder) EditAbnormalOrder(agency, abnormalOrderNo string, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&AbnormalOrder{}).Where("agency = ? and abnormal_order_no = ?", agency, abnormalOrderNo).Updates(data).Error
}

// AddAbnormalOrder add a single abnormalOrder
func (a *AbnormalOrder) AddAbnormalOrder(abnormalOrder *AbnormalOrder) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(abnormalOrder).Error
}

// DeleteAbnormalOrder delete a single abnormalOrder
func (a *AbnormalOrder) DeleteAbnormalOrder(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(AbnormalOrder{}).Error
}

type AbnormalOrderListQueryParams struct {
	// 开始时间
	StartTime string
	// 结束时间
	EndTime string
	// 最大金额
	MaxAmount string
	// 最小金额
	MinAmount string
	// 完成开始时间
	FinishStartTime string
	// 完成结束时间
	FinishEndTime string
}

func (a *AbnormalOrder) GetAcceptorGroup(agency string) ([]string, error) {
	var orders []*AbnormalOrder
	err := a.Session.db.Table(a.TableName()).
		Select("nickname").
		Where("nickname != '' and agency = ?", agency).
		Group("nickname").Find(&orders).Error

	acceptorGroup := make([]string, 0, len(orders))
	for _, item := range orders {
		acceptorGroup = append(acceptorGroup, item.Nickname)
	}
	return acceptorGroup, err
}
