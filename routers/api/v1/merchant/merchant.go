package merchant

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/bill_service"
	"QuickPass/service/fund_service"
	"QuickPass/service/merchant_card_service"
	"QuickPass/service/merchant_rate_service"
	"QuickPass/service/merchant_service"
	"QuickPass/service/order_service"
	"QuickPass/service/user_service"
	"github.com/gin-gonic/gin"
)

// @Summary 添加商家
// @Description 添加商家
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAddMerchantForm body request.ReqAddMerchantForm true "user_name:商家账户名 full_name:商家名称 phone_number:电话号码 address:地址"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/addMerchant [post]
func AddMerchant(c *gin.Context) {
	var (
		form request.ReqAddMerchantForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	merchantService := merchant_service.Merchant{
		MerchantName: form.UserName,
		Agency:       account.Agency,
		FullName:     form.FullName,
		Password:     form.Password,
		PhoneNumber:  form.PhoneNumber,
		Address:      form.Address,
		Rate:         form.Rate,
	}
	exist, err := merchantService.ExistByAgencyAndUsername(account.Agency, form.UserName)
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_MERCHANT_FAIL, err.Error())
		return
	}
	if exist {
		app.ErrorResp(c, e.ExistUsername, "")
		return
	}
	if err := merchantService.Add(); err != nil {
		app.ErrorResp(c, e.ERROR_ADD_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 删除商家
// @Description 删除商家
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqRemoveMerchantForm body request.ReqRemoveMerchantForm true "user_name:商家账户名"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/removeMerchant [post]
func RemoveMerchant(c *gin.Context) {
	var (
		form request.ReqRemoveMerchantForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	merchantService := merchant_service.Merchant{}
	if err := merchantService.Remove(account.Agency, form.UserName); err != nil {
		app.ErrorResp(c, e.ERROR_DELETE_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 商家管理统计
// @Description 商家管理统计
// @Tags 代理平台
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespAcceptorStatistics "response.RespAcceptorStatistics"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getMerchantStatistics [get]
func GetMerchantStatistics(c *gin.Context) {
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 今日承兑数额
	billService := bill_service.Bill{}
	statistics, err := billService.RolesBillStatistics(account.Agency, util.TodayStartTimeStr(), util.TomorrowStartTimeStr(), constant.MERCHANT, constant.INCOME, constant.BusinessTypeBuy)
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespMerchantStatistics{
		TodayAcceptAmount: statistics.Amount,
	})
}

// @Summary 获取商家信息列表
// @Description 获取商家信息列表
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Success 200 {object}  response.RespMerchantInfoList "response.RespMerchantInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getMerchantInfoList [get]
func GetMerchantInfoList(c *gin.Context) {
	var (
		form request.ReqGetMerchantInfoListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	merchantService := merchant_service.Merchant{
		Agency:   account.Agency,
		PageNum:  offset,
		PageSize: limit,
	}

	total, err := merchantService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_MERCHANT_FAIL, err.Error())
		return
	}

	merchants, err := merchantService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_MERCHANTS_FAIL, err.Error())
		return
	}

	startTime := util.TodayStartTimeStr()
	endTime := util.TomorrowStartTimeStr()
	// 3.挨个查询商家今日承兑金额
	billService := bill_service.Bill{}
	userService := user_service.User{}
	merchantRateService := merchant_rate_service.MerchantRate{}
	for _, item := range merchants {
		bill, _ := billService.UserBillStatistics(item.Agency, item.MerchantName, startTime, endTime, constant.INCOME)
		item.TodayAcceptAmount = bill.Amount

		user, _ := userService.GetByAgencyAndUsername(item.Agency, item.MerchantName)
		item.Nickname = user.FullName

		// 商家汇率
		rate, _ := merchantRateService.GetMerchantAllRate(item.Agency, item.MerchantName)
		item.Rate = rate
	}

	app.SuccessResp(c, response.RespMerchantInfoList{
		Total: total,
		List:  merchants,
	})
}

// @Summary 卖出
// @Description 卖出
// @Tags 商家
// @accept json
// @Produce json
// @Param ReqSellForm body request.ReqSellForm true "request.ReqSellForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/sell [post]
func Sell(c *gin.Context) {
	var (
		form request.ReqSellForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	//获取用户信息
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

	// 密码校验
	if user.TradeKey == "" || !util.MD5Equals(form.TradeKey, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}

	// 1、检查卖出卡信息是否存在
	cardService := merchant_card_service.MerchantCard{
		Agency:      account.Agency,
		Merchant:    account.Username,
		CardNo:      form.CardNo,
		CardType:    form.CardType,
		CardAccount: form.CardAccount,
		CardBank:    form.CardBank,
		CardSubBank: form.CardSubBank,
	}
	exist, err := cardService.Exist(account.Agency, account.Username, form.CardNo, form.CardType)
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_MERCHANT_CARD_FAIL, err.Error())
		return
	}

	// 2.如果卡不存在则添加到卡包
	if !exist {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logf.Error(err)
				}
			}()
			err := cardService.Add()
			if err != nil {
				logf.Error(err.Error())
			}
		}()
	}

	// 检查商家资金余额
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

	if fund.AvailableAmount < form.Amount {
		app.ErrorResp(c, e.Insufficient, "")
		return
	}

	// 3.创建卖单
	orderService := order_service.Order{}
	sellOrder, err := orderService.CreateSell(account.Agency, account.Username, account.Username, form.CardType, c.ClientIP(), form.Amount)
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_ORDER_FAIL, err.Error())
		return
	}

	// 4.冻结商家账户金额,修改订单状态
	err = orderService.MatchSell(account.Agency, sellOrder.OrderNo, sellOrder.Agency, form.CardNo, form.CardAccount, form.CardBank, form.CardSubBank)
	if err != nil {
		// 4.1、订单失败
		err2 := orderService.UpdateOrderFailed(sellOrder.Agency, sellOrder.OrderNo)
		if err2 != nil {
			logf.Error(err2)
		}
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	// 4、插入消息队列，通知agency转账

	app.SuccessResp(c, nil)
}

