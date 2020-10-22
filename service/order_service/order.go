package order_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/http_client"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/mq"
	"QuickPass/pkg/rsa"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Order struct {
	//
	Id int64
	// 代理
	Agency string
	// 订单号
	OrderNo string
	// 订单类型（1：转账，2：买入，3：卖出）
	OrderType int64
	// 订单状态（1：创建，2：待支付，3：待放行，4：已取消  5：已失败  6：已完成）
	OrderStatus int64
	// 异常单号
	AbnormalOrderNo string
	// 发起者账号
	FromUserName string
	// 商家账号
	MerchantUserName string
	// 接受者账号
	ToUserName string
	// 订单所匹配的承兑人昵称
	AcceptorNickname string
	// 金额
	Amount int64
	// 提交类型=（1：直接提交，2：接口提交）
	SubmitType int
	// 通道类型= (BANK_CARD,ALIPAY,WECHAT)
	ChannelType string
	// 摘要
	AppendInfo string
	// ip地址
	ClientIp string

	PageNum  int
	PageSize int

	//----------------------------------------------
	// 最大金额
	MaxAmount int64
	// 最小金额
	MinAmount int64
	// 开始时间
	StartTime string
	// 结束时间
	EndTime string
	// 商家订单号
	MerchantOrderNo string
	// 完成开始时间
	FinishStartTime string
	// 完成结束时间
	FinishEndTime string
}

// 创建买单
// agency 代理
// fromUsername 发起者名称
// merchantUsername 商家名称
// payType 支付通道类型= (BANK_CARD,ALIPAY,WECHAT)
// ip 请求IP
// amount 订单金额
// submitType 提交类型[接口，直接]
// merchantOrderNo 商家对接时候的商家订单号，没有为空串
func (a *Order) CreateBuy(agency, fromUsername, merchantUsername, payType, ip string,
	amount int64, submitType int, merchantOrderNo, returnUrl, notifyUrl, appendInfo string) (*models.Order, error) {
	session := models.NewSession()

	session.Begin()
	defer session.Rollback()

	order, err := createOrder(session, agency, fromUsername, merchantUsername, payType, ip, constant.ORDER_TYPE_BUY, amount, submitType)
	if err != nil {
		return nil, err
	}

	if merchantOrderNo == "" {
		merchantOrderNo = order.OrderNo
	}
	merchantOrder := &models.MerchantOrder{
		OrderNo:         order.OrderNo,
		Agency:          agency,
		Username:        merchantUsername,
		MerchantOrderNo: merchantOrderNo,
		SubmitType:      submitType,
		ReturnUrl:       returnUrl,
		CallbackUrl:     notifyUrl,
		CallbackStatus:  constant.CallbackUnprocessed,
		AppendInfo:      appendInfo,
	}
	if err := models.NewMerchantOrderModel(session).AddMerchantOrder(merchantOrder); err != nil {
		return nil, err
	}

	session.Commit()
	return order, nil
}

// 创建卖单
func (a *Order) CreateSell(agency, fromUsername, merchantUsername, payType, ip string, amount int64) (*models.Order, error) {
	session := models.NewSession()
	return createOrder(session, agency, fromUsername, merchantUsername, payType, ip, constant.ORDER_TYPE_SELL, amount, constant.SubmitTypeDirect)
}

