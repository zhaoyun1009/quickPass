package merchant_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/rsa"
	"QuickPass/pkg/util"
)

type Merchant struct {
	// 主键id
	Id int64
	// 商家账号
	MerchantName string
	// 所属代理
	Agency string
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int

	FullName string
	//密码
	Password    string
	PhoneNumber string
	Address     string
	Rate        int64
}

func (a *Merchant) Add() error {
	session := models.NewSession()
	// 1.开启事务
	session.Begin()
	defer session.Rollback()

	// 2.添加用户信息
	// 2.1、在user表中插入数据 agency+userName构成唯一键；交易密码和登录密码默认123456 MD5加密；role 枚举 ACCEPTOR MERCHANT
	user := &models.User{
		Agency:      a.Agency,
		UserName:    a.MerchantName,
		FullName:    a.FullName,
		PhoneNumber: a.PhoneNumber,
		Address:     a.Address,
		Type:        constant.OrdinaryAccount,
		Role:        constant.MERCHANT,
		Status:      constant.ENABLE,
		Password:    util.EncodeMD5(a.Password),
	}
	if err := models.NewUserModel(session).AddUser(user); err != nil {
		return err
	}

	merchantPrivateKey, merchantPublicKey, err := rsa.GenerateKey()
	if err != nil {
		return err
	}
	systemPrivateKey, systemPublicKey, _ := rsa.GenerateKey()
	// 3.添加商家信息
	merchant := &models.Merchant{
		MerchantName:       a.MerchantName,
		Agency:             a.Agency,
		MerchantPrivateKey: merchantPrivateKey,
		MerchantPublicKey:  merchantPublicKey,
		SystemPrivateKey:   systemPrivateKey,
		SystemPublicKey:    systemPublicKey,
	}
	if err := models.NewMerchantModel(session).AddMerchant(merchant); err != nil {
		return err
	}

	// 4.添加商家资金信息
	fund := &models.Fund{
		Agency:   a.Agency,
		UserName: a.MerchantName,
		Type:     constant.OrdinaryAccount,
	}
	if err := models.NewFundModel(session).AddFund(fund); err != nil {
		return err
	}

	// 5.商家通道汇率
	channelGroup, err := models.NewChannelModel(session).GetChannelGroup(a.Agency)
	if err != nil {
		return nil
	}
	if len(channelGroup) > 0 {
		rateModel := models.NewMerchantRateModel(session)
		for _, value := range channelGroup {
			merchantRate := models.MerchantRate{
				Agency:       a.Agency,
				MerchantName: a.MerchantName,
				Channel:      value,
				Rate:         a.Rate,
			}
			if err := rateModel.AddMerchantRate(&merchantRate); err != nil {
				return err
			}
		}
	}

	// 6.提交事务
	session.Commit()
	return nil
}

func (a *Merchant) Get(agency, username string) (*models.Merchant, error) {
	session := models.NewSession()
	merchant, err := models.NewMerchantModel(session).GetMerchant(agency, username)
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func (a *Merchant) GetAll() ([]*models.Merchant, error) {
	session := models.NewSession()
	merchants, err := models.NewMerchantModel(session).GetMerchants(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}
	return merchants, nil
}

//删除商家
func (a *Merchant) Remove(agency string, username string) error {
	session := models.NewSession()

	// 1.检查资金账户资金
	fund, err := models.NewFundModel(session).GetFund(agency, username)
	if err != nil {
		return err
	}
	if fund == nil {
		return errors.New("fund account not found")
	}

	// 2.1账户还存在资金
	if fund.AvailableAmount+fund.FrozenAmount > 0 {
		return errors.New("fund account has amount")
	}

	// 2.2删除商家信息
	// 1.开启事务
	session.Begin()
	defer session.Rollback()

	// 2.删除商家卡信息
	err = models.NewMerchantCardModel(session).DeleteAllMerchantCardByAgencyAndName(agency, username)
	if err != nil {
		return err
	}

	// 3.删除商家信息
	err = models.NewMerchantModel(session).DeleteMerchant(agency, username)
	if err != nil {
		return err
	}

	// 4.删除资金信息
	err = models.NewFundModel(session).Remove(agency, username)
	if err != nil {
		return err
	}

	// 4.删除用户信息
	err = models.NewUserModel(session).DeleteUser(agency, username)
	if err != nil {
		return err
	}

	// 5.删除商家汇率
	err = models.NewMerchantRateModel(session).DeleteMerchantRate(agency, username)
	if err != nil {
		return err
	}

	// 提交事务
	session.Commit()
	return nil
}

func (a *Merchant) ExistByAgencyAndUsername(agency, username string) (bool, error) {
	session := models.NewSession()
	return models.NewMerchantModel(session).ExistMerchantByAgencyAndUsername(agency, username)
}

func (a *Merchant) UpdateReturnUrl(agency, username, returnUrl string) error {
	session := models.NewSession()
	return models.NewMerchantModel(session).EditMerchant(agency, username, map[string]string{
		"return_url": returnUrl,
	})
}

func (a *Merchant) UpdateNotifyUrl(agency, username, notifyUrl string) error {
	session := models.NewSession()
	return models.NewMerchantModel(session).EditMerchant(agency, username, map[string]string{
		"notify_url": notifyUrl,
	})
}

func (a *Merchant) Count() (int64, error) {
	session := models.NewSession()
	return models.NewMerchantModel(session).GetMerchantTotal(a.getMaps())
}

func (a *Merchant) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.MerchantName != "" {
		maps["merchant_name"] = a.MerchantName
	}
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}

	return maps
}
