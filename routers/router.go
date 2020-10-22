package routers

import (
	"QuickPass/middleware/jwt"
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/util"
	"QuickPass/routers/api"
	"QuickPass/routers/api/v1/abnormal_order"
	"QuickPass/routers/api/v1/agency"
	"QuickPass/routers/api/v1/channel"
	"QuickPass/routers/api/v1/fund"
	"QuickPass/routers/api/v1/management"
	"QuickPass/routers/api/v1/open_api"
	"QuickPass/routers/api/v1/order"
	"QuickPass/routers/api/v1/websocket"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"QuickPass/routers/api/v1/acceptor"
	"QuickPass/routers/api/v1/merchant"
	"QuickPass/routers/api/v1/user"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	//日志中间件,所有的异常捕获,logrus
	r.Use(cors(), gin.Recovery(), initLog())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/ws", websocket.TokenWebSocketHandler)

	// 文件上传
	uploadGroup := apiV1.Group("/upload")
	uploadGroup.Use(jwt.JWT())
	{
		uploadGroup.POST("/image", api.UploadImage)
	}

	// 登录
	apiV1.POST("/login", user.Login)
	apiV1.POST("/backend/management/login", management.Login)

	// 获取可用的代理列表
	apiV1.GET("/agencyList", user.GetAgencyList)
	// 获取代理的通道列表
	apiV1.GET("/agency/channels", channel.GetChannels)

	//管理者
	apiV1Management := apiV1.Group("/backend/management")
	apiV1Management.Use(jwt.JWT(), jwt.RequiredRole(constant.SYSTEM_USEER))
	{
		//获取后台管理员信息列表
		//apiV1Management.GET("/getManagementInfoList", management.GetManagementInfoList)
		//后台放行
		//apiV1Management.POST("/backEndDischarge", order.BackEndDischarge)
		//后台放行卖单
		//apiV1Management.POST("/backendDischargeSell", order.BackendDischargeSell)
		//添加代理用户
		apiV1Management.POST("/backendAddAgency", user.AddAgencyUser)
		//转账
		apiV1Management.POST("/transfer", fund.SysTransfer)
		//获取用户信息列表
		apiV1Management.GET("/getUserInfoList", user.GetUserInfoList)
	}

	//代理平台
	apiV1Agency := apiV1.Group("/agency")
	apiV1Agency.Use(jwt.JWT(), jwt.RequiredRole(constant.AGENCY))
	{
		//代理平台统计(已有商家，已有承兑，承兑总数额，手续费收入)
		apiV1Agency.GET("/getAgencyStatistics", agency.GetAgencyStatistics)

		//重置代理下的用户登录密码
		apiV1Agency.POST("/resetUserPassword", user.ResetAgencyUserPassword)
		//重置代理下的用户交易密码
		apiV1Agency.POST("/resetUserTradeKey", user.ResetAgencyUserTradeKey)
		//承兑管理统计
		apiV1Agency.GET("/getAcceptorStatistics", acceptor.GetAcceptorStatistics)
		//获取承兑人信息列表
		apiV1Agency.GET("/getAcceptorInfoList", acceptor.GetAcceptorInfoList)
		//添加承兑人
		apiV1Agency.POST("/addAcceptor", acceptor.AddAcceptor)
		//商家管理统计
		apiV1Agency.GET("/getMerchantStatistics", merchant.GetMerchantStatistics)
		//获取商家信息列表
		apiV1Agency.GET("/getMerchantInfoList", merchant.GetMerchantInfoList)
		//添加商家
		apiV1Agency.POST("/addMerchant", merchant.AddMerchant)
		//查看商家秘钥
		apiV1Agency.POST("/getMerchantKey", merchant.GetMerchantKey)
		//删除商家
		apiV1Agency.POST("/removeMerchant", merchant.RemoveMerchant)
		//商家数据统计和承兑人数据统计
		apiV1Agency.GET("/getUserByRoleBillStatistics", fund.GetUserByRoleBillStatistics)
		//获取用户每日账单统计
		apiV1Agency.GET("/getUserBillStatistics", fund.GetUserBillStatistics)

		//获取代理的通道信息列表
		apiV1Agency.GET("/getChannelInfoList", channel.GetChannelInfoList)
		//更新商家通道汇率
		apiV1Agency.POST("/updateMerchantSingleRate", channel.UpdateMerchantSingleRate)
		//统一更新商家通道汇率
		apiV1Agency.POST("/updateMerchantAllChannelRate", channel.UpdateMerchantAllChannelRate)
		//更新每日承兑上限
		apiV1Agency.POST("/updateLimitAmount", channel.UpdateLimitAmount)
		//统一每日承兑上限
		apiV1Agency.POST("/updateAllLimitAmount", channel.UpdateAllLimitAmount)
		// 更新最大买入金额
		apiV1Agency.POST("/updateBuyMaxAmount", channel.UpdateBuyMaxAmount)
		// 更新最小买入金额
		apiV1Agency.POST("/updateBuyMinAmount", channel.UpdateBuyMinAmount)
		//通道开关
		apiV1Agency.POST("/channelSwitch", channel.ChannelSwitch)

		//代理异常单列表
		apiV1Agency.GET("/getAbnormalOrderInfoList", abnormal_order.GetAgencyAbnormalOrderInfoList)
		//代理异常订单中所有承兑人的列表
		apiV1Agency.GET("/getAbnormalOrderAcceptorGroup", abnormal_order.GetAbnormalAcceptorGroup)
		//更新异常单生成订单
		apiV1Agency.POST("/updateAbnormalOrder", abnormal_order.UpdateAbnormalOrder)

		//代理订单列表
		apiV1Agency.GET("/getOrderInfoList", order.GetAgencyOrderInfoList)
		//代理订单中所有承兑人的列表
		apiV1Agency.GET("/getAcceptorGroup", order.GetAcceptorGroup)
		//代理放行订单
		apiV1Agency.POST("/agencyDischarge", order.AgencyDischarge)
		//代理放行超时取消的订单
		apiV1Agency.POST("/agencyDischargeCancelOrder", order.AgencyDischargeCancelOrder)
		//代理再次回调商家
		apiV1Agency.POST("/againCallback", order.AgencyAgainCallback)

		//代理卖单列表
		apiV1Agency.GET("/getSellOrderInfoList", merchant.GetAgencySellOrderInfoList)
		//代理卖单未处理数
		apiV1Agency.GET("/getSellOrderUnprocessed", order.GetSellOrderUnprocessed)
		//代理卖单确认支付
		apiV1Agency.POST("/confirmSellPay", order.ConfirmSellPay)
		//代理转账给承兑人
		apiV1Agency.POST("/transfer", fund.AgencyTransfer)
		// 代理最近转账信息
		apiV1Agency.GET("/latelyTransferInfo", fund.GetAgencyLatelyTransferInfo)

		// 代理获取转账给承兑人的账单
		apiV1Agency.GET("/getAcceptorBill", fund.GetAcceptorBill)
		//代理更新承兑权限
		apiV1Agency.POST("/updateAcceptSwitch", acceptor.AgencyUpdateAcceptSwitch)
		//代理更新承兑账户状态
		apiV1Agency.POST("/updateAcceptStatus", acceptor.AgencyUpdateAcceptStatus)
		//代理更新自动承兑
		apiV1Agency.POST("/updateIfAutoAccept", acceptor.AgencyUpdateIfAutoAccept)

		// 代理更新商家买入状态
		apiV1Agency.POST("/updateMerchantBuyStatus", merchant.AgencyUpdateMerchantBuyStatus)
	}

	//普通用户
	apiV1User := apiV1.Group("/user")
	apiV1User.Use(jwt.JWT())
	{
		//获取用户信息
		apiV1User.GET("/getUserInfo", user.GetUserInfo)
		//更新用户信息
		apiV1User.POST("/updateUserInfo", user.UpdateUserInfo)
		//修改登录密码
		apiV1User.POST("/modifyPassword", user.ModifyPassword)
		//修改交易密码
		apiV1User.POST("/modifyTradeKey", user.ModifyTradeKey)
		//是否存在交易密码
		apiV1User.GET("/existTradeKey", user.ExistTradeKey)
		//设置交易密码
		apiV1User.POST("/settingTradeKey", user.SettingTradeKey)
	}

	//承兑平台
	apiV1Accept := apiV1.Group("/acceptor")
	apiV1Accept.Use(jwt.JWT(), jwt.RequiredRole(constant.ACCEPTOR))
	{
		//获取承兑人信息
		apiV1Accept.GET("/getAcceptorInfo", acceptor.GetAcceptorInfo)
		//增加承兑人收款方式
		apiV1Accept.POST("/addAcceptorCard", acceptor.AddAcceptorCard)
		//更新承兑权限
		apiV1Accept.POST("/updateAcceptSwitch", acceptor.UpdateAcceptSwitch)
		//更新自动承兑
		apiV1Accept.POST("/updateIfAutoAccept", acceptor.UpdateIfAutoAccept)
		//更新承兑人卡信息
		apiV1Accept.POST("/updateAcceptorCard", acceptor.UpdateAcceptorCard)
		//更新承兑人卡状态
		apiV1Accept.POST("/updateCardStatus", acceptor.UpdateCardStatus)
		//删除承兑人卡
		apiV1Accept.POST("/deleteAcceptorCard", acceptor.DeleteAcceptorCard)
		//获取承兑人卡列表
		apiV1Accept.GET("/getAcceptorCardInfoList", acceptor.GetAcceptorCardInfoList)
		//获取当前承兑人所有卡信息
		apiV1Accept.GET("/getAcceptorAllCardInfo", acceptor.GetAcceptorAllCardInfo)
		//承兑人放行
		apiV1Accept.POST("/acceptorDischarge", acceptor.AcceptorDischarge)

		//订单查询
		apiV1Accept.GET("/getBuyOrderInfoList", order.GetBuyOrderInfoList)

		//异常单查询
		apiV1Accept.GET("/getAbnormalOrderInfoList", abnormal_order.GetAcceptorAbnormalOrderInfoList)
		//新建异常单
		apiV1Accept.POST("/addAbnormalOrder", abnormal_order.AddAcceptorAbnormalOrder)
	}

	//商家平台
	apiV1Merchant := apiV1.Group("/merchant")
	apiV1Merchant.Use(jwt.JWT(), jwt.RequiredRole(constant.MERCHANT))
	{
		//商家卖单
		apiV1Merchant.POST("/sell", merchant.Sell)
		//查看私钥
		apiV1Merchant.POST("/getKey", merchant.GetKey)
		//取消卖单
		apiV1Merchant.POST("/cancelSell", order.CancelSell)
		//买入订单查询
		apiV1Merchant.GET("/getBuyOrderInfoList", merchant.GetBuyOrderInfoList)
		//卖出订单查询
		apiV1Merchant.GET("/getSellOrderInfoList", merchant.GetSellOrderInfoList)
		//商家放行卖单
		apiV1Merchant.POST("/merchantDischargeSell", order.MerchantDischargeSell)
		// 添加卡信息
		apiV1Merchant.POST("/addMerchantCard", merchant.AddMerchantCard)
		// 删除卡信息
		apiV1Merchant.POST("/removeMerchantCard", merchant.RemoveMerchantCard)
		// 卡列表查询
		apiV1Merchant.GET("/getMerchantCardList", merchant.GetMerchantCardList)
		// 修改返回地址
		apiV1Merchant.POST("/updateReturnUrl", merchant.UpdateReturnUrl)
		// 修改回调地址
		apiV1Merchant.POST("/updateNotifyUrl", merchant.UpdateNotifyUrl)
		// 获取对接回调，返回地址
		apiV1Merchant.GET("/getApiCallUrl", merchant.GetApiCallUrl)
	}

	//资金
	apiV1Fund := apiV1.Group("/fund")
	apiV1Fund.Use(jwt.JWT())
	{
		//获取资金账户金额信息
		apiV1Fund.GET("/getFundInfo", fund.GetFundInfo)
		//获取收入/支出统计信息
		apiV1Fund.GET("/bill/inComeStatistics", fund.InComeStatistics)
		//账单
		apiV1Fund.GET("/bill/getUserBillList", fund.GetUserBillList)
	}

	//订单
	apiV1Order := apiV1.Group("/order")
	apiV1Order.Use()
	{
		//买单
		apiV1Order.POST("/buy", order.Buy)
		//订单状态查询
		apiV1Order.GET("/getStatusByOrderNo", order.GetStatusByOrderNo)
		//取消买单
		apiV1Order.POST("/cancelBuy", order.CancelBuy)
		//确认付款
		apiV1Order.POST("/confirmPay", order.ConfirmPay)
	}

	// 开放api（商家对接）
	apiV1OpenApi := apiV1.Group("/openApi")
	apiV1OpenApi.Use()
	{
		//买单
		apiV1OpenApi.POST("/buy", open_api.Buy)
		//商家订单号订单状态查询
		apiV1OpenApi.GET("/getStatusByMerchantOrderNo", open_api.GetStatusByMerchantOrderNo)
	}
	return r
}