// 创建订单
func createOrder(session *models.Session, agency, fromUsername, merchantUsername, payType, ip string,
	orderType, amount int64, submitType int) (*models.Order, error) {
	order := &models.Order{
		Agency:           agency,
		OrderNo:          util.GetUniqueNo(),
		OrderType:        orderType,
		OrderStatus:      constant.ORDER_STATUS_CREATE,
		FromUserName:     fromUsername,
		MerchantUserName: merchantUsername,
		Amount:           amount,
		ChannelType:      payType,
		ClientIp:         ip,
		SubmitType:       submitType,
	}

	if err := models.NewOrderModel(session).AddOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

// 更新订单失败
func (a *Order) UpdateOrderFailed(agency, orderNo string) error {
	session := models.NewSession()
	return models.NewOrderModel(session).UpdateOrderFailed(agency, orderNo)
}

//买入订单放行
func (a *Order) DischargeBuyOrder(agency, orderNo string) error {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	fundModel := models.NewFundModel(session)
	rateModel := models.NewMerchantRateModel(session)
	cardModel := models.NewAcceptorCardModel(session)

	//1.查询订单信息
	order, err := orderModel.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	//2.查询通道汇率
	merchantRate := rateModel.GetMerchantRate(agency, order.MerchantUserName, order.ChannelType)
	if merchantRate == nil {
		return errors.New(fmt.Sprintf("[agency: %s, merchant: %s, channel: %s not found]", agency, order.MerchantUserName, order.ChannelType))
	}
	// 3.1 计算手续费
	fee := merchantRate.Rate * order.Amount / util.GacBase

	session.Begin()
	defer session.Rollback()

	// 2.更新订单完成
	err = orderModel.UpdateOrderFinished(agency, orderNo)
	if err != nil {
		return err
	}

	// 3.判断异常单里面有无关联的，有则更新异常单状态
	err = updateAbnormalOrderFinish(agency, order.AbnormalOrderNo, session)
	if err != nil {
		return err
	}

	// 3.买单资金放行(承兑人冻结账户--》(商家账户, 代理账户))
	agencyFund, fromFund, toFund, err := fundModel.DischargeBuy(agency, order.ToUserName, order.MerchantUserName, order.Amount, fee)
	if err != nil {
		return err
	}

	// 4.承兑人卡冻结金额减少
	err = cardModel.SubtractDayFrozenAmount(agency, order.ToUserName, order.ChannelType, order.CardNo, order.Amount)
	if err != nil {
		return err
	}

	session.Commit()

	// 发送订单状态改变通知
	//1.查询订单信息
	order, _ = orderModel.GetOrderByOrderNo(agency, orderNo)
	go mq.SendOrderStatusChange(order)
	//异常单通知
	if order.AbnormalOrderNo != "" {
		abnormalOrder, _ := models.NewAbnormalOrderModel(session).GetAbnormalOrder(agency, order.AbnormalOrderNo)
		go mq.SendAbnormalOrderStatusChange(abnormalOrder)
	}

	merchantOrder, _ := models.NewMerchantOrderModel(session).GetMerchantOrder(agency, orderNo)
	// 5.生成账单
	go generatorBuyBill(agencyFund, fromFund, toFund, merchantRate.Rate, order.Amount, fee, orderNo, merchantOrder.MerchantOrderNo)
	return nil
}

// 超时取消订单放行
func (a *Order) DischargeCancelOrder(agency, orderNo string) error {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	fundModel := models.NewFundModel(session)
	acceptorCardModel := models.NewAcceptorCardModel(session)
	merchantOrderModel := models.NewMerchantOrderModel(session)

	// 开启事务
	session.Begin()
	defer session.Rollback()

	// 1、重新唤起，更新订单状态为已支付
	if err := orderModel.WakeUpCancelOrder(agency, orderNo); err != nil {
		return err
	}

	// 2、查询订单金额，卡信息
	order, _ := orderModel.GetOrderByOrderNo(agency, orderNo)
	amount := order.Amount

	// 3、资金操作，冻结承兑人的币
	_, err := fundModel.Frozen(agency, order.ToUserName, amount)
	if err != nil {
		return err
	}

	// 4、更新卡中的可承兑金额 和匹配时间
	card, err := acceptorCardModel.GetAcceptorCardByCardNo(agency, order.ToUserName, order.CardNo, order.ChannelType)
	if err != nil {
		return err
	}
	if err = acceptorCardModel.SubtractDayAvailableAmount(card.Id, amount); err != nil {
		return err
	}

	// 5.更新回调状态
	merchantOrder, _ := merchantOrderModel.GetMerchantOrder(agency, orderNo)

	editParams := map[string]interface{}{
		"callback_info":   gorm.Expr("?", ""),
		"callback_times":  gorm.Expr("?", 0),
		"callback_lock":   gorm.Expr("?", constant.CallbackUnLock),
		"callback_status": gorm.Expr("?", constant.CallbackUnprocessed),
	}
	if err = merchantOrderModel.EditMerchantOrder(merchantOrder.Id, editParams); err != nil {
		return err
	}

	// 5.提交事务
	session.Commit()

	// 6.放行
	return a.DischargeBuyOrder(agency, orderNo)
}

// 回调商家
func (a *Order) CallbackMerchant(order *models.Order) {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	session := models.NewSession()
	merchantOrderModel := models.NewMerchantOrderModel(session)

	merchantOrder, err := merchantOrderModel.GetMerchantOrder(order.Agency, order.OrderNo)
	if err != nil {
		logf.Error("GetMerchantOrder", err.Error())
		return
	}
	if merchantOrder == nil || merchantOrder.CallbackUrl == "" {
		return
	}

	// 先锁数据, 请求结束后释放
	lockParams := map[string]interface{}{
		"callback_lock": constant.CallbackLock,
	}
	_ = merchantOrderModel.EditMerchantOrder(merchantOrder.Id, lockParams)

	editParams := make(map[string]interface{}, 0)
	var callbackFlag int

	err = signMerchantParams(session, order, merchantOrder.MerchantOrderNo, merchantOrder.CallbackUrl)
	if err != nil {
		callbackFlag = constant.CallbackFailed
		editParams["callback_info"] = gorm.Expr("?", err.Error())
	} else {
		callbackFlag = constant.CallbackSuccess
	}
	editParams["callback_lock"] = gorm.Expr("?", constant.CallbackUnLock)
	editParams["callback_times"] = gorm.Expr("callback_times + 1")
	editParams["callback_status"] = gorm.Expr("?", callbackFlag)
	_ = merchantOrderModel.EditMerchantOrder(merchantOrder.Id, editParams)
}

// 回调商家
func (a *Order) AgainCallbackMerchant(order *models.Order, callbackUrl string) {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	session := models.NewSession()
	merchantOrderModel := models.NewMerchantOrderModel(session)

	merchantOrder, err := merchantOrderModel.GetMerchantOrder(order.Agency, order.OrderNo)
	if err != nil {
		logf.Error("GetMerchantOrder", err.Error())
		return
	}
	if merchantOrder == nil || merchantOrder.CallbackUrl == "" {
		return
	}

	// 先锁数据, 请求结束后释放
	lockParams := map[string]interface{}{
		"callback_lock": constant.CallbackLock,
	}
	_ = merchantOrderModel.EditMerchantOrder(merchantOrder.Id, lockParams)

	editParams := make(map[string]interface{}, 0)
	var callbackFlag int

	err = signMerchantParams(session, order, merchantOrder.MerchantOrderNo, callbackUrl)
	if err != nil {
		callbackFlag = constant.CallbackFailed
		editParams["callback_info"] = gorm.Expr("?", err.Error())
	} else {
		callbackFlag = constant.CallbackSuccess
	}
	editParams["callback_lock"] = gorm.Expr("?", constant.CallbackUnLock)
	editParams["callback_times"] = gorm.Expr("callback_times + 1")
	editParams["callback_status"] = gorm.Expr("?", callbackFlag)
	_ = merchantOrderModel.EditMerchantOrder(merchantOrder.Id, editParams)
}

func signMerchantParams(session *models.Session, order *models.Order, merchantOrderNo, callbackUrl string) error {
	merchantModel := models.NewMerchantModel(session)
	merchant, err := merchantModel.GetMerchant(order.Agency, order.MerchantUserName)
	if err != nil {
		return err
	}

	params := request.CallbackParams{
		MerchantOrderNo: merchantOrderNo,
		Amount:          order.Amount,
		Member:          order.FromUserName,
		OrderStatus:     order.OrderStatus,
		ChannelType:     order.ChannelType,
		CreateTime:      order.CreateTime,
		FinishTime:      order.FinishTime,
	}

	bt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	sign, err := rsa.RSASign(util.TransHtmlJson(bt), merchant.SystemPrivateKey)
	if err != nil {
		return err
	}

	_, err = http_client.Post(callbackUrl, request.SignCallbackParams{
		Data: params,
		Sign: sign,
	}, "application/json")
	if err != nil {
		return err
	}

	return nil
}

//放行买入生成账单
func generatorBuyBill(agencyFund, acceptorFund, merchantFund *models.Fund, rate, amount, fee int64, orderNo, merchantOrderNo string) {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	session := models.NewSession()
	billModel := models.NewBillModel(session)

	currentRate := fmt.Sprintf("%d", rate)

	session.Begin()
	defer session.Rollback()
	// 4.生成账单
	//承兑人账单
	acceptorBill := &models.Bill{
		Agency:             acceptorFund.Agency,
		OwnRole:            constant.ACCEPTOR,
		OwnUserName:        acceptorFund.UserName,
		OppositeUserName:   merchantFund.UserName,
		OppositeRole:       constant.MERCHANT,
		BillNo:             util.GetUniqueNo(),
		Amount:             amount,
		OrderNo:            orderNo,
		MerchantOrderNo:    merchantOrderNo,
		IncomeExpensesType: constant.EXPEND,
		UsableAmount:       acceptorFund.AvailableAmount,
		FrozenAmount:       acceptorFund.FrozenAmount,
		BusinessType:       constant.BusinessTypeBuy,
		CurrentRate:        currentRate,
	}
	//商家账单
	merchantBill := &models.Bill{
		Agency:             merchantFund.Agency,
		OwnRole:            constant.MERCHANT,
		OwnUserName:        merchantFund.UserName,
		OppositeUserName:   acceptorFund.UserName,
		OppositeRole:       constant.ACCEPTOR,
		BillNo:             util.GetUniqueNo(),
		Amount:             amount - fee,
		OrderNo:            orderNo,
		MerchantOrderNo:    merchantOrderNo,
		IncomeExpensesType: constant.INCOME,
		UsableAmount:       merchantFund.AvailableAmount,
		FrozenAmount:       merchantFund.FrozenAmount,
		BusinessType:       constant.BusinessTypeBuy,
		CurrentRate:        currentRate,
	}
	err := billModel.AddBill(acceptorBill)
	if err != nil {
		logf.Error(err.Error())
		return
	}
	err = billModel.AddBill(merchantBill)
	if err != nil {
		logf.Error(err.Error())
		return
	}
	if fee > 0 {
		//代理手续费账单
		agencyBill := &models.Bill{
			Agency:             agencyFund.Agency,
			OwnRole:            constant.AGENCY,
			OwnUserName:        agencyFund.UserName,
			OppositeUserName:   acceptorFund.UserName,
			OppositeRole:       constant.ACCEPTOR,
			BillNo:             util.GetUniqueNo(),
			Amount:             fee,
			OrderNo:            orderNo,
			MerchantOrderNo:    merchantOrderNo,
			IncomeExpensesType: constant.INCOME,
			UsableAmount:       agencyFund.AvailableAmount,
			FrozenAmount:       agencyFund.FrozenAmount,
			BusinessType:       constant.BusinessTypeBuyFee,
			CurrentRate:        currentRate,
		}
		err = billModel.AddBill(agencyBill)
		if err != nil {
			logf.Error(err.Error())
			return
		}
	}

	session.Commit()
}

//卖出订单放行
func (a *Order) DischargeSellOrder(agency, orderNo string) error {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	fundModel := models.NewFundModel(session)
	merchantRateModel := models.NewMerchantRateModel(session)

	//1.查询订单信息
	order, err := orderModel.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New(fmt.Sprintf("order[%s] not found", orderNo))
	}

	//2.查询通道汇率
	merchantRate := merchantRateModel.GetMerchantRate(agency, order.MerchantUserName, order.ChannelType)
	if merchantRate == nil {
		return errors.New(fmt.Sprintf("merchantRate[%s] not found", order.ChannelType))
	}

	session.Begin()
	defer session.Rollback()

	// 2.更新订单完成
	err = orderModel.UpdateOrderFinished(agency, orderNo)
	if err != nil {
		return err
	}

	// 3.资金放行(商家冻结账户--》系统账户)
	merchantFund, adminFund, err := fundModel.DischargeSell(agency, order.FromUserName, order.Amount)
	if err != nil {
		return err
	}

	session.Commit()

	order, _ = orderModel.GetOrderByOrderNo(agency, orderNo)
	// 发送订单状态改变通知
	go mq.SendOrderStatusChange(order)
	// 5.生成账单
	go generatorSellBill(merchantFund, adminFund, merchantRate.Rate, order.FromUserName, order.Amount, orderNo)
	return nil
}

