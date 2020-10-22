package acceptor_card_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/util"
	"fmt"
)

type AcceptorCard struct {
	//
	Id int64
	// 代理
	Agency string
	// 是否开启={1:关闭，2：开启}
	IfOpen int
	// 承兑人账号
	Acceptor string
	// 卡类型={“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"}
	CardType string
	// 卡号
	CardNo string
	// 卡账户名
	CardAccount string
	// 银行名称
	CardBank string
	// 银行支行
	CardSubBank string
	// 图片地址
	CardImg string
	// 当日剩余可用流水
	DayAvailableAmt int64
	// 当日最大流水
	DayMaxAmt int64
	// 上一次的匹配时间
	LastMatchTime util.JSONTime
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int
}

func (a *AcceptorCard) Add() error {
	session := models.NewSession()
	channelModel := models.NewChannelModel(session)

	channel := channelModel.GetChannel(a.Agency, a.CardType)
	if channel == nil {
		return errors.New(fmt.Sprintf("agency[%s] channel[%s] not found", a.Agency, a.CardType))
	}

	acceptorCard := &models.AcceptorCard{
		Agency:          a.Agency,
		Acceptor:        a.Acceptor,
		CardType:        a.CardType,
		CardNo:          a.CardNo,
		CardAccount:     a.CardAccount,
		CardBank:        a.CardBank,
		CardSubBank:     a.CardSubBank,
		CardImg:         a.CardImg,
		DayAvailableAmt: channel.LimitAmount,
		DayMaxAmt:       channel.LimitAmount,
		IfOpen:          constant.SwitchClose,
		DeleteFlag:      constant.NotDeleted,
	}

	if err := models.NewAcceptorCardModel(session).AddAcceptorCard(acceptorCard); err != nil {
		return err
	}

	return nil
}

func (a *AcceptorCard) UpdateCardInfo() error {
	session := models.NewSession()
	return models.NewAcceptorCardModel(session).EditAcceptorCard(a.Id, map[string]interface{}{
		"card_account":  a.CardAccount,
		"card_bank":     a.CardBank,
		"card_sub_bank": a.CardSubBank,
		"card_img":      a.CardImg,
	})
}

func (a *AcceptorCard) UpdateCardState() error {
	session := models.NewSession()

	// 1.更新卡的承兑状态
	return models.NewAcceptorCardModel(session).EditAcceptorCard(a.Id, map[string]interface{}{
		"if_open": a.IfOpen,
	})
}

// 更新承兑人卡每日可承兑金额
func (a *AcceptorCard) UpdateAcceptorCardDayAvailableAmt() error {
	session := models.NewSession()
	acceptorCardModel := models.NewAcceptorCardModel(session)

	agencyGroup, err := acceptorCardModel.GetAcceptorCardAgencyGroup()
	if err != nil {
		return err
	}

	for _, agency := range agencyGroup {
		go updateAgencyGroup(agency, session, acceptorCardModel)
	}

	return nil
}

//更新代理组
func updateAgencyGroup(agency string, session *models.Session, acceptorCardModel *models.AcceptorCard) {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	channelModel := models.NewChannelModel(session)

	// 获取通道列表
	channelList, err := channelModel.GetChannelGroup(agency)
	if err != nil {
		logf.Error("GetChannelGroup error: ", err.Error())
		return
	}

	for _, item := range channelList {
		channel := channelModel.GetChannel(agency, item)
		if channel == nil {
			continue
		}

		//重置单日可承兑金额
		err := acceptorCardModel.ResetLimitAvailableAmount(agency, channel.Channel, channel.LimitAmount)
		if err != nil {
			logf.Error("ResetLimitAvailableAmount: ", err.Error())
			return
		}
	}
}

func (a *AcceptorCard) Get() (*models.AcceptorCard, error) {
	session := models.NewSession()
	return models.NewAcceptorCardModel(session).GetAcceptorCard(a.Id)
}

func (a *AcceptorCard) GetAll() ([]*models.AcceptorCard, error) {
	session := models.NewSession()
	return models.NewAcceptorCardModel(session).GetAcceptorCards(a.PageNum, a.PageSize, a.getMaps())
}

func (a *AcceptorCard) GetAllCardInfo(agency, acceptor string) (int64, []*models.AcceptorCard, error) {
	session := models.NewSession()
	return models.NewAcceptorCardModel(session).GetAcceptorAllCards(agency, acceptor)
}

func (a *AcceptorCard) CountCardType(agency, acceptor string) (*models.CountCardType, error) {
	session := models.NewSession()
	return models.NewAcceptorCardModel(session).CountCardType(agency, acceptor)
}

func (a *AcceptorCard) Delete() error {
	session := models.NewSession()
	return models.NewAcceptorCardModel(session).DeleteAcceptorCard(a.Id)
}

func (a *AcceptorCard) Count() (int64, error) {
	session := models.NewSession()
	return models.NewAcceptorCardModel(session).GetAcceptorCardTotal(a.getMaps())
}

func (a *AcceptorCard) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.IfOpen != 0 {
		maps["if_open"] = a.IfOpen
	}
	if a.Acceptor != "" {
		maps["acceptor"] = a.Acceptor
	}
	if a.CardType != "" {
		maps["card_type"] = a.CardType
	}
	if a.CardNo != "" {
		maps["card_no"] = a.CardNo
	}
	if a.CardAccount != "" {
		maps["card_account"] = a.CardAccount
	}
	maps["delete_flag"] = constant.NotDeleted
	return maps
}