//跨域中间件
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, "+util.HeaderToken)
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, "+util.HeaderToken)
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func initLog() gin.HandlerFunc {

	return func(c *gin.Context) {
		tranceId := c.Request.Header.Get("trance_id")
		if tranceId == "" {
			c.Request.Header.Set("trance_id", "uuid")
		} else {
			c.Request.Header.Set("trance_id", tranceId+"->uuid")
		}

		// 开始时间
		startTime := time.Now()
		// 请求路由
		path := c.Request.RequestURI

		// 排除文件上传的请求体打印
		isFormData := strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data")
		// requestBody
		var requestBody []byte
		if !isFormData {
			requestBody, _ = c.GetRawData()
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		//处理请求
		c.Next()
		// 处理结果
		result, exists := c.Get(util.LogResponse)
		if exists {
			result = result.(*app.Response)
		}

		// 执行时间
		latencyTime := time.Since(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// http状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		//token := c.GetHeader(tool.HeaderToken)
		// 日志格式
		logf.InfoWithFields(logrus.Fields{
			"trance_id":    tranceId,
			"req_body":     string(requestBody),
			"http_code":    statusCode,
			"latency_time": fmt.Sprintf("%13v", latencyTime),
			"ip":           clientIP,
			"method":       reqMethod,
			"path":         path,
			"result":       result,
		})
	}
}