func generatorSellBill(merchantFund, adminFund *models.Fund, rate int64, formUsername string, amount int64, orderNo string) {
	session := models.NewSession()
	billModel := models.NewBillModel(session)

	currentRate := fmt.Sprintf("%d", rate)

	session.Begin()
	defer session.Rollback()
	// 4.生成账单
	//商家账单
	acceptorBill := &models.Bill{
		Agency:             merchantFund.Agency,
		OwnRole:            constant.MERCHANT,
		OwnUserName:        formUsername,
		OppositeUserName:   adminFund.UserName,
		OppositeRole:       constant.SYSTEM_USEER,
		BillNo:             util.GetUniqueNo(),
		Amount:             amount,
		OrderNo:            orderNo,
		IncomeExpensesType: constant.EXPEND,
		UsableAmount:       merchantFund.AvailableAmount,
		FrozenAmount:       merchantFund.FrozenAmount,
		BusinessType:       constant.BusinessTypeSell,
		CurrentRate:        currentRate,
	}
	//系统账单
	systemBill := &models.Bill{
		Agency:             adminFund.Agency,
		OwnRole:            constant.SYSTEM_USEER,
		OwnUserName:        adminFund.UserName,
		OppositeUserName:   formUsername,
		OppositeRole:       constant.MERCHANT,
		BillNo:             util.GetUniqueNo(),
		Amount:             amount,
		OrderNo:            orderNo,
		IncomeExpensesType: constant.INCOME,
		UsableAmount:       adminFund.AvailableAmount,
		FrozenAmount:       adminFund.FrozenAmount,
		BusinessType:       constant.BusinessTypeSell,
		CurrentRate:        currentRate,
	}

	err := billModel.AddBill(acceptorBill)
	if err != nil {
		logf.Error(err.Error())
		return
	}
	err = billModel.AddBill(systemBill)
	if err != nil {
		logf.Error(err.Error())
		return
	}
	session.Commit()
}

