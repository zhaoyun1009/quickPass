package constant

//websocket消息通道消息类型
const (
	// 1订单状态改变
	MessageTypeOrderStatus = iota + 1
	// 2异常单状态改变
	MessageTypeAbnormalOrderStatus
)
