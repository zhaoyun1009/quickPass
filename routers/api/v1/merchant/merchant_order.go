package merchant

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/order_service"
	"github.com/gin-gonic/gin"
)

// @Summary 买入订单查询
// @Description 买入订单查询
// @Tags 商家
// @accept json
// @Produce  json
// @Param order_no query int false "订单号"
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
// @Router /api/v1/merchant/getBuyOrderInfoList [get]
func GetBuyOrderInfoList(c *gin.Context) {
	var (
		form request.ReqMerchantBuyOrderInfoList
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	orderService := order_service.Order{
		Agency:           account.Agency,
		MerchantOrderNo:  form.MerchantOrderNo,
		MerchantUserName: account.Username,
		OrderNo:          form.OrderNo,
		OrderType:        constant.ORDER_TYPE_BUY,
		OrderStatus:      constant.ORDER_STATUS_FINISHED,
		ChannelType:      form.Channel,
		StartTime:        form.StartTime,
		EndTime:          form.EndTime,
		MaxAmount:        form.MaxAmount,
		MinAmount:        form.MinAmount,
		FinishStartTime:  form.FinishStartTime,
		FinishEndTime:    form.FinishEndTime,
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

	app.SuccessResp(c, response.RespMerchantBuyOrderInfoList{
		Total: total,
		List:  orders,
	})
}

// @Summary 卖出订单查询
// @Description 卖出订单查询
// @Tags 商家
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
// @Success 200 {object}  response.RespSellOrderInfoList "response.RespSellOrderInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/getSellOrderInfoList [get]
func GetSellOrderInfoList(c *gin.Context) {
	var (
		form request.ReqMerchantSellOrderInfoList
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
		FromUserName:    account.Username,
		OrderNo:         form.OrderNo,
		OrderType:       constant.ORDER_TYPE_SELL,
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

	app.SuccessResp(c, response.RespSellOrderInfoList{
		Total: total,
		List:  orders,
	})
}

// @Summary 代理卖单列表
// @Description 代理卖单列表
// @Tags 代理平台
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
// @Success 200 {object}  response.RespSellOrderInfoList "response.RespSellOrderInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getSellOrderInfoList [get]
func GetAgencySellOrderInfoList(c *gin.Context) {
	var (
		form request.ReqAgencySellOrderInfoList
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
		OrderNo:         form.OrderNo,
		OrderType:       constant.ORDER_TYPE_SELL,
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

	app.SuccessResp(c, response.RespSellOrderInfoList{
		Total: total,
		List:  orders,
	})
}
