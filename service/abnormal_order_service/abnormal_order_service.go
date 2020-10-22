package abnormal_order_service

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/mq"
	"QuickPass/pkg/util"
	"fmt"
)

type AbnormalOrder struct {
	//
	Id int64
	// 代理
	Agency string
	// 用户名
	Username string
	// 昵称
	Nickname string
	// 收款通道({“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"})
	Channel string
	// 原始订单号
	OriginalOrderNo string
	// 异常单号
	AbnormalOrderNo string
	// 异常类型（1: 未知, 2：超时取消，3：未生成订单， 4：订单金额不符）
	AbnormalOrderType int64
	// 订单状态（1：未处理  2：已处理  3：处理中）
	AbnormalOrderStatus int64
	// 收款账户
	AcceptCardAccount string
	// 收款卡号
	AcceptCardNo string
	// 银行名称
	AcceptCardBank string
	// 金额
	Amount int64
	// 摘要
	AppendInfo string
	// 付款日期
	PaymentDate util.JSONTime
	// 完成时间
	FinishTime *util.JSONTime
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int

	//----------------------------
	// 最大金额
	MaxAmount string
	// 最小金额
	MinAmount string
	// 开始时间
	StartTime string
	// 结束时间
	EndTime string
	// 完成开始时间
	FinishStartTime string
	// 完成结束时间
	FinishEndTime string
}

func (a *AbnormalOrder) Add() error {
	abnormalOrder := &models.AbnormalOrder{
		Agency:              a.Agency,
		Username:            a.Username,
		Nickname:            a.Nickname,
		OriginalOrderNo:     a.OriginalOrderNo,
		AbnormalOrderNo:     util.GetUniqueNo(),
		AbnormalOrderType:   a.AbnormalOrderType,
		AcceptCardAccount:   a.AcceptCardAccount,
		Amount:              a.Amount,
		AppendInfo:          a.AppendInfo,
		PaymentDate:         a.PaymentDate,
		Channel:             a.Channel,
		AcceptCardNo:        a.AcceptCardNo,
		AcceptCardBank:      a.AcceptCardBank,
		AbnormalOrderStatus: constant.Unprocessed,
	}

	session := models.NewSession()
	if err := models.NewAbnormalOrderModel(session).AddAbnormalOrder(abnormalOrder); err != nil {
		return err
	}

	go mq.SendAbnormalOrderStatusChange(abnormalOrder)
	return nil
}

func (a *AbnormalOrder) UpdateAbnormalOrder(agency, abnormalOrderNo, merchant, member, appendInfo, ip string, abnormalOrderType, amount int64) error {
	session := models.NewSession()
	orderModel := models.NewOrderModel(session)
	merchantOrderModel := models.NewMerchantOrderModel(session)
	abnormalOrderModel := models.NewAbnormalOrderModel(session)
	merchantModel := models.NewMerchantModel(session)
	fundModel := models.NewFundModel(session)
	acceptorCardModel := models.NewAcceptorCardModel(session)

	merchantItem, err := merchantModel.GetMerchant(agency, merchant)
	if err != nil {
		return err
	}

	// 获取异常单信息
	abnormalOrder, err := abnormalOrderModel.GetAbnormalOrder(agency, abnormalOrderNo)
	if err != nil {
		return err
	}
	if abnormalOrder == nil {
		return errors.New(fmt.Sprintf("GetAbnormalOrder agency[%s] orderNo[%s] not exist", agency, abnormalOrderNo))
	}

	// 原始订单号
	originalMerchantOrder, _ := merchantOrderModel.GetMerchantOrder(agency, abnormalOrder.OriginalOrderNo)

	// 获取卡信息
	card, err := acceptorCardModel.GetAcceptorCardByCardNo(agency, abnormalOrder.Username, abnormalOrder.AcceptCardNo, abnormalOrder.Channel)
	if err != nil {
		return err
	}
	if card == nil {
		return errors.New(fmt.Sprintf("GetAcceptorCardByCardNo agency[%s] acceptor[%s] cardNo[%s] cardType[%s] not exist", agency, abnormalOrder.Username, abnormalOrder.AcceptCardNo, abnormalOrder.Channel))
	}

	// 1.开始事务
	session.Begin()
	defer session.Rollback()

	uniqueNo := util.GetUniqueNo()
	order := &models.Order{
		Agency:    agency,
		OrderNo:   uniqueNo,
		OrderType: constant.ORDER_TYPE_BUY,
		// 异常单生成的订单都是待放行状态
		OrderStatus:      constant.ORDER_STATUS_WAIT_DISCHARGE,
		AbnormalOrderNo:  abnormalOrder.AbnormalOrderNo,
		MerchantUserName: merchant,
		// 会员名称
		FromUserName: member,
		// 异常单的创建者为收款方：承兑人
		ToUserName:       card.Acceptor,
		AcceptorNickname: abnormalOrder.Nickname,
		CardNo:           card.CardNo,
		CardAccount:      card.CardAccount,
		CardBank:         card.CardBank,
		Amount:           amount,
		ChannelType:      abnormalOrder.Channel,
		ClientIp:         ip,
		SubmitType:       constant.SubmitTypeDirect,
	}
	// 1.创建新订单
	if err := orderModel.AddOrder(order); err != nil {
		return err
	}

	merchantOrderNo := order.OrderNo
	if originalMerchantOrder != nil {
		merchantOrderNo = originalMerchantOrder.MerchantOrderNo
	}
	merchantOrder := &models.MerchantOrder{
		OrderNo:         order.OrderNo,
		Agency:          agency,
		Username:        merchant,
		MerchantOrderNo: merchantOrderNo,
		ReturnUrl:       merchantItem.ReturnUrl,
		CallbackUrl:     merchantItem.NotifyUrl,
		SubmitType:      constant.SubmitTypeDirect,
		CallbackStatus:  constant.CallbackUnprocessed,
		AppendInfo:      appendInfo,
	}
	if err := models.NewMerchantOrderModel(session).AddMerchantOrder(merchantOrder); err != nil {
		return err
	}

	// 2、资金操作，冻结承兑人的币
	_, err = fundModel.Frozen(agency, abnormalOrder.Username, amount)
	if err != nil {
		return err
	}

	// 3、更新卡中的可承兑金额 和匹配时间
	err = acceptorCardModel.SubtractDayAvailableAmount(card.Id, amount)
	if err != nil {
		return err
	}

	// 24.更新异常单
	err = abnormalOrderModel.EditAbnormalOrder(agency, abnormalOrderNo, map[string]interface{}{
		"original_order_no":     uniqueNo,
		"abnormal_order_type":   abnormalOrderType,
		"abnormal_order_status": constant.Processing,
		"amount":                amount,
		"append_info":           appendInfo,
	})
	if err != nil {
		return err
	}

	session.Commit()

	// 发送订单状态改变通知
	go mq.SendOrderStatusChange(order)
	// 获取异常单信息
	abnormalOrder, _ = abnormalOrderModel.GetAbnormalOrder(agency, abnormalOrderNo)
	go mq.SendAbnormalOrderStatusChange(abnormalOrder)
	return nil
}

