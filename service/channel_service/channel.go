package channel_service

import (
	"QuickPass/models"
	"QuickPass/pkg/util"
)

type Channel struct {
	//
	Id int64
	// 代理
	Agency string
	// 通道名称
	Channel string
	// 通道开关(1: 关闭  2:开启)
	IfOpen int
	// 单账户每日承兑上限
	LimitAmount int64
	// 通道补充信息
	AppendInfo string
	// 买入最小限额
	BuyMin int64
	// 买入最大限额
	BuyMax int64
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int
}

func (a *Channel) Add() error {
	channel := &models.Channel{
		Agency:     a.Agency,
		Channel:    a.Channel,
		IfOpen:     a.IfOpen,
		AppendInfo: a.AppendInfo,
	}

	session := models.NewSession()
	if err := models.NewChannelModel(session).AddChannel(channel); err != nil {
		return err
	}

	return nil
}

func (a *Channel) UpdateLimitAmount() error {
	session := models.NewSession()
	return models.NewChannelModel(session).EditChannel(a.Id, a.Agency, map[string]interface{}{
		"limit_amount": a.LimitAmount,
	})
}

func (a *Channel) UpdateAllLimitAmount() error {
	session := models.NewSession()
	return models.NewChannelModel(session).UpdateByAgency(a.Agency, map[string]interface{}{
		"limit_amount": a.LimitAmount,
	})
}

func (a *Channel) UpdateBuyMaxAmount() error {
	session := models.NewSession()
	return models.NewChannelModel(session).EditChannel(a.Id, a.Agency, map[string]interface{}{
		"buy_max": a.BuyMax,
	})
}

func (a *Channel) UpdateBuyMinAmount() error {
	session := models.NewSession()
	return models.NewChannelModel(session).EditChannel(a.Id, a.Agency, map[string]interface{}{
		"buy_min": a.BuyMin,
	})
}

func (a *Channel) UpdateIfOpen() error {
	session := models.NewSession()
	return models.NewChannelModel(session).EditChannel(a.Id, a.Agency, map[string]interface{}{
		"if_open": a.IfOpen,
		"agency":  a.Agency,
	})
}

func (a *Channel) Get(agency, channel string) *models.Channel {
	session := models.NewSession()
	return models.NewChannelModel(session).GetChannel(agency, channel)
}

func (a *Channel) GetAll() ([]*models.Channel, error) {
	session := models.NewSession()
	channels, err := models.NewChannelModel(session).GetChannels(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (a *Channel) Delete() error {
	session := models.NewSession()
	return models.NewChannelModel(session).DeleteChannel(a.Id)
}

func (a *Channel) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewChannelModel(session).ExistChannelByID(a.Id)
}

func (a *Channel) Count() (int64, error) {
	session := models.NewSession()
	return models.NewChannelModel(session).GetChannelTotal(a.getMaps())
}

func (a *Channel) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.Channel != "" {
		maps["channel"] = a.Channel
	}
	if a.IfOpen != 0 {
		maps["if_open"] = a.IfOpen
	}

	return maps
}
