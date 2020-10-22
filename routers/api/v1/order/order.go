package order

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/mq"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/channel_service"
	"QuickPass/service/match_cache_service"
	"QuickPass/service/merchant_service"
	"QuickPass/service/order_service"
	"QuickPass/service/user_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// @Summary 买入
// @Description 买入
// @Tags 订单
// @accept json
// @Produce  json
// @Param ReqBuyForm body request.ReqBuyForm true "request.ReqBuyForm"
// @Success 200 {object}  response.RespOrderBuy
// @Failure 500 {object}  app.Response
// @Router /api/v1/order/buy [post]
func Buy(c *gin.Context) {
	var (
		form request.ReqBuyForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	orderService := order_service.Order{}
	// 1.ip过滤
	ipOrders, err := orderService.GetCurrentIPOrder(form.Agency, c.ClientIP())
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_ORDER_FAIL, err.Error())
		return
	}
	//该IP有多笔订单
	if len(ipOrders) > 2 {
		app.ErrorResp(c, e.ERROR_IP_RISK, "")
		return
	}

	// 2.判断支付通道
	channelService := channel_service.Channel{}
	channel := channelService.Get(form.Agency, form.PayType)
	if channel == nil || channel.IfOpen == constant.SwitchClose {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_CHANNEL, "")
		return
	}
	if form.Amount < channel.BuyMin || form.Amount > channel.BuyMax {
		app.ErrorResp(c, e.ERROR, "下单金额有误")
		return
	}

	// 3.判断同一个代理下的商家是否存在
	merchantService := merchant_service.Merchant{}
	merchant, err := merchantService.Get(form.Agency, form.Merchant)
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_MERCHANT_FAIL, err.Error())
		return
	}
	if merchant == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_MERCHANT, "")
		return
	}
	// 判断商家可买入状态
	if merchant.BuyStatus == constant.SwitchClose {
		app.ErrorResp(c, e.MerchantBuyStatusClosed, "")
		return
	}
	// 4.设置订单超时监听
	//订单超时时间
	duration := time.Minute * time.Duration(setting.AppSetting.OrderTimeoutMinute)

	// 5.查询当前用户有正在进行中的订单
	memberOrder, err := orderService.GetCurrentMemberOrder(form.Agency, form.Merchant, form.Member)
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_ORDER_FAIL, err.Error())
		return
	}
	if memberOrder != nil {
		// 不是银行卡的需要查询收款码
		cardImg := ""
		if form.PayType != constant.BANK_CARD {
			cardImg = fmt.Sprintf("http://%s/%s", setting.MinioSetting.PreUrl, memberOrder.CardImg)
		}
		now := util.JSONTimeNow()
		app.SuccessRespByCode(c, e.ERROR_HAS_UNFINISH_ORDER, response.RespOrderBuy{
			OrderNo:     memberOrder.OrderNo,
			CardNo:      memberOrder.CardNo,
			CardAccount: memberOrder.CardAccount,
			CardBank:    memberOrder.CardBank,
			CardImg:     cardImg,
			CardSubBank: memberOrder.CardSubBank,
			CreateTime:  memberOrder.CreateTime,
			CurrentTime: &now,
			ExpirationTime: util.JSONTime{
				Time: memberOrder.CreateTime.Add(duration),
			},
		})
		return
	}

	// 6、创建订单
	order, err := orderService.CreateBuy(form.Agency,
		form.Member,
		form.Merchant,
		form.PayType,
		c.ClientIP(),
		form.Amount,
		constant.SubmitTypeDirect,
		"",
		merchant.ReturnUrl,
		merchant.NotifyUrl,
		"")
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_ORDER_FAIL, err.Error())
		return
	}

	// 7、查询出可匹配的卡id
	context := match_cache_service.NewMatchContext("")
	card, err := context.Match(order.Agency, order.OrderNo, form.Amount, form.PayType)
	if err != nil {
		// 7.1、订单失败
		err2 := orderService.UpdateOrderFailed(order.Agency, order.OrderNo)
		if err2 != nil {
			logf.Error(err2)
		}
		app.ErrorResp(c, e.ERROR_MATCH_FAILED, err.Error())
		return
	}

	// 8、插入消息到队列中，通知承兑人
	order, _ = orderService.GetByOrderNo(order.Agency, order.OrderNo)
	// 发送订单状态改变通知
	go mq.SendOrderStatusChange(order)

	app.SuccessResp(c, response.RespOrderBuy{
		OrderNo:     order.OrderNo,
		CardNo:      card.CardNo,
		CardAccount: card.CardAccount,
		CardBank:    card.CardBank,
		CardImg:     card.CardImg,
		CardSubBank: card.CardSubBank,
		CreateTime:  order.CreateTime,
		ExpirationTime: util.JSONTime{
			Time: order.CreateTime.Add(duration),
		},
	})
}

