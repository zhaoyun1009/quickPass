package models

import (
	"QuickPass/pkg/constant"
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type AcceptorCard struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 是否开启={1:关闭，2：开启}
	IfOpen int `json:"if_open" gorm:"column:if_open"`
	// 承兑人账号
	Acceptor string `json:"acceptor" gorm:"column:acceptor"`
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
	// 当日剩余可用流水
	DayAvailableAmt int64 `json:"day_available_amt" gorm:"column:day_available_amt"`
	// 单日冻结流水
	DayFrozenAmount int64 `json:"day_frozen_amount" gorm:"column:day_frozen_amount"`
	// 当日最大流水
	DayMaxAmt int64 `json:"day_max_amt" gorm:"column:day_max_amt"`
	// 删除标记（1：未删除，2：已删除）
	DeleteFlag int `json:"delete_flag" gorm:"column:delete_flag"`
	// 上一次的匹配时间
	LastMatchTime util.JSONTime `json:"last_match_time" gorm:"column:last_match_time"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置AcceptorCard的表名为`acceptor_card`
func (AcceptorCard) TableName() string {
	return "acceptor_card"
}

func NewAcceptorCardModel(session *Session) *AcceptorCard {
	return &AcceptorCard{Session: session}
}

// GetAcceptorCardTotal gets the total number of acceptorCards based on the constraints
func (a *AcceptorCard) GetAcceptorCardTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&AcceptorCard{}).Where(maps).Count(&count).Error
	return count, err
}

// GetAcceptorCards gets a list of acceptorCards based on paging constraints
func (a *AcceptorCard) GetAcceptorCards(pageNum int, pageSize int, maps interface{}) ([]*AcceptorCard, error) {
	var acceptorCards []*AcceptorCard
	err := a.Session.db.Where(maps).
		Where("delete_flag = ?", constant.NotDeleted).
		Offset(pageNum).
		Limit(pageSize).
		Order("id desc").
		Find(&acceptorCards).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return acceptorCards, nil
}

// GetAcceptorCards gets a list of acceptorCards based on paging constraints
func (a *AcceptorCard) GetAcceptorAllCards(agency, acceptor string) (int64, []*AcceptorCard, error) {
	var total int64
	var acceptorCards []*AcceptorCard

	err := a.Session.db.Model(&AcceptorCard{}).Where("agency = ? and acceptor = ? and delete_flag = ?", agency, acceptor, constant.NotDeleted).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}
	err = a.Session.db.Where("agency = ? and acceptor = ? and delete_flag = ?", agency, acceptor, constant.NotDeleted).Find(&acceptorCards).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return total, nil, err
	}

	return total, acceptorCards, nil
}

// GetAcceptorCard Get a single acceptorCard based on ID
func (a *AcceptorCard) GetAcceptorCard(id int64) (*AcceptorCard, error) {
	acceptorCard := new(AcceptorCard)
	err := a.Session.db.Where("id = ? and delete_flag = ?", id, constant.NotDeleted).First(acceptorCard).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return acceptorCard, nil
}

// GetAcceptorCard Get a single acceptorCard based on ID
func (a *AcceptorCard) GetAcceptorCardByCardNo(agency, acceptor, cardNo, cardType string) (*AcceptorCard, error) {
	acceptorCard := new(AcceptorCard)
	err := a.Session.db.
		Where("agency = ? ", agency).
		Where("acceptor = ? ", acceptor).
		Where("card_type = ? ", cardType).
		Where("card_no = ? ", cardNo).
		Order("create_time desc").
		First(acceptorCard).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return acceptorCard, nil
}

// ExistAcceptorCardByID checks if an acceptorCard exists based on ID
func (a *AcceptorCard) getAcceptorCardByFlag(agency, acceptor, cardType, cardNo string, deleteFlag int) (*AcceptorCard, error) {
	acceptorCard := new(AcceptorCard)
	err := a.Session.db.
		Where("agency = ? ", agency).
		Where("acceptor = ? ", acceptor).
		Where("card_type = ? ", cardType).
		Where("card_no = ? ", cardNo).
		Where("delete_flag = ? ", deleteFlag).
		First(acceptorCard).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return acceptorCard, nil
}

// 获取可匹配的承兑人卡
// 条件1：卡的剩余可匹配额度 >= 下单金额
// 条件2：承兑人资金账户可用余额 >= 下单金额
func (a *AcceptorCard) GetMatchedCard(agency string, amount int64, cardType string) (*AcceptorCard, error) {
	card := new(AcceptorCard)

	err := a.Session.db.
		Table("acceptor_card").
		Select("acceptor_card.*").
		Joins("LEFT JOIN fund AS f ON acceptor_card.agency = f.agency and f.user_name = acceptor_card.acceptor").
		Where("acceptor_card.if_open = ?", constant.SwitchOpen).
		Where("acceptor_card.delete_flag = ?", constant.NotDeleted).
		Where("acceptor_card.agency = ?", agency).
		Where("acceptor_card.card_type = ?", cardType).
		Where("acceptor_card.day_available_amt >= ?", amount).
		Where("f.available_amount >= ?", amount).
		Where("acceptor_card.card_no not in "+
			"(select o.card_no from `order` as o where agency = ? AND o.card_no != '' AND o.order_status IN ( ?, ? ) AND o.order_type = ? AND o.amount = ?)", agency, constant.ORDER_STATUS_WAIT_PAY, constant.ORDER_STATUS_WAIT_DISCHARGE, constant.ORDER_TYPE_BUY, amount).
		Where("acceptor_card.acceptor in (select acceptor from `acceptor` where agency = ? and accept_switch = ?)", agency, constant.SwitchOpen).
		Order("acceptor_card.last_match_time").First(card).Error
	if err != nil {
		return nil, err
	}

	return card, nil
}

//更新缓存卡最大匹配金额
func (a *AcceptorCard) AddDayAvailableAmount(agency, username, cardType, cardNo string, matchedAmount int64) error {
	tx := GetSessionTx(a.Session)

	data := map[string]interface{}{
		"day_available_amt": gorm.Expr("day_available_amt + ?", matchedAmount),
		"day_frozen_amount": gorm.Expr("day_frozen_amount - ?", matchedAmount),
	}
	return tx.Model(&AcceptorCard{}).Where("agency = ? and acceptor = ? and card_no = ? and card_type = ?", agency, username, cardNo, cardType).Updates(data).Error
}

//更新缓存卡最大匹配金额
func (a *AcceptorCard) SubtractDayAvailableAmount(id int64, matchedAmount int64) error {
	tx := GetSessionTx(a.Session)

	data := map[string]interface{}{
		"day_available_amt": gorm.Expr("day_available_amt - ?", matchedAmount),
		"day_frozen_amount": gorm.Expr("day_frozen_amount + ?", matchedAmount),
		"last_match_time":   util.JSONTimeNow(),
	}
	return tx.Model(&AcceptorCard{}).Where("id = ?", id).Updates(data).Error
}

//更新缓存卡冻结金额
func (a *AcceptorCard) SubtractDayFrozenAmount(agency, username, cardType, cardNo string, matchedAmount int64) error {
	tx := GetSessionTx(a.Session)

	data := map[string]interface{}{
		"day_frozen_amount": gorm.Expr("day_frozen_amount - ?", matchedAmount),
	}
	return tx.Model(&AcceptorCard{}).Where("agency = ? and acceptor = ? and card_no = ? and card_type = ?", agency, username, cardNo, cardType).Updates(data).Error
}

//调度器用
//重置当日剩余可用流水
func (a *AcceptorCard) ResetLimitAvailableAmount(agency, channel string, amount int64) error {
	tx := GetSessionTx(a.Session)

	data := map[string]interface{}{
		"day_available_amt": amount,
		"day_frozen_amount": 0,
		"day_max_amt":       amount,
	}
	return tx.Model(&AcceptorCard{}).Where("agency = ? and card_type = ?", agency, channel).Updates(data).Error
}

// EditAcceptorCard modify a single acceptorCard
func (a *AcceptorCard) EditAcceptorCard(id int64, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&AcceptorCard{}).Where("id = ?", id).Updates(data).Error
}

// AddAcceptorCard add a single acceptorCard
func (a *AcceptorCard) AddAcceptorCard(acceptorCard *AcceptorCard) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(acceptorCard).Error
}

// DeleteAcceptorCard delete a single acceptorCard
// 逻辑删除
func (a *AcceptorCard) DeleteAcceptorCard(id int64) error {
	tx := GetSessionTx(a.Session)
	// 获取卡信息
	card, err := a.GetAcceptorCard(id)
	if err != nil {
		return err
	}

	// 判断有没有被删除的相同的卡信息,有则先删除，再改变原来卡的状态
	deleteCard, err := a.getAcceptorCardByFlag(card.Agency, card.Acceptor, card.CardType, card.CardNo, constant.Deleted)
	if err != nil {
		return err
	}
	if deleteCard != nil {
		err := a.deleteFromDb(deleteCard.Id)
		if err != nil {
			return err
		}
	}

	return tx.Model(&AcceptorCard{}).Where("id = ?", id).Updates(map[string]interface{}{
		"delete_flag": constant.Deleted,
	}).Error
}

// 物理删除
func (a *AcceptorCard) deleteFromDb(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(AcceptorCard{}).Error
}

type OnlineAcceptorCardFundInfo struct {
	Agency               string
	Acceptor             string
	CardType             string
	CardId               int64
	TotalAmount          int64
	FrozenAmount         int64
	FundVersion          int64
	Deposit              int64
	SignalMaxAmt         int64
	MaxAcceptedAmount    int64
	MinAcceptedAmount    int64
	AcceptedMerchantList string
}

// 获取承兑人卡代理组
func (a *AcceptorCard) GetAcceptorCardAgencyGroup() ([]string, error) {
	tx := GetSessionTx(a.Session)

	var acceptorCards []*AcceptorCard

	err := tx.Table(a.TableName()).Select("agency").Group("agency").Find(&acceptorCards).Error
	if err != nil {
		return nil, err
	}

	group := make([]string, 0, len(acceptorCards))
	for _, item := range acceptorCards {
		group = append(group, item.Agency)
	}

	return group, nil
}

type CountCardType struct {
	// 银行卡数量
	CountBankCard int64 `json:"count_bank_card" gorm:"column:count_bank_card"`
	// 支付宝数量
	CountAliPay int64 `json:"count_ali_pay" gorm:"column:count_ali_pay"`
	// 微信数量
	CountWeChat int64 `json:"count_we_chat" gorm:"column:count_we_chat"`
}

func (a *AcceptorCard) CountCardType(agency string, acceptor string) (*CountCardType, error) {
	tx := GetSessionTx(a.Session)

	count := new(CountCardType)
	selectArr := []string{
		"SUM(CASE WHEN card_type = ? then 1 else 0 end) as count_bank_card",
		"SUM(CASE WHEN card_type = ? then 1 else 0 end) as count_ali_pay",
		"SUM(CASE WHEN card_type = ? then 1 else 0 end) as count_we_chat",
	}
	err := tx.Table(a.TableName()).
		Select(selectArr, constant.BANK_CARD, constant.ALIPAY, constant.WECHAT).
		Where("agency = ? and acceptor = ? and delete_flag = ?", agency, acceptor, constant.NotDeleted).First(count).Error

	return count, err
}
