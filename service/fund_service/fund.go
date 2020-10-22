package fund_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/util"
)

type Fund struct {
	//
	Id int64
	// 代理
	Agency string
	// 用户名
	UserName string
	// 账户类型(1：系统资金账户  2：普通资金账户)
	Type int
	// 可用资金
	AvailableAmount int64
	// 冻结资金
	FrozenAmount int64
	// 资金版本号
	Version int64
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int
}

func (a *Fund) GetByAgencyAndName(agency string, name string) (*models.Fund, error) {
	session := models.NewSession()
	return models.NewFundModel(session).GetFund(agency, name)
}

func (a *Fund) GetAll() ([]*models.Fund, error) {
	session := models.NewSession()
	return models.NewFundModel(session).GetFunds(a.PageNum, a.PageSize, a.getMaps())
}

func (a *Fund) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewFundModel(session).ExistFundByID(a.Id)
}

func (a *Fund) Count() (int, error) {
	session := models.NewSession()
	return models.NewFundModel(session).GetFundTotal(a.getMaps())
}

func (a *Fund) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.UserName != "" {
		maps["user_name"] = a.UserName
	}

	return maps
}

func (a *Fund) Transfer(fromAgency, toAgency, fromUserName, toUserName, appendInfo string, amount int64) error {
	session := models.NewSession()
	fundModel := models.NewFundModel(session)
	// from账户查询
	fromFund, err := fundModel.GetFund(fromAgency, fromUserName)
	if err != nil {
		return err
	}
	if fromFund == nil {
		return errors.New("转出资金账户不存在")
	}

	// 资金余额判断
	if fromFund.AvailableAmount < amount {
		return errors.New("资金不足")
	}

	// to账户查询
	toFund, err := fundModel.GetFund(toAgency, toUserName)
	if err != nil {
		return err
	}
	if toFund == nil {
		return errors.New("转入资金账户不存在")
	}

	session.Begin()
	defer session.Rollback()

	//1.转账交易
	fromFund, toFund, err = fundModel.Transaction(fromAgency, toAgency, fromUserName, toUserName, amount)
	if err != nil {
		return err
	}
	session.Commit()

	//2.生成账单
	go generatorBill(fromFund, toFund, appendInfo, amount)
	return nil
}

//转账账单生成
func generatorBill(fromFund, toFund *models.Fund, appendInfo string, amount int64) {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	session := models.NewSession()
	billModel := models.NewBillModel(session)
	userModel := models.NewUserModel(session)

	// 1.获取转出用户角色
	fromUser, _ := userModel.GetUserByUsernameAndAgency(fromFund.Agency, fromFund.UserName)
	// 2.获取转入用户角色
	toUser, _ := userModel.GetUserByUsernameAndAgency(toFund.Agency, toFund.UserName)

	session.Begin()
	defer session.Rollback()
	// 2、生成转出方流水账单
	fromBill := &models.Bill{
		Agency:             fromFund.Agency,
		OwnRole:            fromUser.Role,
		OwnUserName:        fromUser.UserName,
		OppositeUserName:   toUser.UserName,
		OppositeRole:       toUser.Role,
		BillNo:             util.GetUniqueNo(),
		Amount:             amount,
		IncomeExpensesType: constant.EXPEND,
		UsableAmount:       fromFund.AvailableAmount,
		FrozenAmount:       fromFund.FrozenAmount,
		BusinessType:       constant.BusinessTypeTransfer,
		AppendInfo:         appendInfo,
	}
	// 4、生成转入方账单
	toBill := &models.Bill{
		Agency:             toFund.Agency,
		OwnRole:            toUser.Role,
		OwnUserName:        toUser.UserName,
		OppositeUserName:   fromUser.UserName,
		OppositeRole:       fromUser.Role,
		BillNo:             util.GetUniqueNo(),
		Amount:             amount,
		IncomeExpensesType: constant.INCOME,
		UsableAmount:       toFund.AvailableAmount,
		FrozenAmount:       toFund.FrozenAmount,
		BusinessType:       constant.BusinessTypeTransfer,
		AppendInfo:         appendInfo,
	}
	err := billModel.AddBill(fromBill)
	if err != nil {
		logf.Error(err.Error())
		return
	}
	err = billModel.AddBill(toBill)
	if err != nil {
		logf.Error(err.Error())
		return
	}

	session.Commit()
}

// 获取所有的承兑人和商家资金账户
func (a *Fund) GetMerchantAndAcceptor() ([]*models.FundExt, error) {
	session := models.NewSession()
	return models.NewFundModel(session).GetMerchantAndAcceptor()
}
