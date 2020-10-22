package abnormal_order

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/abnormal_order_service"
	"QuickPass/service/merchant_service"
	"QuickPass/service/user_service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取承兑人异常单信息列表
// @Description 获取承兑人异常单信息列表
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param abnormal_order_id query string false "异常单ID"
// @Param abnormal_order_type query int false "异常类型（1：超时取消，2：未生成订单， 3：订单金额不符）"
// @Param abnormal_order_status query int false "订单状态（1：未处理  2：已处理  3：处理中）"
// @Param max_amount query string false "最大金额"
// @Param min_amount query int false "最小金额"
// @Param start_time query int false "开始时间"
// @Param end_time query string false "结束时间"
// @Param finish_start_time query int false "完成开始时间"
// @Param finish_end_time query string false "完成结束时间"
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Success 200 {object}  response.RespAbnormalOrderInfoList
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/getAbnormalOrderInfoList [get]
func GetAcceptorAbnormalOrderInfoList(c *gin.Context) {
	var (
		form request.ReqGetAcceptorAbnormalOrderInfoList
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	abnormalService := abnormal_order_service.AbnormalOrder{
		Agency:              account.Agency,
		Username:            account.Username,
		AbnormalOrderNo:     form.AbnormalOrderNo,
		AbnormalOrderType:   form.AbnormalOrderType,
		AbnormalOrderStatus: form.AbnormalOrderStatus,
		PageNum:             offset,
		PageSize:            limit,
		MaxAmount:           form.MaxAmount,
		MinAmount:           form.MinAmount,
		StartTime:           form.StartTime,
		EndTime:             form.EndTime,
		FinishStartTime:     form.FinishStartTime,
		FinishEndTime:       form.FinishEndTime,
	}

	// 1.查询匹配条件总数
	total, err := abnormalService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ABNORMAL_ORDER_FAIL, err.Error())
		return
	}

	// 2.分页获取异常单信息
	abnormalOrders, err := abnormalService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ABNORMAL_ORDERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespAbnormalOrderInfoList{
		Total: total,
		List:  abnormalOrders,
	})
}

// @Summary 新建异常单
// @Description 新建异常单
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqAddAcceptorAbnormalOrderFrom body request.ReqAddAcceptorAbnormalOrderFrom true  "request.ReqAddAcceptorAbnormalOrderFrom"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/addAbnormalOrder [post]
func AddAcceptorAbnormalOrder(c *gin.Context) {
	var (
		form request.ReqAddAcceptorAbnormalOrderFrom
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

	abnormalOrderService := abnormal_order_service.AbnormalOrder{
		Agency:            account.Agency,
		Username:          account.Username,
		Nickname:          user.FullName,
		Channel:           form.Channel,
		Amount:            form.Amount,
		AcceptCardAccount: form.AcceptCardAccount,
		AcceptCardNo:      form.AcceptCardNo,
		AcceptCardBank:    form.AcceptCardBank,
		PaymentDate:       util.JSONTimeParse(form.PaymentDate),
	}
	err = abnormalOrderService.Add()
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_ABNORMAL_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 获取代理异常单信息列表
// @Description 获取代理异常单信息列表
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param abnormal_order_id query string true "异常单ID"
// @Param abnormal_order_type query int true "异常类型（1：无，2：超时取消，3：未生成订单， 4：订单金额不符）"
// @Param abnormal_order_status query int true "订单状态（1：未处理  2：已处理）"
// @Param max_amount query string true "最大金额"
// @Param min_amount query int true "最小金额"
// @Param start_time query int true "开始时间"
// @Param end_time query string true "结束时间"
// @Param start_page query int true "起始页"
// @Param page_size query int true "页面大小"
// @Success 200 {object}  response.RespAbnormalOrderInfoList
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getAbnormalOrderInfoList [get]
func GetAgencyAbnormalOrderInfoList(c *gin.Context) {
	var (
		form request.ReqGetAgencyAbnormalOrderInfoList
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	abnormalService := abnormal_order_service.AbnormalOrder{
		Agency:              account.Agency,
		Nickname:            form.AcceptorName,
		AbnormalOrderNo:     form.AbnormalOrderNo,
		AbnormalOrderType:   form.AbnormalOrderType,
		AbnormalOrderStatus: form.AbnormalOrderStatus,
		PageNum:             offset,
		PageSize:            limit,
		MaxAmount:           form.MaxAmount,
		MinAmount:           form.MinAmount,
		StartTime:           form.StartTime,
		EndTime:             form.EndTime,
	}

	// 1.查询匹配条件总数
	total, err := abnormalService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ABNORMAL_ORDER_FAIL, err.Error())
		return
	}

	// 2.分页获取异常单信息
	abnormalOrders, err := abnormalService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ABNORMAL_ORDERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespAbnormalOrderInfoList{
		Total: total,
		List:  abnormalOrders,
	})
}

// @Summary 更新异常单生成订单
// @Description 更新异常单生成订单
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqUpdateAbnormalOrderFrom body request.ReqUpdateAbnormalOrderFrom true  "request.ReqUpdateAbnormalOrderFrom"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateAbnormalOrder [post]
func UpdateAbnormalOrder(c *gin.Context) {
	var (
		form request.ReqUpdateAbnormalOrderFrom
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 3.判断同一个代理下的商家是否存在
	merchantService := merchant_service.Merchant{}
	exist, err := merchantService.ExistByAgencyAndUsername(account.Agency, form.Merchant)
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_MERCHANT_FAIL, err.Error())
		return
	}
	if !exist {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_MERCHANT, "")
		return
	}

	abnormalOrderService := abnormal_order_service.AbnormalOrder{}
	err = abnormalOrderService.UpdateAbnormalOrder(account.Agency, form.AbnormalOrderNo,
		form.Merchant, form.Member, form.AppendInfo, c.ClientIP(),
		form.AbnormalOrderType, form.Amount)
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 代理异常订单中所有承兑人的列表
// @Description 代理异常订单中所有承兑人的列表
// @Tags 代理平台
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespGetAcceptorGroup
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getAbnormalOrderAcceptorGroup [get]
func GetAbnormalAcceptorGroup(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	abnormalOrderService := abnormal_order_service.AbnormalOrder{}
	group, err := abnormalOrderService.GetAcceptorGroup(account.Agency)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ABNORMAL_ORDERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespGetAcceptorGroup{
		Total: int64(len(group)),
		List:  group,
	})
}
