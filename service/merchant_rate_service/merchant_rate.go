package merchant_rate_service

import (
	"QuickPass/models"
)

type MerchantRate struct {
	// 主键id
	Id int64
	// 商家账号
	MerchantName string
	// 所属代理
	Agency string
	// 通道名称
	Channel string
	// 通道汇率
	Rate int64
}

func (a *MerchantRate) GetMerchantAllRate(agency, username string) ([]*models.MerchantRate, error) {
	session := models.NewSession()
	return models.NewMerchantRateModel(session).GetMerchantRates(agency, username)
}

func (a *MerchantRate) UpdateMerchantSingleRate(agency, username, channel string, rate int64) error {
	session := models.NewSession()
	return models.NewMerchantRateModel(session).EditMerchantSingleRate(agency, username, channel, rate)
}

func (a *MerchantRate) UpdateMerchantAllChannelRate(agency, username string, rate int64) error {
	session := models.NewSession()
	return models.NewMerchantRateModel(session).EditMerchantAllChannelRate(agency, username, rate)
}
