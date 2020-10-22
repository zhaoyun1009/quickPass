package constant

// 异常单状态
const (
	//1.未处理
	Unprocessed = iota + 1
	//2.已处理
	Processed
	//3.处理中
	Processing
)

// 异常单收款通道
const (
	// 1.支付宝转银行卡
	AliPayToBankCard = iota + 1
	// 2.微信转银行卡
	WeiXinToBankCard
	// 3.银行卡转银行卡
	BankCardToBankCard
)
