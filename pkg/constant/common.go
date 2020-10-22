package constant

//自动承兑
const (
	// 手动
	Manual = iota + 1
	// 自动
	Auto
)

// 承兑开关
const (
	// 关闭
	SwitchClose = iota + 1
	// 开启
	SwitchOpen
)

//卡状态
const (
	// 未删除
	NotDeleted = iota + 1
	// 已删除
	Deleted
)
