package constant

// 订单类型
const (
	_ = iota
	// 转账
	ORDER_TYPE_TRANSFER
	// 买单
	ORDER_TYPE_BUY
	// 卖单
	ORDER_TYPE_SELL
)

// 订单状态
const (
	_ = iota
	// 订单创建
	ORDER_STATUS_CREATE
	// 订单等待支付
	ORDER_STATUS_WAIT_PAY
	// 等待放行
	ORDER_STATUS_WAIT_DISCHARGE
	// 订单取消
	ORDER_STATUS_CANCEL
	// 订单失败
	ORDER_STATUS_FAILED
	// 订单完成
	ORDER_STATUS_FINISHED
)

//提交类型
const (
	_ = iota
	//直接提交
	SubmitTypeDirect
	//接口提交
	SubmitTypeInterface
)

//接口模式
const (
	ModelPage = "page"
	ModelApi  = "api"
)

// 回调请求状态（1、未请求，2、请求中
const (
	_ = iota
	CallbackUnLock
	CallbackLock
)
