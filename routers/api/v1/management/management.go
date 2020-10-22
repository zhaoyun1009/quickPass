package management

import (
	"QuickPass/models"
	"QuickPass/pkg/app"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/e"
	"QuickPass/pkg/util"
	"QuickPass/request"
	"QuickPass/response"
	"QuickPass/service/management_service"
	"QuickPass/service/user_service"
	"github.com/gin-gonic/gin"
)

// @Summary 后台管理登录
// @Description 后台管理登录
// @Tags 后台管理
// @accept json
// @Produce  json
// @Param ReqManagementLoginForm body request.ReqManagementLoginForm true "user_name:用户名  password:密码 "
// @Success 200 {object} response.RespManagementLogin
// @Failure 500 {object}  app.Response
// @Router /api/v1/backend/management/login [post]
func Login(c *gin.Context) {
	var (
		form request.ReqManagementLoginForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 1、密码验证 返回角色 和权限列表
	// 1.1 获取用户信息
	userService := user_service.User{}
	user, err := userService.GetByAgencyAndUsername("", form.UserName)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	// 1.2 角色/密码验证
	if user == nil || user.Role != constant.SYSTEM_USEER || !util.MD5Equals(form.Password, user.Password) {
		app.ErrorResp(c, e.LOGIN_FAILED, "")
		return
	}

	role := int8(user.Role)
	// 2、token生成
	token, err := util.GenerateToken("", user.UserName, user.FullName, role)
	if err != nil {
		app.ErrorResp(c, e.ERROR_AUTH_TOKEN, "token生成失败")
		return
	}

	app.SuccessResp(c, response.RespManagementLogin{
		Token: token,
		Role:  role,
		Rules: "",
	})
}

type getManagementInfoListForm struct {
	//代理
	Agency string `form:"agency" json:"agency" valid:"Required"`
	//页面大小
	StartPage int `form:"start_page" json:"start_page"`
	//起始页
	PageSize int `form:"page_size" json:"page_size"`
}

// @Summary 获取后台管理员信息列表
// @Description 获取后台管理员信息列表
// @Tags 后台管理
// @accept json
// @Produce  json
// @Param agency query string true "代理"
// @Param start_page query int true "起始页"
// @Param page_size query int true "页面大小"
// @Success 200 {object}  response.RespManagementInfoList "response.RespManagementInfoList"
// @Failure 500 {object}  app.Response
// @Router /api/v1/backend/management/getManagementInfoList [get]
func GetManagementInfoList(c *gin.Context) {
	var (
		form getManagementInfoListForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	offset, limit := util.GetPaginationParams(form.StartPage, form.PageSize)
	managementService := management_service.Management{
		Agency:   form.Agency,
		Role:     -1,
		PageNum:  offset,
		PageSize: limit,
	}

	total, err := managementService.Count()
	if err != nil {
		app.ErrorResp(c, e.ERROR_COUNT_USER_FAIL, err.Error())
		return
	}

	users, err := managementService.GetAll()
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USERS_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, response.RespManagementInfoList{
		List:  users,
		Total: total,
	})
}

type modifyPasswordForm struct {
	Agency           string `form:"agency" json:"agency" valid:"Required"`
	UserName         string `form:"user_name" json:"user_name" valid:"Required"`
	OriginalPassword string `form:"original_password" json:"original_password" valid:"Required"`
	UpdatePassword   string `form:"update_password" json:"update_password" valid:"Required"`
}

// @Summary 修改登录密码
// @Description 修改登录密码
// @Tags 后台管理
// @accept json
// @Produce  json
// @Param agency body management.modifyPasswordForm true "agency代理 所属系统 默认admin"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/backend/management/modifyPassword [post]
func ModifyPassword(c *gin.Context) {
	var (
		form modifyPasswordForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// 1、获取用户信息
	management, err := getManagementByAgencyAndUsername(form.Agency, form.UserName)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if management == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	// 2、校验原始密码
	if !util.MD5Equals(form.OriginalPassword, management.Password) {
		app.ErrorResp(c, e.ERROR_VALIDATE_PASSWORD, "")
		return
	}

	// 3、修改密码
	updateManagement := management_service.Management{
		Id:       management.Id,
		Password: util.EncodeMD5(form.UpdatePassword),
	}

	err = updateManagement.UpdatePassword()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_PASSWORD, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

type updateManagementInfoForm struct {
	Agency      string `form:"agency" json:"agency" valid:"Required"`
	UserName    string `form:"user_name" json:"user_name" valid:"Required"`
	PhoneNumber string `form:"phone_number" json:"phone_number" `
	FullName    string `form:"full_name" json:"full_name" `
}

// @Summary 更新后台管理员账号信息
// @Description 更新后台管理员账号信息
// @Tags 后台管理
// @accept json
// @Produce  json
// @Param agency body management.updateManagementInfoForm true "代理 所属系统 默认admin"
// @Success 200 {object}  app.Response
// @Failure 500 {object}  app.Response
// @Router /api/v1/backend/management/updateManagementInfo [post]
func UpdateManagementInfo(c *gin.Context) {
	var (
		form updateManagementInfoForm
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}

	//1.获取用户信息
	management, err := getManagementByAgencyAndUsername(form.Agency, form.UserName)
	if err != nil {
		app.ErrorResp(c, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	if management == nil {
		app.ErrorResp(c, e.ERROR_NOT_EXIST_USER, "")
		return
	}

	//2.更新用户信息
	updateManagement := management_service.Management{
		Id:          management.Id,
		FullName:    form.FullName,
		PhoneNumber: form.PhoneNumber,
	}

	err = updateManagement.UpdateFullNameAndPhone()
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPDATE_USER_FAIL, err.Error())
		return
	}
	app.SuccessResp(c, nil)
}

//获取后台管理员信息
func getManagementByAgencyAndUsername(agency string, username string) (*models.Management, error) {
	managementService := management_service.Management{
		Agency:   agency,
		UserName: username,
	}
	return managementService.GetByAgencyAndUsername()
}