// @Summary 订单状态查询
// @Description 订单状态查询
// @Tags 订单
// @accept json
// @Produce  json
// @Param agency query string tue "代理"
// @Param order_no query string true "订单号"
// @Success 200 {object}  response.RespGetOrderStatus "response.RespGetOrderStatus"
// @Failure 500 {object}  app.Response
// @Router /api/v1/order/getStatusByOrderNo [get]
func GetStatusByOrderNo(c *gin.Context) {
	var (
		form request.ReqGetStatusByOrderNoForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	orderService := order_service.Order{}
	order, err := orderService.GetByOrderNo(form.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDER_FAIL, err.Error())
		return
	}
	if order == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ORDER, "")
		return
	}
	merchantOrderService := order_service.MerchantOrder{}
	merchantOrder, err := merchantOrderService.GetByOrderNo(form.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDER_FAIL, err.Error())
		return
	}
	if merchantOrder == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ORDER, "")
		return
	}

	// 不是银行卡的需要查询收款码
	cardImg := ""
	if order.ChannelType != constant.BANK_CARD {
		cardImg = fmt.Sprintf("http://%s/%s", setting.MinioSetting.PreUrl, order.CardImg)
	}
	now := util.JSONTimeNow()
	//订单超时时间
	duration := time.Minute * time.Duration(setting.AppSetting.OrderTimeoutMinute)

	app.SuccessResp(c, response.RespGetOrderStatus{
		Agency:      order.Agency,
		OrderNo:     order.OrderNo,
		OrderType:   order.OrderType,
		OrderStatus: order.OrderStatus,
		Amount:      order.Amount,
		ChannelType: order.ChannelType,
		FinishTime:  order.FinishTime,
		CardNo:      order.CardNo,
		CardAccount: order.CardAccount,
		CardBank:    order.CardBank,
		CardImg:     cardImg,
		CardSubBank: order.CardSubBank,
		ReturnUrl:   merchantOrder.ReturnUrl,
		CreateTime:  order.CreateTime,
		CurrentTime: &now,
		ExpirationTime: util.JSONTime{
			Time: order.CreateTime.Add(duration),
		},
	})
}

