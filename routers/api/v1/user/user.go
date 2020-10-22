package user

import (
	"QuickPass/models"
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/user_service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取用户信息列表
// @Description 获取用户信息列表人
// @Tags 后台管理
// @accept json
// @Produce  json
// @Param agency query string true "代理 所属系统 默认admin"
// @Param start_page query int true "起始页"
// @Param page_size query int true "页面大小"
// @Success 200 {object}  response.RespUserInfoList "response.RespUserInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/backend/management/getUserInfoList [get]
func GetUserInfoList(c *gin.Context) {
	var (
		form request.ReqGetUserInfoListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	userService := user_service.User{
		Agency:   form.Agency,
		PageNum:  offset,
		PageSize: limit,
	}

	total, err := userService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_USER_FAIL, err.Error())
		return
	}

	users, err := userService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespUserInfoList{
		Total: total,
		List:  users,
	})
}

// @Summary 添加代理用户
// @Description 添加代理用户
// @Tags 后台管理
// @accept json
// @Produce  json
// @Param ReqAddAgencyUserForm body request.ReqAddAgencyUserForm true "request.ReqAddAgencyUserForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/backend/management/backendAddAgency [post]
func AddAgencyUser(c *gin.Context) {
	var (
		form request.ReqAddAgencyUserForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	userService := user_service.User{
		Agency:      form.UserName,
		UserName:    form.UserName,
		Password:    form.Password,
		FullName:    form.FullName,
		PhoneNumber: form.PhoneNumber,
		Address:     form.Address,
	}

	err = userService.AddAgency()
	if err != nil {
		app.ErrorResp(c, e.ERROR_ADD_USER_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, nil)
}

// @Summary 登录
// @Description 登录
// @Tags 用户
// @accept json
// @Produce  json
// @Param request body request.ReqUserLoginForm true "agency:代理 userName:用户名  password:密码"
// @Success 200 {object}  response.RespUserLogin
// @Failure 500 {object}  app.Response
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	var (
		form request.ReqUserLoginForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 1、密码验证 身份验证
	// 1.1、获取用户信息
	user, err := getUserByAgencyAndUsername(form.Agency, form.UserName)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	// 2、校验密码
	if user == nil || !util.MD5Equals(form.Password, user.Password) {
		app.ErrorResp(c, e.LOGIN_FAILED, "")
		return
	}

	// 2、token生成
	token, err := util.GenerateToken(user.Agency, user.UserName, user.FullName, int8(user.Role))
	if err != nil {
		app.ErrorResp(c, e.ERROR_AUTH_TOKEN, err.Error())
		return
	}

	app.SuccessResp(c, response.RespUserLogin{
		Token: token,
		Role:  user.Role,
	})
}

// @Summary 获取代理组
// @Description 获取代理组
// @Tags 用户
// @accept json
// @Produce  json
// @Success 200 {object}  []string
// @Failure 500 {object}  app.Response
// @Router /api/v1/agencyList [get]
func GetAgencyList(c *gin.Context) {
	userService := user_service.User{}
	list, err := userService.AgencyList()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.ReqAgencyList{
		List: list,
	})
}

// @Summary 修改登录密码
// @Description 修改登录密码
// @Tags 用户
// @accept json
// @Produce  json
// @Param ReqPasswordModifyForm body request.ReqPasswordModifyForm true "original_password:原始密码 update_password:更新密码"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/user/modifyPassword [post]
func ModifyPassword(c *gin.Context) {
	var (
		form request.ReqPasswordModifyForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 1、获取用户信息
	user, err := getUserByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	// 2、校验原始密码
	if !util.MD5Equals(form.OriginalPassword, user.Password) {
		app.ErrorResp(c, e.ERROR_VALIDATE_PASSWORD, "")
		return
	}

	// 3、修改密码
	updateUser := user_service.User{
		Id:       user.Id,
		Password: util.EncodeMD5(form.UpdatePassword),
	}

	err = updateUser.UpdatePassword()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_PASSWORD, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

// @Summary 修改交易密码
// @Description 修改交易密码
// @Tags 用户
// @accept json
// @Produce  json
// @Param ReqTradeKeyModifyForm body request.ReqTradeKeyModifyForm true "original_trade_password:原始交易密码 update_trade_password:更新交易密码"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/user/modifyTradeKey  [post]
func ModifyTradeKey(c *gin.Context) {
	var (
		form request.ReqTradeKeyModifyForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 1、获取用户信息
	user, err := getUserByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}
	// 2、校验交易密码
	if user.TradeKey == "" || !util.MD5Equals(form.OriginalTradePassword, user.TradeKey) {
		app.ErrorResp(c, e.ERROR_VALIDATE_TRADE_PASSWORD, "")
		return
	}
	// 3、修改交易密码
	updateUser := user_service.User{
		Id:       user.Id,
		TradeKey: util.EncodeMD5(form.UpdateTradePassword),
	}

	err = updateUser.UpdateTradeKey()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_TRADE_PASSWORD, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

// @Summary 是否存在交易密码
// @Description 是否存在交易密码
// @Tags 用户
// @accept json
// @Produce  json
// @Success 200 {object}  response.RespExistTradeKey
// @Failure 500 {object}  app.Response
// @Router /api/v1/user/existTradeKey  [get]
func ExistTradeKey(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	user, err := getUserByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	app.SuccessResp(c, response.RespExistTradeKey{
		Exist: user.TradeKey != "",
	})
}

// @Summary 设置交易密码
// @Description 设置交易密码
// @Tags 用户
// @accept json
// @Produce  json
// @Param ReqSettingTradeKeyForm body request.ReqSettingTradeKeyForm true "request.ReqSettingTradeKeyForm"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/user/settingTradeKey  [post]
func SettingTradeKey(c *gin.Context) {
	var (
		form request.ReqSettingTradeKeyForm
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)
	// 1、获取用户信息
	user, err := getUserByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	// 2、校验交易密码
	if user.TradeKey != "" {
		app.ErrorResp(c, e.ERROR_UPDATE_TRADE_PASSWORD, "")
		return
	}

	// 3、修改交易密码
	updateUser := user_service.User{
		Id:       user.Id,
		TradeKey: util.EncodeMD5(form.TradeKey),
	}

	err = updateUser.UpdateTradeKey()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_TRADE_PASSWORD, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 用户
// @accept json
// @Produce  json
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/user/getUserInfo  [get]
func GetUserInfo(c *gin.Context) {
	account := c.MustGet(util.TokenKey).(*util.Claims)

	user, err := getUserByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	app.SuccessResp(c, user)
}

// @Summary 更新用户信息
// @Description 更新用户信息
// @Tags 用户
// @accept json
// @Produce  json
// @Param ReqUserUpdateForm body request.ReqUserUpdateForm true "phone_number:电话号码 full_name:姓名 address:地址"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/user/updateUserInfo  [post]
func UpdateUserInfo(c *gin.Context) {
	var (
		form request.ReqUserUpdateForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	//1.获取用户信息
	user, err := getUserByAgencyAndUsername(account.Agency, account.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, "")
		return
	}

	//2.更新用户信息
	updateUser := user_service.User{
		Id:          user.Id,
		FullName:    form.FullName,
		PhoneNumber: form.PhoneNumber,
		Address:     form.Address,
	}

	err = updateUser.UpdateInfo()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_USER_FAIL, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

// @Summary 重置代理下的用户密码
// @Description 重置代理下的用户密码
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAgencyResetUserPasswordForm body request.ReqAgencyResetUserPasswordForm true "username:用户名"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/resetUserPassword  [post]
func ResetAgencyUserPassword(c *gin.Context) {
	var (
		form request.ReqAgencyResetUserPasswordForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	user, err := getUserByAgencyAndUsername(account.Agency, form.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	// 3、修改密码
	updateUser := user_service.User{
		Id:       user.Id,
		Password: util.EncodeMD5(util.DefaultPassword),
	}

	err = updateUser.UpdatePassword()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_PASSWORD, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

// @Summary 重置代理下的用户资金交易密码
// @Description 重置代理下的用户资金交易密码
// @Tags 代理平台
// @accept json
// @Produce  json
// @Param ReqAgencyResetUserTradeKeyForm body request.ReqAgencyResetUserTradeKeyForm true "username:用户名"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/agency/resetUserTradeKey  [post]
func ResetAgencyUserTradeKey(c *gin.Context) {
	var (
		form request.ReqAgencyResetUserTradeKeyForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	account := c.MustGet(util.TokenKey).(*util.Claims)

	user, err := getUserByAgencyAndUsername(account.Agency, form.Username)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if user == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	// 3、修改交易密码
	updateUser := user_service.User{
		Id:       user.Id,
		TradeKey: util.EncodeMD5(util.DefaultTradeKey),
	}

	err = updateUser.UpdateTradeKey()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_TRADE_PASSWORD, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

//获取用户信息
func getUserByAgencyAndUsername(agency string, username string) (*models.User, error) {
	userService := user_service.User{}
	return userService.GetByAgencyAndUsername(agency, username)
}
