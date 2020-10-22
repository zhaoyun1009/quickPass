package constant

//收支类型
const (
	_ = iota
	//支出
	EXPEND
	//收入
	INCOME
)

const (
	//冻结
	FROZEN = iota + 1
	//解冻
	UNFROZEN
	//放行
	DISCHARGE
)

//会计科目
const (
	_ = iota
	//1.转账
	BusinessTypeTransfer
	//2.买入手续费
	BusinessTypeBuyFee
	//3.买入
	BusinessTypeBuy
	//4.卖出
	BusinessTypeSell
)
