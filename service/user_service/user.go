package user_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/util"
)

type User struct {
	// 主键id
	Id int64
	// 用户名
	UserName string
	// 登录密码
	Password string
	// 交易密码
	TradeKey string
	// 所属代理
	Agency string
	// 姓名
	FullName string
	// 电话号码
	PhoneNumber string
	// 地址
	Address string
	// 账户类型(1:系统账户 2:普通账户)
	Type int
	// 角色（1:代理，2：承兑人、3：商家）
	Role int
	// token盐值
	SecretKey string
	// 账户状态（1：暂停，2：启用）
	Status int
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int
}

//添加代理
func (a *User) AddAgency() error {
	session := models.NewSession()
	// 1.开启事务
	session.Begin()
	defer session.Rollback()

	// 2.添加用户信息
	user := &models.User{
		Agency:      a.Agency,
		UserName:    a.UserName,
		FullName:    a.FullName,
		PhoneNumber: a.PhoneNumber,
		Address:     a.Address,
		Type:        constant.OrdinaryAccount,
		Role:        constant.AGENCY,
		Status:      constant.ENABLE,
		Password:    util.EncodeMD5(a.Password),
	}
	if err := models.NewUserModel(session).AddUser(user); err != nil {
		return err
	}

	// 3.添加代理资金信息
	fund := &models.Fund{
		Agency:   a.Agency,
		UserName: a.UserName,
		Type:     constant.OrdinaryAccount,
	}
	if err := models.NewFundModel(session).AddFund(fund); err != nil {
		return err
	}

	// 4.添加通道信息（目前默认三种）
	err := addChannel(a.Agency, session, constant.BANK_CARD, constant.ALIPAY, constant.WECHAT)
	if err != nil {
		return err
	}

	// 5.提交事务
	session.Commit()
	return nil
}

// 代理名称列表
func (a *User) AgencyList() ([]string, error) {
	session := models.NewSession()
	return models.NewUserModel(session).GetAgencyGroup()
}

func addChannel(agency string, session *models.Session, channels ...string) error {
	//单账户每日承兑上限
	limitAmount := int64(49999 * util.GacBase)

	model := models.NewChannelModel(session)
	for _, channel := range channels {
		channel := &models.Channel{
			Agency:      agency,
			Channel:     channel,
			IfOpen:      constant.SwitchOpen,
			BuyMax:      50000 * util.GacBase,
			BuyMin:      9 * util.GacBase,
			LimitAmount: limitAmount,
		}
		if err := model.AddChannel(channel); err != nil {
			return err
		}
	}

	return nil
}

func (a *User) UpdateInfo() error {
	session := models.NewSession()
	return models.NewUserModel(session).EditUser(a.Id, map[string]interface{}{
		"full_name":    a.FullName,
		"phone_number": a.PhoneNumber,
		"address":      a.Address,
	})
}

func (a *User) UpdatePassword() error {
	session := models.NewSession()
	return models.NewUserModel(session).EditUser(a.Id, map[string]interface{}{
		"password": a.Password,
	})
}

func (a *User) UpdateTradeKey() error {
	session := models.NewSession()
	return models.NewUserModel(session).EditUser(a.Id, map[string]interface{}{
		"trade_key": a.TradeKey,
	})
}

func (a *User) Get() (*models.User, error) {
	session := models.NewSession()
	user, err := models.NewUserModel(session).GetUser(a.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *User) GetByAgencyAndUsername(agency, username string) (*models.User, error) {
	session := models.NewSession()
	return models.NewUserModel(session).GetUserByUsernameAndAgency(agency, username)
}

func (a *User) GetAll() ([]*models.User, error) {
	session := models.NewSession()
	users, err := models.NewUserModel(session).GetUsers(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (a *User) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewUserModel(session).ExistUserByID(a.Id)
}

func (a *User) Count() (int64, error) {
	session := models.NewSession()
	return models.NewUserModel(session).GetUserTotal(a.getMaps())
}

func (a *User) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.UserName != "" {
		maps["user_name"] = a.UserName
	}
	if a.Password != "" {
		maps["password"] = a.Password
	}
	if a.TradeKey != "" {
		maps["trade_key"] = a.TradeKey
	}
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.FullName != "" {
		maps["full_name"] = a.FullName
	}
	if a.PhoneNumber != "" {
		maps["phone_number"] = a.PhoneNumber
	}
	if a.Address != "" {
		maps["address"] = a.Address
	}
	if a.Role != 0 {
		maps["role"] = a.Role
	}
	if a.SecretKey != "" {
		maps["secret_key"] = a.SecretKey
	}
	if a.Status != 0 {
		maps["status"] = a.Status
	}

	return maps
}