// @Summary 取消买单
// @Description 取消买单
// @Tags 订单
// @accept json
// @Produce  json
// @Param ReqCancelBuyForm body request.ReqCancelBuyForm true "request.ReqCancelBuyForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/order/cancelBuy [post]
func CancelBuy(c *gin.Context) {
	var (
		form request.ReqCancelBuyForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 1、取消订单
	orderService := order_service.Order{}
	err = orderService.CancelBuyOrder(form.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 确认付款
// @Description 确认付款
// @Tags 代理平台
// @accept json
// @Produce json
// @Param ReqConfirmSellPayForm body request.ReqConfirmSellPayForm true "request.ReqConfirmSellPayForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/confirmSellPay [post]
func ConfirmSellPay(c *gin.Context) {
	var (
		form request.ReqConfirmSellPayForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 1、修改订单状态
	orderService := order_service.Order{}
	err = orderService.ConfirmSellOrder(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	// 2、通知商家(您有新的收款待放行)
	order, _ := orderService.GetByOrderNo(account.Agency, form.OrderNo)
	// 发送订单状态改变通知
	go mq.SendOrderStatusChange(order)

	app.SuccessResp(c, nil)
}

// @Summary 确认付款
// @Description 确认付款
// @Tags 订单
// @accept json
// @Produce json
// @Param ReqConfirmPayForm body request.ReqConfirmPayForm true "order_no:订单号"
// @Success 200 {object}  response.RespConfirmPay
// @Failure 500 {object}  app.Response
// @Router /api/v1/order/confirmPay [post]
func ConfirmPay(c *gin.Context) {
	var (
		form request.ReqConfirmPayForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	// 1、修改订单状态
	orderService := order_service.Order{}
	err = orderService.ConfirmBuyOrder(form.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}
	// 2、通知承兑人(您有新的收款待放行)
	order, _ := orderService.GetByOrderNo(form.Agency, form.OrderNo)
	// 发送订单状态改变通知
	go mq.SendOrderStatusChange(order)

	//订单超时时间
	duration := time.Minute * time.Duration(setting.AppSetting.OrderTimeoutMinute)

	app.SuccessResp(c, response.RespConfirmPay{
		ExpirationTime: util.JSONTime{
			Time: time.Now().Add(duration),
		},
	})
}

// @Summary 商家放行卖单
// @Description 商家放行卖单
// @Tags 商家
// @accept json
// @Produce json
// @Param ReqMerchantDischargeSellForm body request.ReqMerchantDischargeSellForm true "request.ReqMerchantDischargeSellForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/merchantDischargeSell [post]
func MerchantDischargeSell(c *gin.Context) {
	var (
		form request.ReqMerchantDischargeSellForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	userService := user_service.User{}
	user, err := userService.GetByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	// 1、核对交易密码
	if user.TradeKey == "" || !util.MD5Equals(form.TradeKey, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}

	// 2、修改订单状态  where agency order_id to_user_name
	orderService := order_service.Order{}
	err = orderService.DischargeSellOrder(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 取消卖单
// @Description 取消卖单
// @Tags 商家
// @accept json
// @Produce  json
// @Param request body request.ReqCancelSellForm true "order_no:订单号"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/cancelSell [post]
func CancelSell(c *gin.Context) {
	var (
		form request.ReqCancelSellForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 1、修改订单状态  where agency order_id
	service := order_service.Order{}
	order, err := service.GetByOrderNo(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDER_FAIL, err.Error())
		return
	}
	if order == nil || order.FromUserName != account.Username || order.Agency != account.Agency {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ORDER, "")
		return
	}

	err = service.CancelSellOrder(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	// 卖单通知代理

	app.SuccessResp(c, nil)
}

// @Summary 订单查询
// @Description 订单查询
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param order_no query int false "订单号"
// @Param order_status query int false "订单状态"
// @Param channel query string false "交易通道(BANK_CARD,ALIPAY,WECHAT)"
// @Param max_amount query string false "最大金额"
// @Param min_amount query string false "最小金额"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param finish_start_time query string false "完成开始时间"
// @Param finish_end_time query string false "完成结束时间"
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Success 200 {object}  response.RespBuyOrderInfoList "response.RespBuyOrderInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/getBuyOrderInfoList [get]
func GetBuyOrderInfoList(c *gin.Context) {
	var (
		form request.ReqBuyOrderInfoList
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	orderService := order_service.Order{
		Agency:          account.Agency,
		ToUserName:      account.Username,
		OrderNo:         form.OrderNo,
		OrderType:       constant.ORDER_TYPE_BUY,
		OrderStatus:     form.OrderStatus,
		ChannelType:     form.Channel,
		StartTime:       form.StartTime,
		EndTime:         form.EndTime,
		MaxAmount:       form.MaxAmount,
		MinAmount:       form.MinAmount,
		FinishStartTime: form.FinishStartTime,
		FinishEndTime:   form.FinishEndTime,
		PageNum:         offset,
		PageSize:        limit,
	}

	total, err := orderService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ORDER_FAIL, err.Error())
		return
	}

	orders, err := orderService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespBuyOrderInfoList{
		Total: total,
		List:  orders,
	})
}

// @Summary 订单查询
// @Description 订单查询
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param order_id query int false "订单ID"
// @Param order_status query int false "订单状态"
// @Param channel query string false "交易通道(BANK_CARD,ALIPAY,WECHAT)"
// @Param max_amount query string false "最大金额"
// @Param min_amount query string false "最小金额"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param acceptor_name query string false "所属承兑人"
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Success 200 {object}  response.RespBuyOrderInfoList "response.RespBuyOrderInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getOrderInfoList [get]
func GetAgencyOrderInfoList(c *gin.Context) {
	var (
		form request.ReqGetAgencyOrderInfoList
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	orderService := order_service.Order{
		Agency: account.Agency,
		//ToUserName:  form.AcceptorName,
		MerchantOrderNo:  form.MerchantOrderNo,
		AcceptorNickname: form.AcceptorName,
		OrderNo:          form.OrderNo,
		OrderType:        constant.ORDER_TYPE_BUY,
		OrderStatus:      form.OrderStatus,
		ChannelType:      form.Channel,
		StartTime:        form.StartTime,
		EndTime:          form.EndTime,
		MaxAmount:        form.MaxAmount,
		MinAmount:        form.MinAmount,
		PageNum:          offset,
		PageSize:         limit,
	}

	total, err := orderService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ORDER_FAIL, err.Error())
		return
	}

	orders, err := orderService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespAgencyBuyOrderInfoList{
		Total: total,
		List:  orders,
	})
}

// @Summary 代理订单中所有承兑人的列表
// @Description 代理订单中所有承兑人的列表
// @Tags 代理平台
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespGetAcceptorGroup "response.RespGetAcceptorGroup"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getAcceptorGroup [get]
func GetAcceptorGroup(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	orderService := order_service.Order{}
	group, err := orderService.GetAcceptorGroup(account.Agency)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespGetAcceptorGroup{
		Total: int64(len(group)),
		List:  group,
	})
}

// @Summary 代理放行订单
// @Description 代理放行订单
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAgencyDischargeForm body request.ReqAgencyDischargeForm true "order_no:订单号"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/agencyDischarge [post]
func AgencyDischarge(c *gin.Context) {
	var (
		form request.ReqAgencyDischargeForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 2、订单放行
	orderService := order_service.Order{}
	err = orderService.DischargeBuyOrder(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 代理放行超时取消的订单
// @Description 代理放行超时取消的订单
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAgencyDischargeCancelOrderForm body request.ReqAgencyDischargeCancelOrderForm true "order_no:订单号"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/agencyDischargeCancelOrder [post]
func AgencyDischargeCancelOrder(c *gin.Context) {
	var (
		form request.ReqAgencyDischargeCancelOrderForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 已取消的订单放行
	orderService := order_service.Order{}
	err = orderService.DischargeCancelOrder(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)

}

// @Summary 代理再次回调商家
// @Description 代理再次回调商家
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAgencyDischargeForm body request.ReqAgencyAgainCallback true "order_no:订单号"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/againCallback [post]
func AgencyAgainCallback(c *gin.Context) {
	var (
		form request.ReqAgencyAgainCallback
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 订单查询
	orderService := order_service.Order{}
	order, err := orderService.GetByOrderNo(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_AGAIN_CALLBACK, err.Error())
		return
	}
	// 开始回调
	orderService.AgainCallbackMerchant(order, form.CallbackUrl)
	app.SuccessResp(c, nil)
}

// @Summary 代理卖单未处理数
// @Description 代理卖单未处理数
// @Tags 代理平台
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespGetSellOrderUnprocessed "response.RespGetSellOrderUnprocessed"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getSellOrderUnprocessed [get]
func GetSellOrderUnprocessed(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	orderService := order_service.Order{
		Agency:      account.Agency,
		OrderType:   constant.ORDER_TYPE_SELL,
		OrderStatus: constant.ORDER_STATUS_WAIT_PAY,
	}

	count, err := orderService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespGetSellOrderUnprocessed{
		UnprocessedCount: count,
	})
}
