package channel

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/channel_service"
	"QuickPass/service/merchant_rate_service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取代理的通道列表
// @Description 获取代理的通道列表
// @Tags 用户
// @accept json
// @Produce  json
// @Param agency query string true "代理"
// @Success 200 {object}  response.RespChannelInfoList "response.RespChannelInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/channels [get]
func GetChannels(c *gin.Context) {
	var (
		form request.ReqGetChannelsForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	offset, limit := util.GetPaginationParams(1, setting.AppSetting.PageSizeLimit)
	channelService := channel_service.Channel{
		PageNum:  offset,
		PageSize: limit,
		Agency:   form.Agency,
	}

	channels, err := channelService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_CHANNELS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespChannelInfoList{
		Total: int64(len(channels)),
		List:  channels,
	})
}

// @Summary 获取代理的通道信息列表
// @Description 获取代理的通道信息列表
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param start_page query int true "起始页"
// @Param page_size query int true "页面大小"
// @Success 200 {object}  response.RespChannelInfoList "response.RespChannelInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/getChannelInfoList [get]
func GetChannelInfoList(c *gin.Context) {
	var (
		form request.ReqGetChannelInfoListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	channelService := channel_service.Channel{
		PageNum:  offset,
		PageSize: limit,
		Agency:   account.Agency,
	}

	total, err := channelService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_CHANNEL_FAIL, err.Error())
		return
	}

	channels, err := channelService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_CHANNELS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespChannelInfoList{
		Total: total,
		List:  channels,
	})
}

// @Summary 更新通道汇率
// @Description 更新通道汇率
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqUpdateMerchantSingleRate body request.ReqUpdateMerchantSingleRate true "request.ReqUpdateMerchantSingleRate"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateMerchantSingleRate [post]
func UpdateMerchantSingleRate(c *gin.Context) {
	var (
		form request.ReqUpdateMerchantSingleRate
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	merchantRateService := merchant_rate_service.MerchantRate{}
	err = merchantRateService.UpdateMerchantSingleRate(account.Agency, form.MerchantName, form.Channel, form.Rate)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_MERCHANT_RATE_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 统一更新通道汇率
// @Description 统一更新通道汇率
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqUpdateMerchantAllChannelRate body request.ReqUpdateMerchantAllChannelRate true "request.ReqUpdateMerchantAllChannelRate"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateMerchantAllChannelRate [post]
func UpdateMerchantAllChannelRate(c *gin.Context) {
	var (
		form request.ReqUpdateMerchantAllChannelRate
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	merchantRateService := merchant_rate_service.MerchantRate{}
	err = merchantRateService.UpdateMerchantAllChannelRate(account.Agency, form.MerchantName, form.Rate)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_MERCHANT_RATE_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 更新每日承兑上限
// @Description 更新每日承兑上限
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqUpdateLimitAmountForm body request.ReqUpdateLimitAmountForm true "id:通道ID limit_amount:每日承兑上限金额 "
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateLimitAmount [post]
func UpdateLimitAmount(c *gin.Context) {
	var (
		form request.ReqUpdateLimitAmountForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	channelService := channel_service.Channel{
		Id:          form.Id,
		Agency:      account.Agency,
		LimitAmount: form.LimitAmount,
	}

	err = channelService.UpdateLimitAmount()
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_CHANNEL_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 统一更新每日承兑上限
// @Description 统一更新每日承兑上限
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqUpdateAllLimitAmountForm body request.ReqUpdateAllLimitAmountForm true "limit_amount:每日承兑上限金额 "
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/updateAllLimitAmount [post]
func UpdateAllLimitAmount(c *gin.Context) {
	var (
		form request.ReqUpdateAllLimitAmountForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	channelService := channel_service.Channel{
		Agency:      account.Agency,
		LimitAmount: form.LimitAmount,
	}

	err = channelService.UpdateAllLimitAmount()
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_CHANNEL_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

func UpdateBuyMaxAmount(c *gin.Context) {
	var (
		form request.ReqUpdateBuyMaxAmountForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	channelService := channel_service.Channel{
		Id:     form.Id,
		Agency: account.Agency,
		BuyMax: form.Amount,
	}

	err = channelService.UpdateBuyMaxAmount()
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_CHANNEL_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

func UpdateBuyMinAmount(c *gin.Context) {
	var (
		form request.ReqUpdateBuyMinAmountForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	channelService := channel_service.Channel{
		Id:     form.Id,
		Agency: account.Agency,
		BuyMin: form.Amount,
	}

	err = channelService.UpdateBuyMinAmount()
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_CHANNEL_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 通道开关
// @Description 通道开关
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqChannelSwitchForm body request.ReqChannelSwitchForm true "channel_id:通道ID if_open:通道开关"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/channelSwitch [post]
func ChannelSwitch(c *gin.Context) {
	var (
		form request.ReqChannelSwitchForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 获取用户信息
	account := c.MustGet(util.TokenKey).(*util.Claims)

	channelService := channel_service.Channel{
		Id:     form.ChannelId,
		IfOpen: form.IfOpen,
		Agency: account.Agency,
	}

	err = channelService.UpdateIfOpen()
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_CHANNEL_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}
