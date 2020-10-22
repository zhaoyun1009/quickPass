package models

import (
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type Acceptor struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 承兑人账号
	Acceptor string `json:"acceptor" gorm:"column:acceptor"`
	// 所属代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 是否自动承兑（1：不自动，2：自动）
	IfAutoAccept int `json:"if_auto_accept" gorm:"column:if_auto_accept"`
	// 承兑开关(1：关闭，2：开启)
	AcceptSwitch int `json:"accept_switch" gorm:"column:accept_switch"`
	// 承兑账户状态(1：关闭，2：开启)
	AcceptStatus int `json:"accept_status" gorm:"column:accept_status"`
	// 保证金
	Deposit int64 `json:"deposit" gorm:"column:deposit"`
	// 最大承兑金额
	MaxAcceptedAmount int64 `json:"max_accepted_amount" gorm:"column:max_accepted_amount"`
	// 最小承兑金额
	MinAcceptedAmount int64 `json:"min_accepted_amount" gorm:"column:min_accepted_amount"`
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
	// 银行卡数量
	CountBankCard int64 `json:"count_bank_card" gorm:"-"`
	// 支付宝数量
	CountAliPay int64 `json:"count_ali_pay" gorm:"-"`
	// 微信数量
	CountWeChat int64 `json:"count_we_chat" gorm:"-"`
	// 承兑人昵称
	Nickname string `json:"nickname" gorm:"-"`
}

// 设置Acceptor的表名为`acceptor`
func (Acceptor) TableName() string {
	return "acceptor"
}

func NewAcceptorModel(session *Session) *Acceptor {
	return &Acceptor{Session: session}
}

// ExistAcceptorByID checks if an acceptor exists based on ID
func (a *Acceptor) ExistAcceptorByID(id int64) (bool, error) {
	acceptor := new(Acceptor)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(acceptor).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if acceptor.Id > 0 {
		return true, nil
	}

	return false, nil
}

// ExistAcceptorByAgencyAndUsername checks if an acceptor exists based on Agency and Username
func (a *Acceptor) ExistAcceptorByAgencyAndUsername(agency string, username string) (bool, error) {
	acceptor := new(Acceptor)
	err := a.Session.db.Select("id").Where("agency = ? and acceptor = ?", agency, username).First(acceptor).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if acceptor.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetAcceptorTotal gets the total number of acceptors based on the constraints
func (a *Acceptor) GetAcceptorTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&Acceptor{}).Where(maps).Count(&count).Error
	return count, err
}

// GetAcceptors gets a list of acceptors based on paging constraints
func (a *Acceptor) GetAcceptors(pageNum int, pageSize int, maps interface{}) ([]*Acceptor, error) {
	var acceptors []*Acceptor
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&acceptors).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return acceptors, nil
}

// GetAcceptor Get a single acceptor based on ID
func (a *Acceptor) GetAcceptor(id int64) (*Acceptor, error) {
	var acceptor Acceptor
	err := a.Session.db.Where("id = ?", id).First(acceptor).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &acceptor, nil
}

func (a *Acceptor) GetAcceptorByAgencyAndUsername(agency string, username string) (*Acceptor, error) {
	acceptor := new(Acceptor)
	err := db.Where("acceptor = ? and agency = ?", username, agency).First(acceptor).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return acceptor, nil
}

// EditAcceptor modify a single acceptor
func (a *Acceptor) EditAcceptor(agency, acceptor string, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&Acceptor{}).Where("agency = ? and acceptor = ?", agency, acceptor).Updates(data).Error
}

// AddAcceptor add a single acceptor
func (a *Acceptor) AddAcceptor(acceptor *Acceptor) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(&acceptor).Error
}

// DeleteAcceptor delete a single acceptor
func (a *Acceptor) DeleteAcceptor(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(Acceptor{}).Error
}
