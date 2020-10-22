package fund

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/bill_service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取用户账单
// @Description 获取用户账单
// @Tags 资金
// @accept json
// @Produce json
// @Param bill_no query string false "账单号"
// @Param opposite_user_name query string false "对方交易账户"
// @Param income_expenses_type query int false "收支类型（1：支出，2：收入）"
// @Param start_time query int false "开始日期 yyyy-mm-dd hh:mm:ss"
// @Param end_time query string false "结束日期 yyyy-mm-dd hh:mm:ss"
// @Param max_amount query int false "最大金额"
// @Param min_amount query int false "最小金额"
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Success 200 {object}  response.RespUserBillList "response.RespUserBillList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/fund/bill/getUserBillList [get]
func GetUserBillList(c *gin.Context) {
	var (
		form request.ReqGetUserBillListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	billService := bill_service.Bill{
		BillNo:             form.BillNo,
		IncomeExpensesType: form.IncomeExpensesType,
		OppositeUserName:   form.OppositeUserName,
		OrderNo:            form.OrderNo,
		MerchantOrderNo:    form.MerchantOrderNo,
		Agency:             account.Agency,
		OwnUserName:        account.Username,
		PageNum:            offset,
		PageSize:           limit,

		StartTime: form.StartTime,
		EndTime:   form.EndTime,
		MaxAmount: form.MaxAmount,
		MinAmount: form.MinAmount,
	}

	total, err := billService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_BILL_FAIL, err.Error())
		return
	}

	bills, err := billService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_BILLS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespUserBillList{
		Total: total,
		List:  bills,
	})
}

// @Summary 获取收入/支出统计信息
// @Description 获取收入/支出统计信息
// @Tags 资金
// @accept json
// @Produce json
// @Param income_expenses_type query int true "收支类型（1：支出，2：收入）"
// @Param start_time query int false "开始日期 yyyy-mm-dd hh:mm:ss"
// @Param end_time query string false "结束日期 yyyy-mm-dd hh:mm:ss"
// @Success 200 {object}  response.RespInComeStatistics "response.RespInComeStatistics"
// @Failure 500 {object}  app.Response
// @Router /api/v1/fund/bill/inComeStatistics [get]
func InComeStatistics(c *gin.Context) {
	var (
		form request.ReqInComeStatisticsForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	billService := bill_service.Bill{}
	statistics, err := billService.UserBillStatistics(account.Agency, account.Username, form.StartTime, form.EndTime, form.IncomeExpensesType)
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_BILL_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespInComeStatistics{
		Agency:             statistics.Agency,
		Username:           statistics.OwnUserName,
		Amount:             statistics.Amount,
		IncomeExpensesType: statistics.IncomeExpensesType,
	})
}

// @Summary 代理最近转账信息
// @Description 代理最近转账信息
// @Tags 代理平台
// @accept json
// @Produce json
// @Success 200 {object}  response.RespUserBillList "response.RespUserBillList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/latelyTransferInfo [get]
func GetAgencyLatelyTransferInfo(c *gin.Context) {
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(1, 6)
	billService := bill_service.Bill{
		//转账类型
		BusinessType: constant.BusinessTypeTransfer,
		//支出
		IncomeExpensesType: constant.EXPEND,
		Agency:             account.Agency,
		OwnUserName:        account.Username,
		PageNum:            offset,
		PageSize:           limit,
	}
	bills, err := billService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_BILLS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespUserBillList{
		Total: int64(len(bills)),
		List:  bills,
	})
}

// @Summary 代理获取转账给承兑人的账单
// @Description 代理获取转账给承兑人的账单
// @Tags 代理平台
// @accept json
// @Produce json
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Success 200 {object}  response.RespUserBillList "response.RespUserBillList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getAcceptorBill [get]
func GetAcceptorBill(c *gin.Context) {
	var (
		form request.ReqGetAcceptorBillForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	billService := bill_service.Bill{
		Agency:             account.Agency,
		OwnUserName:        account.Username,
		PageNum:            offset,
		PageSize:           limit,
		BusinessType:       constant.BusinessTypeTransfer,
		IncomeExpensesType: constant.EXPEND,
	}

	total, err := billService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_BILL_FAIL, err.Error())
		return
	}

	bills, err := billService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_BILLS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespUserBillList{
		Total: total,
		List:  bills,
	})
}
