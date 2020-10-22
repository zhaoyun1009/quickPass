package fund

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/fund_service"
	"QuickPass/service/user_service"
	"github.com/gin-gonic/gin"
)

// @Summary 系统转账
// @Description 系统转账
// @Tags 后台管理
// @accept json
// @Produce  json
// @Param ReqSysTransferForm body request.ReqSysTransferForm true "request.ReqSysTransferForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/backend/management/transfer [post]
func SysTransfer(c *gin.Context) {
	var (
		form request.ReqSysTransferForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 1.获取用户信息
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

	// 2.校验交易密码
	if user.TradeKey == "" || !util.MD5Equals(form.TradeKey, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}

	err = transfer(account.Agency, account.Username, form.ToAgency, form.ToUserName, form.AppendInfo, form.Amount)
	if err != nil {
		app.ErrorResp(c, e.ERROR, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 代理转账
// @Description 代理转账
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAgencyTransferForm body request.ReqAgencyTransferForm true "request.ReqAgencyTransferForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/transfer [post]
func AgencyTransfer(c *gin.Context) {
	var (
		form request.ReqAgencyTransferForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 1.获取用户信息
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

	// 2.校验交易密码
	if user.TradeKey == "" || !util.MD5Equals(form.TradeKey, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}

	// 同一代理之间用户转账
	agency := account.Agency

	err = transfer(agency, account.Username, agency, form.ToUserName, form.AppendInfo, form.Amount)
	if err != nil {
		app.ErrorResp(c, e.ERROR, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

func transfer(fromAgency, fromUsername, toAgency, toUsername, appendInfo string, amount int64) error {
	fundService := fund_service.Fund{}
	return fundService.Transfer(fromAgency, toAgency, fromUsername, toUsername, appendInfo, amount)
}

// @Summary 获取资金账户金额信息
// @Description 获取资金账户金额信息
// @Tags 资金
// @accept json
// @Produce json
// @Success 200 {object}  response.RespGetFundInfo "response.RespGetFundInfo"
// @Failure 500 {object}  app.Response
// @Router /api/v1/fund/getFundInfo [get]
func GetFundInfo(c *gin.Context) {
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	fundService := fund_service.Fund{}
	fund, err := fundService.GetByAgencyAndName(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_FUND_FAIL, err.Error())
		return
	}

	if fund == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_FUND, "")
		return
	}

	app.SuccessResp(c, response.RespGetFundInfo{
		Agency:          fund.Agency,
		UserName:        fund.UserName,
		AvailableAmount: fund.AvailableAmount,
		FrozenAmount:    fund.FrozenAmount,
	})
}