// 卖单匹配,冻结商家账户金额,修改订单状态
func (a *Order) MatchSell(agency, orderNo, toUsername, cardNo, cardAccount, cardBank, cardSubBank string) error {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	fundModel := models.NewFundModel(session)

	// 获取订单信息
	order, err := orderModel.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New(fmt.Sprintf("MatchSell order[%s] not found", orderNo))
	}

	// 开启事务
	session.Begin()
	defer session.Rollback()

	// 更新订单待支付
	// 卖单都是银行卡支付 不存在收款码
	err = orderModel.UpdateOrderWaitPay(agency, orderNo, toUsername, "", cardNo, cardAccount, cardBank, cardSubBank, "")
	if err != nil {
		return err
	}

	// 冻结商家账户资金
	_, err = fundModel.Frozen(order.Agency, order.FromUserName, order.Amount)
	if err != nil {
		return err
	}

	session.Commit()

	// 4、插入消息到队列中，通知代理
	order, _ = orderModel.GetOrderByOrderNo(order.Agency, order.OrderNo)
	// 发送订单状态改变通知
	go mq.SendOrderStatusChange(order)
	return nil
}

//买入订单取消
func (a *Order) CancelBuyOrder(agency, orderNo string) error {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	fundModel := models.NewFundModel(session)
	acceptorCardModel := models.NewAcceptorCardModel(session)

	//1.查询订单信息
	order, err := orderModel.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New(fmt.Sprintf("CancelBuyOrder order[%s] not found", orderNo))
	}

	session.Begin()
	defer session.Rollback()

	// 2.更新订单被取消
	err = orderModel.UpdateOrderCancel(agency, orderNo)
	if err != nil {
		return err
	}

	// 更新异常单完成
	err = updateAbnormalOrderFinish(agency, order.AbnormalOrderNo, session)
	if err != nil {
		return err
	}

	// 3.解冻承兑人资金
	_, err = fundModel.UnFrozen(order.Agency, order.ToUserName, order.Amount)
	if err != nil {
		return err
	}

	// 4.订单被取消，卡的可承兑金额增量为订单金额
	err = acceptorCardModel.AddDayAvailableAmount(order.Agency, order.ToUserName, order.ChannelType, order.CardNo, order.Amount)
	if err != nil {
		return err
	}

	session.Commit()

	// 4、消息队列中
	// 发送订单状态改变通知
	order, _ = orderModel.GetOrderByOrderNo(agency, orderNo)
	go mq.SendOrderStatusChange(order)
	// 异常单通知
	if order.AbnormalOrderNo != "" {
		abnormalOrder, _ := models.NewAbnormalOrderModel(session).GetAbnormalOrder(agency, order.AbnormalOrderNo)
		go mq.SendAbnormalOrderStatusChange(abnormalOrder)
	}
	return nil
}

