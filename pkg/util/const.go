package util

const (
	TIME_TEMPLATE_YEAR   = "2006"
	TIME_TEMPLATE_MOHTH  = "01"
	TIME_TEMPLATE_DAY    = "02"
	TIME_TEMPLATE_HOUR   = "15"
	TIME_TEMPLATE_MINUTE = "04"
	TIME_TEMPLATE_SECOND = "05"
	TIME_TEMPLATE_1      = "2006-01-02 15:04:05"
	TIME_TEMPLATE_2      = "2006/01/02 15:04:05"
	TIME_TEMPLATE_3      = "2006-01-02"
	TIME_TEMPLATE_4      = "15:04:05"
	TIME_TEMPLATE_5      = "20060102"

	TIME_TEMPLATE_6 = "00:00:00"
	TIME_TEMPLATE_7 = "23:59:59"
)

//常量
const (
	//gin的上下文的token_key
	TokenKey     = "claims"
	LogResponse  = "log_response"
	OperationLog = "operation_log"
	HeaderToken  = "token"
	// redis消息队列
	RedisMessageQueueKey = "FujiaMessageQueue"
	// redis频道：开奖
	RedisChannelLottery = "Fujia_Lottery"
	// redis频道：开奖
	RedisChannelConsultation = "Fujia_Consultation"
)

// 资金账户类型
const (
	// 普通资金账户
	AccountTypeOrdinary = 2
	// 系统账户
	AccountTypeSystem = 1
)

// 提现最小金额 100*10000
const MinCashAmount = 1000000

var UserRoleName = map[int]string{
	1: "普通用户",
	2: "系统账户",
}

const DefaultPassword = "12345678"
const DefaultTradeKey = "123456"
const GameRoomCodeString = "E5FCDG3HQA4B1NOPIJ2RSTUV67MWX89KLYZ"
const GacBase = 10000
const DefaultRoomName = "这是一个很帅的房间"
const DefaultRoomInfo = "这家伙很懒，没有任何介绍"
const DefaultRoomImg = "http://img02file.tooopen.com/images/20160408/tooopen_sy_158723161481.jpg"
const DefaultHeadImg = "http://img02file.tooopen.com/images/20160408/tooopen_sy_158723161481.jpg"
const MAX_ROOM_DEPOSIT = 1000000 * GacBase
const MIN_ROOM_DEPOSIT = 0 * GacBase
