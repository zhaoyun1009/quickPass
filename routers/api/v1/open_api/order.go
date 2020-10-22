package open_api

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/mq"
	"QuickPass/pkg/rsa"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/channel_service"
	"QuickPass/service/match_cache_service"
	"QuickPass/service/merchant_service"
	"QuickPass/service/order_service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Buy(c *gin.Context) {
	var (
		form request.ReqOpenApiBuyForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	data := form.SignData
	merchantService := merchant_service.Merchant{}
	merchant, err := merchantService.Get(data.Agency, data.Merchant)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_MERCHANT_FAIL, err.Error())
		return
	}
	if merchant == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_MERCHANT, "")
		return
	}

	// 判断商家可买入状态
	if merchant.BuyStatus == constant.SwitchClose {
		app.ErrorResp(c, e.MerchantBuyStatusClosed, "")
		return
	}

	jsonData, _ := json.Marshal(data)
	// 过滤特殊字符
	transHtmlJson := util.TransHtmlJson(jsonData)

	// 商家公钥验证签名的合法性
	err = rsa.RSAVerify(transHtmlJson, form.Sign, merchant.MerchantPublicKey)
	if err != nil {
		app.ErrorResp(c, e.CheckRsaError, "")
		return
	}

	// 2.判断支付通道
	channelService := channel_service.Channel{}
	channel := channelService.Get(data.Agency, data.PayType)
	if channel == nil || channel.IfOpen == constant.SwitchClose {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_CHANNEL, "")
		return
	}
	if data.Amount < channel.BuyMin || data.Amount > channel.BuyMax {
		app.ErrorResp(c, e.ERROR, "下单金额有误")
		return
	}

	// 6、创建订单
	orderService := order_service.Order{}

	returnUrl := data.ReturnUrl
	if returnUrl == "" {
		returnUrl = merchant.ReturnUrl
	}
	if returnUrl == "" {
		app.ErrorResp(c, e.ERROR, "return_url为空,请填写或前往商家平台设置")
		return
	}
	notifyUrl := data.NotifyUrl
	if notifyUrl == "" {
		notifyUrl = merchant.NotifyUrl
	}
	if notifyUrl == "" {
		app.ErrorResp(c, e.ERROR, "notify_url为空,请填写或前往商家平台设置")
		return
	}

	order, err := orderService.CreateBuy(data.Agency,
		data.Member,
		data.Merchant,
		data.PayType,
		c.ClientIP(),
		data.Amount,
		constant.SubmitTypeInterface,
		data.MerchantOrderNo,
		returnUrl,
		notifyUrl,
		data.AppendInfo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_ORDER_FAIL, err.Error())
		return
	}

	// 7、查询出可匹配的卡id
	context := match_cache_service.NewMatchContext("")
	card, err := context.Match(order.Agency, order.OrderNo, data.Amount, data.PayType)
	if err != nil {
		// 7.1、订单失败
		err2 := orderService.UpdateOrderFailed(order.Agency, order.OrderNo)
		if err2 != nil {
			logf.Error(err2)
		}
		app.ErrorResp(c, e.ERROR_MATCH_FAILED, err.Error())
		return
	}

	// 8、插入消息到队列中，通知承兑人
	order, _ = orderService.GetByOrderNo(order.Agency, order.OrderNo)
	// 发送订单状态改变通知
	go mq.SendOrderStatusChange(order)

	//订单超时时间
	duration := time.Minute * time.Duration(setting.AppSetting.OrderTimeoutMinute)

	if data.Model == constant.ModelApi {
		app.SuccessResp(c, response.RespOrderBuy{
			OrderNo:     order.OrderNo,
			CardNo:      card.CardNo,
			CardAccount: card.CardAccount,
			CardBank:    card.CardBank,
			CardImg:     card.CardImg,
			CardSubBank: card.CardSubBank,
			CreateTime:  order.CreateTime,
			ExpirationTime: util.JSONTime{
				Time: order.CreateTime.Add(duration),
			},
		})
		return
	}
	if data.Model == constant.ModelPage {
		pageUrl := fmt.Sprintf("%s/#/payment/index?agency=%s&orderNo=%s", setting.AppSetting.OrderOpenApiBuyUrl, order.Agency, order.OrderNo)
		app.SuccessPureResp(c, RespOpenApiOrderBuyModelPage{
			PageUrl: pageUrl,
		})
	}
}

type RespOpenApiOrderBuyModelPage struct {
	PageUrl string `json:"page_url"`
}

func GetStatusByMerchantOrderNo(c *gin.Context) {
	var (
		form request.ReqGetStatusByMerchantOrderNoForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	merchantOrderService := order_service.MerchantOrder{}
	merchantOrder, err := merchantOrderService.GetByMerchantOrderNo(form.Agency, form.MerchantUsername, form.MerchantOrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDER_FAIL, err.Error())
		return
	}
	if merchantOrder == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ORDER, "")
		return
	}

	orderService := order_service.Order{}
	order, err := orderService.GetByOrderNo(form.Agency, merchantOrder.OrderNo)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_ORDER_FAIL, err.Error())
		return
	}
	if order == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_ORDER, "")
		return
	}

	// 不是银行卡的需要查询收款码
	cardImg := ""
	if order.ChannelType != constant.BANK_CARD {
		cardImg = fmt.Sprintf("http://%s/%s", setting.MinioSetting.PreUrl, order.CardImg)
	}
	now := util.JSONTimeNow()
	//订单超时时间
	duration := time.Minute * time.Duration(setting.AppSetting.OrderTimeoutMinute)

	app.SuccessResp(c, response.RespGetOrderStatus{
		Agency:      order.Agency,
		OrderNo:     order.OrderNo,
		OrderType:   order.OrderType,
		OrderStatus: order.OrderStatus,
		Amount:      order.Amount,
		ChannelType: order.ChannelType,
		FinishTime:  order.FinishTime,
		CardNo:      order.CardNo,
		CardAccount: order.CardAccount,
		CardBank:    order.CardBank,
		CardImg:     cardImg,
		CardSubBank: order.CardSubBank,
		ReturnUrl:   merchantOrder.ReturnUrl,
		CreateTime:  order.CreateTime,
		CurrentTime: &now,
		ExpirationTime: util.JSONTime{
			Time: order.CreateTime.Add(duration),
		},
	})
}
