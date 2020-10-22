package request

type ReqGetAgencyStatisticsForm struct {
	// 开始时间
	StartTime string `form:"start_time" json:"start_time" binding:"timeValidated"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time" binding:"timeValidated"`
}

type ReqGetUserBillStatisticsForm struct {
	// 用户名
	Username string `form:"username" json:"username" binding:"required,len=6"`

	// 开始日期
	StartDate string `form:"start_date" json:"start_date" binding:"dateValidated"`
	// 结束日期
	EndDate string `form:"end_date" json:"end_date" binding:"dateValidated"`

	StartPage int `form:"start_page" json:"start_page"`
	PageSize  int `form:"page_size" json:"page_size"`
}

type ReqUserByRoleBillStatisticsForm struct {
	// 用户名
	Username string `form:"username" json:"username" binding:"omitempty,len=6"`
	// 角色(2：承兑人、3：商家)
	Role int `form:"role" json:"role" binding:"required,oneof=2 3"`
	// 开始日期
	StartDate string `form:"start_date" json:"start_date" binding:"dateValidated"`
	// 结束日期
	EndDate string `form:"end_date" json:"end_date" binding:"dateValidated"`
}
