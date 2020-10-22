package acceptor

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/acceptor_card_service"
	"QuickPass/service/acceptor_service"
	"QuickPass/service/bill_service"
	"QuickPass/service/fund_service"
	"QuickPass/service/user_service"
	"github.com/gin-gonic/gin"
)

// @Summary 增加承兑人
// @Description 增加承兑人
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAddAcceptorForm body request.ReqAddAcceptorForm true  "user_name:承兑人账户名 full_name:承兑人姓名 phone_number:电话号码 address:地址"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/AddAcceptor [post]
func AddAcceptor(c *gin.Context) {
	var (
		form request.ReqAddAcceptorForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorService := acceptor_service.Acceptor{
		Acceptor:    form.UserName,
		Agency:      account.Agency,
		Password:    form.Password,
		FullName:    form.FullName,
		PhoneNumber: form.PhoneNumber,
		Address:     form.Address,
	}
	exist, err := acceptorService.ExistByAgencyAndUsername()
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_ACCEPTOR_FAIL, err.Error())
		return
	}
	if exist {
		app.ErrorResp(c, e.ExistUsername, "")
		return
	}
	if err := acceptorService.Add(); err != nil {
		app.ErrorResp(c, e.ERROR_ADD_ACCEPTOR_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 更新承兑权限
// @Description 更新承兑权限
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqUpdateAcceptSwitchFrom body request.ReqUpdateAcceptSwitchFrom true "request.ReqUpdateAcceptSwitchFrom"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/updateAcceptSwitch [post]
func UpdateAcceptSwitch(c *gin.Context) {
	var (
		form request.ReqUpdateAcceptSwitchFrom
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorService := acceptor_service.Acceptor{}
	err = acceptorService.UpdateAcceptSwitch(account.Agency, account.Username, form.AcceptSwitch)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ACCEPTOR_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 代理更新承兑权限
// @Description 代理更新承兑权限
// @Tags 代理
// @accept json
// @Produce  json
// @Param ReqAgencyUpdateAcceptSwitchFrom body request.ReqAgencyUpdateAcceptSwitchFrom true "request.ReqAgencyUpdateAcceptSwitchFrom"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateAcceptSwitch [post]
func AgencyUpdateAcceptSwitch(c *gin.Context) {
	var (
		form request.ReqAgencyUpdateAcceptSwitchFrom
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorService := acceptor_service.Acceptor{}
	err = acceptorService.UpdateAcceptSwitch(account.Agency, form.Username, form.AcceptSwitch)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ACCEPTOR_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 代理更新承兑账户状态
// @Description 代理更新承兑账户状态
// @Tags 代理
// @accept json
// @Produce  json
// @Param ReqAgencyUpdateAcceptStatusFrom body request.ReqAgencyUpdateAcceptStatusFrom true "request.ReqAgencyUpdateAcceptStatusFrom"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateAcceptStatus [post]
func AgencyUpdateAcceptStatus(c *gin.Context) {
	var (
		form request.ReqAgencyUpdateAcceptStatusFrom
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorService := acceptor_service.Acceptor{}
	err = acceptorService.UpdateAcceptStatus(account.Agency, form.Username, form.AcceptStatus)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ACCEPTOR_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 更新承兑权限
// @Description 更新承兑权限
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqUpdateIfAutoAcceptFrom body request.ReqUpdateIfAutoAcceptFrom true "request.ReqUpdateIfAutoAcceptFrom"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/updateIfAutoAccept [post]
func UpdateIfAutoAccept(c *gin.Context) {
	var (
		form request.ReqUpdateIfAutoAcceptFrom
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorService := acceptor_service.Acceptor{}
	err = acceptorService.UpdateIfAutoAccept(account.Agency, account.Username, form.IfAutoAccept)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ACCEPTOR_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 代理更新承兑权限
// @Description 代理更新承兑权限
// @Tags 代理
// @accept json
// @Produce  json
// @Param ReqAgencyUpdateIfAutoAcceptFrom body request.ReqAgencyUpdateIfAutoAcceptFrom true "request.ReqAgencyUpdateIfAutoAcceptFrom"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateIfAutoAccept [post]
func AgencyUpdateIfAutoAccept(c *gin.Context) {
	var (
		form request.ReqAgencyUpdateIfAutoAcceptFrom
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorService := acceptor_service.Acceptor{}
	err = acceptorService.UpdateIfAutoAccept(account.Agency, form.Username, form.IfAutoAccept)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ACCEPTOR_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 承兑管理统计
// @Description 承兑管理统计
// @Tags 代理平台
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespAcceptorStatistics "response.RespAcceptorStatistics"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getAcceptorStatistics [get]
func GetAcceptorStatistics(c *gin.Context) {
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	//已有卡商
	acceptorService := acceptor_service.Acceptor{Agency: account.Agency}
	count1, err := acceptorService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ACCEPTOR_FAIL, err.Error())
		return
	}

	//当前在线
	acceptorService.AcceptSwitch = constant.SwitchOpen
	count2, err := acceptorService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ACCEPTOR_FAIL, err.Error())
		return
	}

	//自动承兑（开启）
	acceptorService.AcceptSwitch = 0
	acceptorService.IfAutoAccept = constant.Auto
	count3, err := acceptorService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ACCEPTOR_FAIL, err.Error())
		return
	}

	// 今日承兑数额
	billService := bill_service.Bill{}
	statistics, err := billService.RolesBillStatistics(account.Agency, util.TodayStartTimeStr(), util.TomorrowStartTimeStr(), constant.ACCEPTOR, constant.EXPEND, constant.BusinessTypeBuy)
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ACCEPTOR_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespAcceptorStatistics{
		TotalCount:        count1,
		OnLineCount:       count2,
		AutoAcceptCount:   count3,
		TodayAcceptAmount: statistics.Amount,
	})
}

// @Summary 获取承兑人信息列表
// @Description 获取承兑人信息列表
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param start_page query int true "起始页"
// @Param page_size query int true "页面大小"
// @Success 200 {object}  response.RespAcceptorInfoList "response.RespAcceptorInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getAcceptorInfoList [get]
func GetAcceptorInfoList(c *gin.Context) {
	var (
		form request.ReqGetAcceptorInfoListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	acceptorService := acceptor_service.Acceptor{
		Agency:   account.Agency,
		PageNum:  offset,
		PageSize: limit,
	}

	// 1.查询匹配条件总数
	total, err := acceptorService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ACCEPTOR_FAIL, err.Error())
		return
	}

	// 2.分页获取承兑人信息
	acceptors, err := acceptorService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ACCEPTORS_FAIL, err.Error())
		return
	}

	startTime := util.TodayStartTimeStr()
	endTime := util.TomorrowStartTimeStr()
	// 3.挨个查询承兑人今日承兑金额
	billService := bill_service.Bill{}
	cardService := acceptor_card_service.AcceptorCard{}
	userService := user_service.User{}
	for _, item := range acceptors {
		//今日承兑金额
		bill, _ := billService.UserBillStatistics(item.Agency, item.Acceptor, startTime, endTime, constant.EXPEND)
		item.TodayAcceptAmount = bill.Amount

		count, _ := cardService.CountCardType(item.Agency, item.Acceptor)
		item.CountBankCard = count.CountBankCard
		item.CountAliPay = count.CountAliPay
		item.CountWeChat = count.CountWeChat

		user, _ := userService.GetByAgencyAndUsername(item.Agency, item.Acceptor)
		item.Nickname = user.FullName
	}

	app.SuccessResp(c, response.RespAcceptorInfoList{
		Total: total,
		List:  acceptors,
	})
}

// @Summary 获取承兑人信息
// @Description 获取承兑人信息
// @Tags 承兑人
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespAcceptorInfo "response.RespAcceptorInfo"
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/getAcceptorInfo [get]
func GetAcceptorInfo(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorService := acceptor_service.Acceptor{}

	acceptor, err := acceptorService.GetByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ACCEPTOR_FAIL, err.Error())
		return
	}
	if acceptor == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ACCEPTOR, "")
		return
	}

	fundService := fund_service.Fund{}
	fund, _ := fundService.GetByAgencyAndName(account.Agency, account.Username)

	app.SuccessResp(c, response.RespAcceptorInfo{
		TotalAmount:     fund.AvailableAmount + fund.FrozenAmount,
		AvailableAmount: fund.AvailableAmount,
		FrozenAmount:    fund.FrozenAmount,
		AcceptSwitch:    acceptor.AcceptSwitch,
		IfAutoAccept:    acceptor.IfAutoAccept,
	})
}
