package merchant_card_service

import (
	"QuickPass/models"
	"QuickPass/pkg/util"
)

type MerchantCard struct {
	//
	Id int64
	// 代理
	Agency string
	// 商家账号
	Merchant string
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
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int
}

func (a *MerchantCard) Add() error {
	merchantCard := &models.MerchantCard{
		Agency:      a.Agency,
		Merchant:    a.Merchant,
		CardType:    a.CardType,
		CardNo:      a.CardNo,
		CardAccount: a.CardAccount,
		CardBank:    a.CardBank,
		CardSubBank: a.CardSubBank,
		CardImg:     a.CardImg,
	}

	session := models.NewSession()
	if err := models.NewMerchantCardModel(session).AddMerchantCard(merchantCard); err != nil {
		return err
	}

	a.Id = merchantCard.Id
	return nil
}

func (a *MerchantCard) Edit() error {
	session := models.NewSession()
	return models.NewMerchantCardModel(session).EditMerchantCard(a.Id, map[string]interface{}{
		"agency":        a.Agency,
		"merchant":      a.Merchant,
		"card_type":     a.CardType,
		"card_no":       a.CardNo,
		"card_account":  a.CardAccount,
		"card_bank":     a.CardBank,
		"card_sub_bank": a.CardSubBank,
		"card_img":      a.CardImg,
	})
}

func (a *MerchantCard) Get() (*models.MerchantCard, error) {
	session := models.NewSession()
	return models.NewMerchantCardModel(session).GetMerchantCard(a.Id)
}

func (a *MerchantCard) GetAll() ([]*models.MerchantCard, error) {
	session := models.NewSession()
	return models.NewMerchantCardModel(session).GetMerchantCards(a.PageNum, a.PageSize, a.getMaps())
}

func (a *MerchantCard) Delete(id int64, agency, name string) error {
	session := models.NewSession()
	return models.NewMerchantCardModel(session).DeleteMerchantCardByAgencyAndName(id, agency, name)
}

func (a *MerchantCard) Exist(agency, merchant, cardNo, cardType string) (bool, error) {
	session := models.NewSession()
	return models.NewMerchantCardModel(session).ExistMerchantCard(agency, merchant, cardNo, cardType)
}

func (a *MerchantCard) Count() (int64, error) {
	session := models.NewSession()
	return models.NewMerchantCardModel(session).GetMerchantCardTotal(a.getMaps())
}

func (a *MerchantCard) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.Merchant != "" {
		maps["merchant"] = a.Merchant
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
	if a.CardBank != "" {
		maps["card_bank"] = a.CardBank
	}
	if a.CardSubBank != "" {
		maps["card_sub_bank"] = a.CardSubBank
	}
	if a.CardImg != "" {
		maps["card_img"] = a.CardImg
	}

	return maps
}
