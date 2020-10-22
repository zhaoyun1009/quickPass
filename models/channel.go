package models

import (
	"QuickPass/pkg/errors"
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type Channel struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 通道名称
	Channel string `json:"channel" gorm:"column:channel"`
	// 通道开关(1: 关闭  2:开启)
	IfOpen int `json:"if_open" gorm:"column:if_open"`
	// 单账户每日承兑上限
	LimitAmount int64 `json:"limit_amount" gorm:"column:limit_amount"`
	// 通道补充信息
	AppendInfo string `json:"append_info" gorm:"column:append_info"`
	// 买入最小限额
	BuyMin int64 `json:"buy_min" gorm:"buy_min"`
	// 买入最大限额
	BuyMax int64 `json:"buy_max" gorm:"buy_max"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置Channel的表名为`channel`
func (Channel) TableName() string {
	return "channel"
}

func NewChannelModel(session *Session) *Channel {
	return &Channel{Session: session}
}

// ExistChannelByID checks if an channel exists based on ID
func (a *Channel) ExistChannelByID(id int64) (bool, error) {
	channel := new(Channel)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(channel).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if channel.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetChannelTotal gets the total number of channels based on the constraints
func (a *Channel) GetChannelTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&Channel{}).Where(maps).Count(&count).Error
	return count, err
}

// GetChannels gets a list of channels based on paging constraints
func (a *Channel) GetChannels(pageNum int, pageSize int, maps interface{}) ([]*Channel, error) {
	var channels []*Channel
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&channels).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return channels, nil
}

// GetChannel Get a single channel based on channel
func (a *Channel) GetChannel(agency, channelName string) *Channel {
	channel := new(Channel)
	err := a.Session.db.Where("agency =? and channel = ?", agency, channelName).First(channel).Error
	if err != nil {
		return nil
	}

	return channel
}

// EditChannel modify a single channel
func (a *Channel) EditChannel(id int64, agency string, data interface{}) error {
	tx := GetSessionTx(a.Session)
	updates := tx.Model(&Channel{}).Where("id = ? and agency = ?", id, agency).Updates(data)

	if !(updates.RowsAffected > 0) {
		return errors.New("data not found")
	}

	return updates.Error
}

func (a *Channel) UpdateByAgency(agency string, data interface{}) error {
	tx := GetSessionTx(a.Session)
	updates := tx.Model(&Channel{}).Where("agency = ?", agency).Update(data)

	if !(updates.RowsAffected > 0) {
		return errors.New("data not found")
	}

	return updates.Error
}

// AddChannel add a single channel
func (a *Channel) AddChannel(channel *Channel) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(channel).Error
}

// DeleteChannel delete a single channel
func (a *Channel) DeleteChannel(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(Channel{}).Error
}

// 获取代理下的通道组
func (a *Channel) GetChannelGroup(agency string) ([]string, error) {
	var channels []*Channel
	err := db.Table(a.TableName()).Select("channel").Where("agency = ?", agency).Group("channel").Find(&channels).Error
	if err != nil {
		return nil, err
	}

	group := make([]string, 0, len(channels))
	for _, item := range channels {
		group = append(group, item.Channel)
	}

	return group, nil
}