func (a *AbnormalOrder) Get(agency, abnormalOrderNo string) (*models.AbnormalOrder, error) {
	session := models.NewSession()
	return models.NewAbnormalOrderModel(session).GetAbnormalOrder(agency, abnormalOrderNo)
}

func (a *AbnormalOrder) GetAcceptorGroup(agency string) ([]string, error) {
	session := models.NewSession()
	return models.NewAbnormalOrderModel(session).GetAcceptorGroup(agency)
}

func (a *AbnormalOrder) GetAll() ([]*models.AbnormalOrder, error) {
	session := models.NewSession()
	return models.NewAbnormalOrderModel(session).GetAbnormalOrders(&models.AbnormalOrderListQueryParams{
		StartTime:       a.StartTime,
		EndTime:         a.EndTime,
		FinishStartTime: a.FinishStartTime,
		FinishEndTime:   a.FinishEndTime,
		MaxAmount:       a.MaxAmount,
		MinAmount:       a.MinAmount,
	}, a.PageNum, a.PageSize, a.getMaps())
}

func (a *AbnormalOrder) Delete() error {
	session := models.NewSession()
	return models.NewAbnormalOrderModel(session).DeleteAbnormalOrder(a.Id)
}

func (a *AbnormalOrder) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewAbnormalOrderModel(session).ExistAbnormalOrderByID(a.Id)
}

func (a *AbnormalOrder) Count() (int64, error) {
	session := models.NewSession()
	return models.NewAbnormalOrderModel(session).GetAbnormalOrderTotal(&models.AbnormalOrderListQueryParams{
		StartTime:       a.StartTime,
		EndTime:         a.EndTime,
		FinishStartTime: a.FinishStartTime,
		FinishEndTime:   a.FinishEndTime,
		MaxAmount:       a.MaxAmount,
		MinAmount:       a.MinAmount,
	}, a.getMaps())
}

func (a *AbnormalOrder) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.Username != "" {
		maps["username"] = a.Username
	}
	if a.Nickname != "" {
		maps["nickname"] = a.Nickname
	}
	if a.OriginalOrderNo != "" {
		maps["original_order_no"] = a.OriginalOrderNo
	}
	if a.AbnormalOrderNo != "" {
		maps["abnormal_order_no"] = a.AbnormalOrderNo
	}
	if a.AbnormalOrderType != 0 {
		maps["abnormal_order_type"] = a.AbnormalOrderType
	}
	if a.AbnormalOrderStatus != 0 {
		maps["abnormal_order_status"] = a.AbnormalOrderStatus
	}

	return maps
}
