package constant

// 用户角色
//ORDINARY普通用户
const (
	_ = iota
	//1:代理
	AGENCY
	//2：承兑人
	ACCEPTOR
	//3：商家
	MERCHANT
	//4：系统用户
	SYSTEM_USEER
	//5: 客服
	RECEPTONIST
)

// 账户类型
const (
	// 1:系统账户
	SystemAccount = iota + 1
	// 2:普通账户
	OrdinaryAccount
)

// 账户状态
const (
	_ = iota
	//1:暂停
	UNENABLE
	//2:启用
	ENABLE
)

//超级管理员用户名
var SuperAdministrator = "admin"