func updateAbnormalOrderFinish(agency, abnormalOrderNo string, session *models.Session) error {
	if abnormalOrderNo != "" {
		return models.NewAbnormalOrderModel(session).EditAbnormalOrder(agency, abnormalOrderNo, map[string]interface{}{
			"abnormal_order_status": constant.Processed,
			"finish_time":           util.JSONTimeNow(),
		})
	}
	return nil
}

//卖出订单取消
func (a *Order) CancelSellOrder(agency, orderNo string) error {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	fundModel := models.NewFundModel(session)

	//1.查询订单信息
	order, err := orderModel.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New(fmt.Sprintf("CancelSellOrder order[%s] not found", orderNo))
	}

	session.Begin()
	defer session.Rollback()

	// 2.更新订单被取消
	err = orderModel.UpdateOrderCancel(agency, orderNo)
	if err != nil {
		return err
	}

	// 3.解冻商家资金
	_, err = fundModel.UnFrozen(order.Agency, order.FromUserName, order.Amount)
	if err != nil {
		return err
	}
	session.Commit()

	return nil
}

//买入订单确认
func (a *Order) ConfirmBuyOrder(agency, orderNo string) error {
	session := models.NewSession()
	return models.NewOrderModel(session).UpdateOrderWaitDischarge(agency, orderNo)
}