// @Summary 查看私钥
// @Description 查看私钥
// @Tags 商家
// @accept json
// @Produce json
// @Param ReqGetKeyForm body request.ReqGetKeyForm true "request.ReqGetKeyForm"
// @Success 200 {object}  response.RespGetKey
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/getKey [post]
func GetKey(c *gin.Context) {
	var (
		form request.ReqGetKeyForm
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

	// 密码校验
	if user.TradeKey == "" || !util.MD5Equals(form.TradeKey, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}

	merchantService := merchant_service.Merchant{}
	merchant, err := merchantService.Get(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespGetKey{
		PublicKey:  merchant.SystemPublicKey,
		PrivateKey: merchant.MerchantPrivateKey,
	})
}

// @Summary 代理查看商家秘钥
// @Description 代理查看商家秘钥
// @Tags 代理平台
// @accept json
// @Produce json
// @Param ReqAgencyGetMerchantKeyForm body request.ReqAgencyGetMerchantKeyForm true "request.ReqAgencyGetMerchantKeyForm"
// @Success 200 {object}  response.RespGetKey
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getMerchantKey [post]
func GetMerchantKey(c *gin.Context) {
	var (
		form request.ReqAgencyGetMerchantKeyForm
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

	// 密码校验
	if user.TradeKey == "" || !util.MD5Equals(form.TradeKey, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}

	merchantService := merchant_service.Merchant{}
	merchant, err := merchantService.Get(account.Agency, form.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespGetKey{
		PublicKey:  merchant.SystemPublicKey,
		PrivateKey: merchant.MerchantPrivateKey,
	})
}

// @Summary 添加卡信息
// @Description 添加卡信息
// @Tags 商家
// @accept json
// @Produce json
// @Param ReqAddMerchantCard body request.ReqAddMerchantCard true "request.ReqAddMerchantCard"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/addMerchantCard [post]
func AddMerchantCard(c *gin.Context) {
	var (
		form request.ReqAddMerchantCard
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	cardService := merchant_card_service.MerchantCard{
		Agency:      account.Agency,
		Merchant:    account.Username,
		CardType:    form.CardType,
		CardNo:      form.CardNo,
		CardAccount: form.CardAccount,
		CardBank:    form.CardBank,
		CardSubBank: form.CardSubBank,
		CardImg:     form.CardImg,
	}

	err = cardService.Add()
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_MERCHANT_CARD_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 删除卡信息
// @Description 删除卡信息
// @Tags 商家
// @accept json
// @Produce json
// @Param ReqRemoveMerchantCard body request.ReqRemoveMerchantCard true "request.ReqRemoveMerchantCard"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/removeMerchantCard [post]
func RemoveMerchantCard(c *gin.Context) {
	var (
		form request.ReqRemoveMerchantCard
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	cardService := merchant_card_service.MerchantCard{}
	err = cardService.Delete(form.CardId, account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_DELETE_MERCHANT_CARD_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 获取商家卡信息列表
// @Description 获取商家卡信息列表
// @Tags 商家
// @accept json
// @Produce  json
// @Param start_page query int false "起始页"
// @Param page_size query int false "页面大小"
// @Success 200 {object}  response.RespMerchantCardInfoList "response.RespMerchantCardInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/getMerchantCardList [get]
func GetMerchantCardList(c *gin.Context) {
	var (
		form request.ReqGetMerchantCardInfoListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	cardService := merchant_card_service.MerchantCard{
		Agency:   account.Agency,
		Merchant: account.Username,
		PageNum:  offset,
		PageSize: limit,
	}

	total, err := cardService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_MERCHANT_CARD_FAIL, err.Error())
		return
	}

	merchants, err := cardService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_MERCHANT_CARDS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespMerchantCardInfoList{
		Total: total,
		List:  merchants,
	})
}

// @Summary 修改返回地址
// @Description 修改返回地址
// @Tags 商家
// @accept json
// @Produce json
// @Param ReqUpdateReturnUrl body request.ReqUpdateReturnUrl true "request.ReqUpdateReturnUrl"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/updateReturnUrl [post]
func UpdateReturnUrl(c *gin.Context) {
	var (
		form request.ReqUpdateReturnUrl
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	merchantService := merchant_service.Merchant{}
	err = merchantService.UpdateReturnUrl(account.Agency, account.Username, form.ReturnUrl)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 修改回调地址
// @Description 修改回调地址
// @Tags 商家
// @accept json
// @Produce json
// @Param ReqUpdateNotifyUrl body request.ReqUpdateNotifyUrl true "request.ReqUpdateNotifyUrl"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/updateNotifyUrl [post]
func UpdateNotifyUrl(c *gin.Context) {
	var (
		form request.ReqUpdateNotifyUrl
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	merchantService := merchant_service.Merchant{}
	err = merchantService.UpdateNotifyUrl(account.Agency, account.Username, form.NotifyUrl)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 获取对接回调，返回地址
// @Description 获取对接回调，返回地址
// @Tags 商家
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespApiCallUrl "response.RespApiCallUrl"
// @Failure 500 {object}  app.Response
// @Router /api/v1/merchant/getApiCallUrl [get]
func GetApiCallUrl(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	merchantService := merchant_service.Merchant{}
	merchant, err := merchantService.Get(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_MERCHANT_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespApiCallUrl{
		ReturnUrl: merchant.ReturnUrl,
		NotifyUrl: merchant.NotifyUrl,
	})
}
