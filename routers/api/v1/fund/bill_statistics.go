package fund

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/bill_service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取用户每日账单统计
// @Description 获取用户每日账单统计
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param username query int true "用户名"
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Param start_date query int false "开始日期"
// @Param end_date query int false "结束日期"
// @Success 200 {object}  response.RespUserBillStatisticsList "response.RespUserBillStatisticsList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getUserBillStatistics [get]
func GetUserBillStatistics(c *gin.Context) {
	var (
		form request.ReqGetUserBillStatisticsForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	billStatisticsService := bill_service.BillStatistics{
		Agency:    account.Agency,
		Username:  form.Username,
		PageNum:   offset,
		PageSize:  limit,
		StartDate: form.StartDate,
		EndDate:   form.EndDate,
	}
	count, err := billStatisticsService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_BILL_FAIL, err.Error())
		return
	}
	list, err := billStatisticsService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_BILLS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespUserBillStatisticsList{
		Total: count,
		List:  list,
	})
}

// @Summary 商家数据统计和承兑人数据统计
// @Description 商家数据统计和承兑人数据统计
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param role query int true "角色(2：承兑人、3：商家)"
// @Param username query int false "用户名"
// @Param start_date query int false "开始日期"
// @Param end_date query int false "结束日期"
// @Success 200 {object}  response.RespUserByRoleBillStatistics "response.RespUserByRoleBillStatistics"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getUserByRoleBillStatistics [get]
func GetUserByRoleBillStatistics(c *gin.Context) {
	var (
		form request.ReqUserByRoleBillStatisticsForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	billStatisticsService := bill_service.BillStatistics{}
	statistics, err := billStatisticsService.GetUserByRoleBillStatistics(account.Agency, form.Username, form.StartDate, form.EndDate, form.Role)
	if err != nil {
		app.ErrorResp(c, e.ERROR, "服务异常")
		return
	}
	app.SuccessResp(c, response.RespUserByRoleBillStatistics{
		AcceptAmount:     statistics.AcceptAmount,
		AcceptCount:      statistics.AcceptCount,
		WithdrawalAmount: statistics.WithdrawalAmount,
		WithdrawalCount:  statistics.WithdrawalCount,
		RechargeAmount:   statistics.RechargeAmount,
		RechargeCount:    statistics.RechargeCount,
	})
}
