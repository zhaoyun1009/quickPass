package crontab

import (
	"QuickPass/pkg/constant"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/token_cache"
	"QuickPass/pkg/util"
	"QuickPass/service/acceptor_card_service"
	"QuickPass/service/acceptor_service"
	"QuickPass/service/bill_service"
	"QuickPass/service/fund_service"
	"QuickPass/service/order_service"
	"github.com/robfig/cron"
)

func Setup() {
	c := cron.New()
	// 每天午夜跑一次 等式(0 0 0 * * *)0 0/5 * * * ? 更新每天卡可用余额缓存
	c.AddFunc("@daily", updateAcceptorCardDayAvailableAmt)
	// 每天午夜跑一次,统计账单数据
	c.AddFunc("@daily", billStatistics)
	// 每1秒执行一次 检测订单超时
	c.AddFunc("@every 1s", updateBuyOrderTimeoutCancel)
	// 每1秒执行一次 商家回调
	c.AddFunc("@every 1s", orderCallback)
	// 已付款卖单30分钟自动放行
	c.AddFunc("@every 10s", updateSellOrderAutoDischarge)
	// 每5秒执行一次 检测承兑人下线关闭承兑状态
	c.AddFunc("@every 5s", checkAcceptorOffLine)
	// 启动服务
	c.Start()
}

// 每天午夜跑一次
// 更新承兑人卡当日剩余可用流水
func updateAcceptorCardDayAvailableAmt() {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	acceptorCard := acceptor_card_service.AcceptorCard{}
	err := acceptorCard.UpdateAcceptorCardDayAvailableAmt()
	if err != nil {
		logf.Error("updateAcceptorCardDayAvailableAmt: ", err.Error())
	}
}

// 更新订单表中的超时取消订单
func updateBuyOrderTimeoutCancel() {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	// 查询表中过期的订单
	orderService := order_service.Order{}
	timeoutList, err := orderService.GetOrderTimeoutList(setting.AppSetting.OrderTimeoutMinute)
	if err != nil {
		logf.Error(err, "GetOrderTimeoutList")
		return
	}

	for _, order := range timeoutList {
		// 发送订单状态改变通知
		err = orderService.CancelBuyOrder(order.Agency, order.OrderNo)
		if err != nil {
			logf.Error(err, "CancelBuyOrder", order.Agency, order.OrderNo)
			continue
		}
	}
}

// 订单回调
func orderCallback() {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	// 查询表中需要回调的订单
	orderService := order_service.Order{}
	callbackList, err := orderService.GetOrderCallbackList()
	if err != nil {
		logf.Error(err, "GetOrderCallbackList")
		return
	}

	for _, order := range callbackList {
		go orderService.CallbackMerchant(order)
	}
}

// 检测承兑人下线
func checkAcceptorOffLine() {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	offset, limit := util.GetPaginationParams(1, setting.AppSetting.PageSizeLimit)
	acceptorService := acceptor_service.Acceptor{
		PageNum:      offset,
		PageSize:     limit,
		AcceptSwitch: constant.SwitchOpen,
	}

	list, err := acceptorService.GetAll()
	if err != nil {
		logf.Error(err)
		return
	}

	for _, item := range list {
		exist, err := token_cache.CheckToken(item.Agency, item.Acceptor)
		if err != nil {
			logf.Error("token_cache.CheckToken:", err)
			continue
		}
		if !exist {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logf.Error(err)
					}
				}()
				_ = acceptorService.UpdateAcceptSwitch(item.Agency, item.Acceptor, constant.SwitchClose)
			}()
		}
	}
}

//统计账单数据
func billStatistics() {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	fundService := fund_service.Fund{}
	list, err := fundService.GetMerchantAndAcceptor()
	if err != nil {
		logf.Error(err.Error())
		return
	}

	yesterdayDateStr := util.YesterdayDateStr()

	startTimeStr := util.YesterdayStartTimeStr()
	endTimeStr := util.TodayStartTimeStr()

	for _, fund := range list {
		item := fund
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logf.Error(err)
				}
			}()
			var incomeExpense int
			switch item.Role {
			case constant.ACCEPTOR:
				incomeExpense = constant.EXPEND
			case constant.MERCHANT:
				incomeExpense = constant.INCOME
			}

			billService := bill_service.Bill{}
			// 承兑统计(账单类型：承兑人是支出，商家收入  账单科目：买入)
			acceptStatistics, err := billService.CronBillStatistics(item.Agency, item.UserName, startTimeStr, endTimeStr, incomeExpense, constant.BusinessTypeBuy)
			if err != nil {
				logf.Error(err.Error())
				return
			}

			// 提现统计
			withdrawalStatistics, err := billService.CronBillStatistics(item.Agency, item.UserName, startTimeStr, endTimeStr, constant.EXPEND, constant.BusinessTypeSell)
			if err != nil {
				logf.Error(err.Error())
				return
			}

			// 充值统计
			rechargeStatistics, err := billService.CronBillStatistics(item.Agency, item.UserName, startTimeStr, endTimeStr, constant.INCOME, constant.BusinessTypeTransfer)
			if err != nil {
				logf.Error(err.Error())
				return
			}

			billStatisticsService := bill_service.BillStatistics{
				Agency:         item.Agency,
				Username:       item.UserName,
				Role:           item.Role,
				StatisticsDate: yesterdayDateStr,
				LeftAmount:     item.AvailableAmount + item.FrozenAmount,
				// 承兑金额
				AcceptAmount: acceptStatistics.Amount,
				// 承兑次数
				AcceptCount: acceptStatistics.Count,
				// 提现金额
				WithdrawalAmount: withdrawalStatistics.Amount,
				// 提现次数
				WithdrawalCount: withdrawalStatistics.Count,
				// 充值金额
				RechargeAmount: rechargeStatistics.Amount,
				// 提现次数
				RechargeCount: rechargeStatistics.Count,
			}
			err = billStatisticsService.Add()
			if err != nil {
				logf.Error(err.Error())
				return
			}
		}()
	}
}

// 已付款卖单30分钟自动放行
func updateSellOrderAutoDischarge() {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	// 查询表中自动放行的提现订单
	orderService := order_service.Order{}
	autoDischargeList, err := orderService.GetSellOrderAutoDischargeList(setting.AppSetting.SellOrderAutoDischargeMinute)
	if err != nil {
		logf.Error(err, "GetOrderTimeoutList")
		return
	}
	for _, order := range autoDischargeList {
		err := orderService.DischargeSellOrder(order.Agency, order.OrderNo)
		if err != nil {
			logf.Error(err, "AutoDischargeSellOrder", order.Agency, order.OrderNo)
			continue
		}
	}
}