//卖出订单确认
func (a *Order) ConfirmSellOrder(agency, orderNo string) error {
	session := models.NewSession()
	return models.NewOrderModel(session).UpdateSellOrderWaitDischarge(agency, orderNo)
}

func (a *Order) GetByOrderNo(agency, orderNo string) (*models.Order, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).GetOrderByOrderNo(agency, orderNo)
}

// 获取超时的订单列表（不区分代理）
func (a *Order) GetOrderTimeoutList(minute int64) ([]*models.Order, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).GetOrderTimeoutList(minute)
}

// 获取回调的订单列表（不区分代理）
func (a *Order) GetOrderCallbackList() ([]*models.Order, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).GetOrderCallbackList()
}

// 获取自动放行的卖出订单列表（不区分代理）
func (a *Order) GetSellOrderAutoDischargeList(minute int64) ([]*models.Order, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).GetSellOrderAutoDischargeList(minute)
}

func (a *Order) GetAcceptorGroup(agency string) ([]string, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).GetAcceptorGroup(agency)
}

func (a *Order) GetAll() ([]*models.Order, error) {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	merchantOrderModel := models.NewMerchantOrderModel(session)

	var orders []*models.Order
	var err error
	if a.MerchantOrderNo != "" {
		var no []*models.MerchantOrder
		if a.MerchantUserName == "" {
			no, err = merchantOrderModel.GetMerchantOrderByMerchantOrderNo(a.Agency, a.MerchantOrderNo)
		} else {
			var merchantOrderNo *models.MerchantOrder
			merchantOrderNo, err = merchantOrderModel.GetOneMerchantOrderByMerchantOrderNo(a.Agency, a.MerchantUserName, a.MerchantOrderNo)
			no = append(no, merchantOrderNo)
		}
		if err != nil || no == nil {
			return orders, err
		}
		for _, mNo := range no {
			order, _ := orderModel.GetOrderByOrderNo(a.Agency, mNo.OrderNo)
			orders = append(orders, order)
		}
	} else {
		orders, err = orderModel.GetOrders(&models.OrderListQueryParams{
			MaxAmount:       a.MaxAmount,
			MinAmount:       a.MinAmount,
			StartTime:       a.StartTime,
			EndTime:         a.EndTime,
			FinishStartTime: a.FinishStartTime,
			FinishEndTime:   a.FinishEndTime,
		}, a.PageNum, a.PageSize, a.getMaps())
	}
	if err == nil {
		for _, item := range orders {
			merchantOrder, err := merchantOrderModel.GetMerchantOrder(item.Agency, item.OrderNo)
			if err == nil && merchantOrder != nil {
				item.MerchantOrderNo = merchantOrder.MerchantOrderNo
				item.CallbackStatus = merchantOrder.CallbackStatus
				item.CallbackInfo = merchantOrder.CallbackInfo
				item.CallbackUrl = merchantOrder.CallbackUrl
			}
		}
	}
	return orders, err
}

