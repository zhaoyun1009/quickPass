package acceptor_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/util"
)

type Acceptor struct {
	//
	Id int64
	// 承兑人账号
	Acceptor string
	// 所属代理
	Agency string
	// 是否自动承兑（1：不自动，2：自动）
	IfAutoAccept int
	// 承兑开关(1：关闭，2：开启)
	AcceptSwitch int
	// 承兑账户状态(1：关闭，2：开启)
	AcceptStatus int
	// 保证金
	Deposit int64
	// 最大承兑金额
	MaxAcceptedAmount int64
	// 最小承兑金额
	MinAcceptedAmount int64
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int

	//承兑人姓名
	FullName string
	//电话号码
	PhoneNumber string
	//密码
	Password string
	//地址
	Address string
}

func (a *Acceptor) Add() error {
	session := models.NewSession()
	// 1.开启事务
	session.Begin()
	defer session.Rollback()

	// 2.添加用户信息
	// 2.1、在user表中插入数据 agency+userName构成唯一键；交易密码和登录密码默认123456 MD5加密；role 枚举 ACCEPTOR MERCHANT
	user := &models.User{
		Agency:      a.Agency,
		UserName:    a.Acceptor,
		FullName:    a.FullName,
		PhoneNumber: a.PhoneNumber,
		Address:     a.Address,
		Type:        constant.OrdinaryAccount,
		Role:        constant.ACCEPTOR,
		Status:      constant.ENABLE,
		Password:    util.EncodeMD5(a.Password),
	}
	if err := models.NewUserModel(session).AddUser(user); err != nil {
		return err
	}

	acceptor := &models.Acceptor{
		Acceptor:          a.Acceptor,
		Agency:            a.Agency,
		IfAutoAccept:      constant.Manual,
		AcceptSwitch:      constant.SwitchClose,
		AcceptStatus:      constant.SwitchOpen,
		Deposit:           a.Deposit,
		MaxAcceptedAmount: a.MaxAcceptedAmount,
		MinAcceptedAmount: a.MinAcceptedAmount,
	}
	// 3.添加承兑人信息
	if err := models.NewAcceptorModel(session).AddAcceptor(acceptor); err != nil {
		return err
	}

	// 4.添加承兑人资金信息
	fund := &models.Fund{
		Agency:   a.Agency,
		UserName: a.Acceptor,
		Type:     constant.OrdinaryAccount,
	}
	if err := models.NewFundModel(session).AddFund(fund); err != nil {
		return err
	}
	// 提交事务
	session.Commit()
	return nil
}

func (a *Acceptor) UpdateAcceptSwitch(agency, acceptorName string, acceptSwitch int) error {
	session := models.NewSession()
	acceptorModel := models.NewAcceptorModel(session)
	acceptor, err := acceptorModel.GetAcceptorByAgencyAndUsername(agency, acceptorName)
	if err != nil {
		return err
	}
	if acceptor == nil {
		return errors.New(e.GetMsg(e.ERROR_NOT_EXIST_ACCEPTOR))
	}
	if acceptor.AcceptStatus == constant.SwitchClose {
		return errors.New(e.GetMsg(e.AcceptorStatusClosed))
	}

	return models.NewAcceptorModel(session).EditAcceptor(agency, acceptorName, map[string]interface{}{
		"accept_switch": acceptSwitch,
	})
}

func (a *Acceptor) UpdateAcceptStatus(agency, acceptor string, acceptStatus int) error {
	m := make(map[string]interface{}, 0)
	if acceptStatus == constant.SwitchClose {
		m["accept_switch"] = constant.SwitchClose
	}
	m["accept_status"] = acceptStatus
	session := models.NewSession()
	return models.NewAcceptorModel(session).EditAcceptor(agency, acceptor, m)
}

func (a *Acceptor) UpdateIfAutoAccept(agency, acceptor string, ifAutoAccept int) error {
	session := models.NewSession()
	return models.NewAcceptorModel(session).EditAcceptor(agency, acceptor, map[string]interface{}{
		"if_auto_accept": ifAutoAccept,
	})
}

func (a *Acceptor) Get() (*models.Acceptor, error) {
	session := models.NewSession()
	acceptor, err := models.NewAcceptorModel(session).GetAcceptor(a.Id)
	if err != nil {
		return nil, err
	}

	return acceptor, nil
}

func (a *Acceptor) GetByAgencyAndUsername(agency string, username string) (*models.Acceptor, error) {
	session := models.NewSession()
	return models.NewAcceptorModel(session).GetAcceptorByAgencyAndUsername(agency, username)
}

func (a *Acceptor) GetAll() ([]*models.Acceptor, error) {
	session := models.NewSession()
	acceptors, err := models.NewAcceptorModel(session).GetAcceptors(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	return acceptors, nil
}

func (a *Acceptor) Delete() error {
	session := models.NewSession()
	return models.NewAcceptorModel(session).DeleteAcceptor(a.Id)
}

func (a *Acceptor) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewAcceptorModel(session).ExistAcceptorByID(a.Id)
}

func (a *Acceptor) ExistByAgencyAndUsername() (bool, error) {
	session := models.NewSession()
	return models.NewAcceptorModel(session).ExistAcceptorByAgencyAndUsername(a.Agency, a.Acceptor)
}

func (a *Acceptor) Count() (int64, error) {
	session := models.NewSession()
	return models.NewAcceptorModel(session).GetAcceptorTotal(a.getMaps())
}

func (a *Acceptor) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Acceptor != "" {
		maps["acceptor"] = a.Acceptor
	}
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.IfAutoAccept != 0 {
		maps["if_auto_accept"] = a.IfAutoAccept
	}
	if a.AcceptSwitch != 0 {
		maps["accept_switch"] = a.AcceptSwitch
	}
	if a.AcceptStatus != 0 {
		maps["accept_status"] = a.AcceptStatus
	}

	return maps
}
