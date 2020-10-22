package constant

// 商家订单回调状态(回调结果（1、未处理，2、回调成功，3、回调失败)
const (
	//1.未处理
	CallbackUnprocessed = iota + 1
	// 2、回调成功
	CallbackSuccess
	// 3、回调失败
	CallbackFailed
)
