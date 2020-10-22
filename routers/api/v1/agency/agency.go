package agency

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/acceptor_service"
	"QuickPass/service/bill_service"
	"QuickPass/service/merchant_service"
	"github.com/gin-gonic/gin"
)

// @Summary 代理平台统计
// @Description 代理平台统计
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param start_time query int false "开始日期 yyyy-mm-dd hh:mm:ss"
// @Param end_time query string false "结束日期 yyyy-mm-dd hh:mm:ss"
// @Success 200 {object}  response.RespAgencyStatistics "response.RespAgencyStatistics"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getAgencyStatistics [get]
func GetAgencyStatistics(c *gin.Context) {
	var (
		form request.ReqGetAgencyStatisticsForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 1.已有商家
	merchantService := merchant_service.Merchant{
		Agency: account.Agency,
	}
	merchantCount, err := merchantService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_MERCHANT_FAIL, err.Error())
		return
	}

	// 2.已有承兑
	acceptorService := acceptor_service.Acceptor{
		Agency: account.Agency,
	}
	acceptorCount, err := acceptorService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ACCEPTOR_FAIL, err.Error())
		return
	}

	// 3.承兑总数额
	billService := bill_service.Bill{}
	totalAmountBill, err := billService.RolesBillStatistics(account.Agency, form.StartTime, form.EndTime, constant.ACCEPTOR, constant.EXPEND, constant.BusinessTypeBuy)
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_BILL_FAIL, err.Error())
		return
	}
	// 4.手续费收入
	feeAmountBill, err := billService.RolesBillStatistics(account.Agency, form.StartTime, form.EndTime, constant.AGENCY, constant.INCOME, constant.BusinessTypeBuyFee)
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_BILL_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespAgencyStatistics{
		AcceptorCount: acceptorCount,
		MerchantCount: merchantCount,
		TotalAmount:   totalAmountBill.Amount,
		FeeAmount:     feeAmountBill.Amount,
	})
}