//查询当前用户正在进行的订单
func (a *Order) GetCurrentMemberOrder(agency, merchant, member string) (*models.Order, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).GetCurrentOrder(agency, merchant, member)
}

//  查询当前用户的IP下的订单
func (a *Order) GetCurrentIPOrder(agency, ip string) ([]models.Order, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).GetCurrentIPOrder(agency, ip)
}

func (a *Order) Delete() error {
	session := models.NewSession()
	return models.NewOrderModel(session).DeleteOrder(a.Id)
}

func (a *Order) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewOrderModel(session).ExistOrderByID(a.Id)
}

func (a *Order) Count() (int64, error) {
	session := models.NewSession()
	if a.MerchantOrderNo != "" {
		var no []*models.MerchantOrder
		var err error
		merchantOrderModel := models.NewMerchantOrderModel(session)
		if a.MerchantUserName == "" {
			no, err = merchantOrderModel.GetMerchantOrderByMerchantOrderNo(a.Agency, a.MerchantOrderNo)
		} else {
			var merchantOrder *models.MerchantOrder
			merchantOrder, err = merchantOrderModel.GetOneMerchantOrderByMerchantOrderNo(a.Agency, a.MerchantUserName, a.MerchantOrderNo)
			no = append(no, merchantOrder)
		}
		if err != nil {
			return 0, err
		}
		return int64(len(no)), nil
	}
	return models.NewOrderModel(session).GetOrderTotal(&models.OrderListQueryParams{
		MaxAmount:       a.MaxAmount,
		MinAmount:       a.MinAmount,
		StartTime:       a.StartTime,
		EndTime:         a.EndTime,
		FinishStartTime: a.FinishStartTime,
		FinishEndTime:   a.FinishEndTime,
	}, a.getMaps())
}

func (a *Order) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.OrderNo != "" {
		maps["order_no"] = a.OrderNo
	}
	if a.OrderType != 0 {
		maps["order_type"] = a.OrderType
	}
	if a.OrderStatus != 0 {
		maps["order_status"] = a.OrderStatus
	}
	if a.FromUserName != "" {
		maps["from_user_name"] = a.FromUserName
	}
	if a.MerchantUserName != "" {
		maps["merchant_user_name"] = a.MerchantUserName
	}
	if a.ToUserName != "" {
		maps["to_user_name"] = a.ToUserName
	}
	if a.AcceptorNickname != "" {
		maps["acceptor_nickname"] = a.AcceptorNickname
	}
	if a.ChannelType != "" {
		maps["channel_type"] = a.ChannelType
	}
	if a.AppendInfo != "" {
		maps["append_info"] = a.AppendInfo
	}
	if a.ClientIp != "" {
		maps["client_ip"] = a.ClientIp
	}

	return maps
}
