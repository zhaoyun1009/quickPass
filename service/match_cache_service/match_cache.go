package match_cache_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
	"fmt"
	"sync"
)

type MatchCache struct {
	//
	Id int64
	// 代理
	Agency string
	// 承兑人账号
	Acceptor string
	// 卡id
	CardId int64
	// 卡类型
	CardType string
	// 最大承兑金额
	MaxMatchedAmount int64
	// 最小承兑金额
	MinMatchedAmount int64
	// 上一次的匹配时间
	LastMatchTime util.JSONTime
	// 资金版本号
	FundVersion int64
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int

	Amount int64
}

type MatchCacheCard struct {
	// 代理
	Agency string
	// 承兑人
	Acceptor string
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
}

type MatchStrategy interface {
	Match(agency string, orderNo string, amount int64, cardType string) (*MatchCacheCard, error)
}

//默认匹配策略
type DefaultMatchStrategy struct{}

// 全局匹配锁
var matchMutex sync.Mutex

func (strategy *DefaultMatchStrategy) Match(agency string, orderNo string, amount int64, cardType string) (*MatchCacheCard, error) {
	matchMutex.Lock()
	defer matchMutex.Unlock()
	// 1、查询出可匹配的卡id
	session := models.NewSession()
	acceptorCardModel := models.NewAcceptorCardModel(session)

	card, err := acceptorCardModel.GetMatchedCard(agency, amount, cardType)

	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, errors.New("匹配卡失败")
	}

	userModel := models.NewUserModel(session)
	user, _ := userModel.GetUserByUsernameAndAgency(agency, card.Acceptor)

	// 开始事务
	session.Begin()
	defer session.Rollback()
	// 3.2、订单待支付
	err = models.NewOrderModel(session).UpdateOrderWaitPay(agency, orderNo, card.Acceptor, user.FullName, card.CardNo, card.CardAccount, card.CardBank, card.CardSubBank, card.CardImg)
	if err != nil {
		return nil, err
	}

	// 4、资金操作，冻结承兑人的币
	_, err = models.NewFundModel(session).Frozen(agency, card.Acceptor, amount)
	if err != nil {
		return nil, err
	}

	// 5、更新卡中的可承兑金额 和匹配时间
	err = acceptorCardModel.SubtractDayAvailableAmount(card.Id, amount)
	if err != nil {
		return nil, err
	}

	session.Commit()

	if cardType != constant.BANK_CARD {
		card.CardImg = fmt.Sprintf("http://%s/%s", setting.MinioSetting.PreUrl, card.CardImg)
	}
	return &MatchCacheCard{
		Agency:      agency,
		Acceptor:    card.Acceptor,
		CardNo:      card.CardNo,
		CardAccount: card.CardAccount,
		CardBank:    card.CardBank,
		CardImg:     card.CardImg,
		CardSubBank: card.CardSubBank,
	}, nil
}

type MatchContext struct {
	MatchStrategy
}

func NewMatchContext(strategy string) *MatchContext {
	context := new(MatchContext)
	switch strategy {
	default:
		context.MatchStrategy = &DefaultMatchStrategy{}
	}

	return context
}

func (a *MatchCache) UpdateMaxMatchAmount() error {
	session := models.NewSession()
	return models.NewMatchCacheModel(session).EditMatchCache(a.Id, map[string]interface{}{
		"max_matched_amount": a.MaxMatchedAmount,
		"last_match_time":    util.JSONTimeNow(),
	})
}

func (a *MatchCache) Delete() error {
	session := models.NewSession()
	return models.NewMatchCacheModel(session).DeleteMatchCache(a.Id)
}

func (a *MatchCache) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewMatchCacheModel(session).ExistMatchCacheByID(a.Id)
}

func (a *MatchCache) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.Acceptor != "" {
		maps["acceptor"] = a.Acceptor
	}
	if a.CardId != 0 {
		maps["card_id"] = a.CardId
	}
	if a.CardType != "" {
		maps["card_type"] = a.CardType
	}
	if a.MaxMatchedAmount != 0 {
		maps["max_matched_amount"] = a.MaxMatchedAmount
	}
	if a.MinMatchedAmount != 0 {
		maps["min_matched_amount"] = a.MinMatchedAmount
	}
	if a.FundVersion != 0 {
		maps["fund_version"] = a.FundVersion
	}

	return maps
}
