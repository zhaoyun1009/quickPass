package acceptor

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/acceptor_card_service"
	"QuickPass/service/acceptor_service"
	"QuickPass/service/order_service"
	"QuickPass/service/user_service"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @Summary 增加承兑人收款方式
// @Description 增加承兑人收款方式
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqAddAcceptorCardForm body request.ReqAddAcceptorCardForm true "request.ReqAddAcceptorCardForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/addAcceptorCard [post]
func AddAcceptorCard(c *gin.Context) {
	var (
		form request.ReqAddAcceptorCardForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 1.查询是否存在承兑人
	exist, err := existAcceptor(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_ACCEPTOR_FAIL, err.Error())
		return
	}
	if !exist {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ACCEPTOR, "")
		return
	}

	// 3.定义数据
	acceptorCardService := acceptor_card_service.AcceptorCard{
		Agency:      account.Agency,
		Acceptor:    account.Username,
		CardType:    form.CardType,
		CardNo:      form.CardNo,
		CardAccount: form.CardAccount,
		CardBank:    form.CardBank,
		CardSubBank: form.CardSubBank,
		CardImg:     form.CardImg,
	}

	// 4.插入数据库数据
	err = acceptorCardService.Add()
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 更新承兑人卡信息
// @Description 更新承兑人卡信息
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqUpdateAcceptorCardForm body request.ReqUpdateAcceptorCardForm true "card_id：卡id card_account：卡账户名 card_bank：银行卡名称 card_sub_bank:支行 card_img:图片地址 signal_max_amt:单笔最大金额 day_max_amt:日最大金额"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/updateAcceptorCard [post]
func UpdateAcceptorCard(c *gin.Context) {
	var (
		form request.ReqUpdateAcceptorCardForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	account := c.MustGet(util.TokenKey).(*util.Claims)

	acceptorCardService := acceptor_card_service.AcceptorCard{
		Id:          form.CardId,
		CardAccount: form.CardAccount,
		CardBank:    form.CardBank,
		CardSubBank: form.CardSubBank,
		CardImg:     form.CardImg,
	}

	// 2.检查用户卡信息
	card, err := acceptorCardService.Get()
	if err != nil {
		app.ErrorResp(c, e.ERROR_CHECK_EXIST_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}
	if card == nil || (account.Username != card.Acceptor) || (account.Agency != card.Agency) {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ACCEPTOR_CARD, "")
		return
	}

	// 4.更新数据
	err = acceptorCardService.UpdateCardInfo()
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 更新承兑人卡状态
// @Description 更新承兑人卡状态
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqUpdateCardStatusForm body request.ReqUpdateCardStatusForm true "card_id:卡id if_open:是否开启在线"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/updateCardStatus [post]
func UpdateCardStatus(c *gin.Context) {
	var (
		form request.ReqUpdateCardStatusForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 1.定义查询和更新参数
	acceptorCardService := acceptor_card_service.AcceptorCard{
		Id:     form.CardId,
		IfOpen: form.IfOpen,
	}

	// 2.获取卡信息
	card, err := acceptorCardService.Get()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}
	// 3.判断承兑人是否一致
	if (card.Agency != account.Agency) || (card.Acceptor != account.Username) {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ACCEPTOR_CARD, "")
		return
	}

	// 4.更新卡状态
	err = acceptorCardService.UpdateCardState()
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 删除承兑人卡
// @Description 删除承兑人卡
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqDeleteAcceptorCardForm body request.ReqDeleteAcceptorCardForm true "card_id：卡id"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/deleteAcceptorCard [post]
func DeleteAcceptorCard(c *gin.Context) {
	var (
		form request.ReqDeleteAcceptorCardForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 1.定义删除参数
	acceptorCardService := acceptor_card_service.AcceptorCard{
		Id: form.CardId,
	}

	// 2.获取卡信息
	card, err := acceptorCardService.Get()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}
	// 3.判断承兑人是否一致
	if card == nil || (card.Agency != account.Agency) || (card.Acceptor != account.Username) {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ACCEPTOR_CARD, "")
		return
	}

	// 4.删除卡信息
	err = acceptorCardService.Delete()
	if err != nil {
		app.ErrorResp(c, e.ERROR_DELETE_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 获取承兑人卡列表
// @Description 获取承兑人卡列表
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param card_type query string true "卡类型"
// @Param start_page query int true "起始页"
// @Param page_size query int true "页面大小"
// @Success 200 {object}  response.RespAcceptorCardInfoList "response.RespAcceptorCardInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/getAcceptorCardInfoList [get]
func GetAcceptorCardInfoList(c *gin.Context) {
	var (
		form request.ReqGetAcceptorCardInfoListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	cardService := acceptor_card_service.AcceptorCard{
		Agency:   account.Agency,
		Acceptor: account.Username,
		CardType: form.CardType,
		PageNum:  offset,
		PageSize: limit,
	}

	total, err := cardService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_ACCEPTOR_CARD_FAIL, err.Error())
		return
	}

	cards, err := cardService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ACCEPTOR_CARDS_FAIL, err.Error())
		return
	}

	for _, item := range cards {
		if item.CardType != constant.BANK_CARD {
			item.CardImg = fmt.Sprintf("http://%s/%s", setting.MinioSetting.PreUrl, item.CardImg)
		}
	}

	app.SuccessResp(c, response.RespAcceptorCardInfoList{
		Total: total,
		List:  cards,
	})
}

// @Summary 获取当前承兑人所有卡信息
// @Description 获取当前承兑人所有卡信息
// @Tags 承兑人
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespAcceptorAllCardInfo "response.RespAcceptorAllCardInfo"
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/getAcceptorAllCardInfo [get]
func GetAcceptorAllCardInfo(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	cardService := acceptor_card_service.AcceptorCard{}
	total, cards, err := cardService.GetAllCardInfo(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ACCEPTOR_CARDS_FAIL, err.Error())
		return
	}

	for _, item := range cards {
		if item.CardType != constant.BANK_CARD {
			item.CardImg = fmt.Sprintf("http://%s/%s", setting.MinioSetting.PreUrl, item.CardImg)
		}
	}

	app.SuccessResp(c, response.RespAcceptorAllCardInfo{
		Total: total,
		List:  cards,
	})
}

// @Summary 放行
// @Description 放行
// @Tags 承兑人
// @accept json
// @Produce  json
// @Param ReqAcceptorDischargeForm body request.ReqAcceptorDischargeForm true "request.AcceptorDischargeForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/acceptor/acceptorDischarge [post]
func AcceptorDischarge(c *gin.Context) {
	var (
		form request.ReqAcceptorDischargeForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	// 1、核对交易密码
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
	if user.TradeKey == "" || !util.MD5Equals(form.TradeKey, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}

	// 2、订单放行
	orderService := order_service.Order{}
	err = orderService.DischargeBuyOrder(account.Agency, form.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_EDIT_ORDER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// 判断承兑人是否存在
func existAcceptor(agency string, username string) (bool, error) {
	acceptorService := acceptor_service.Acceptor{
		Agency:   agency,
		Acceptor: username,
	}
	return acceptorService.ExistByAgencyAndUsername()
}
