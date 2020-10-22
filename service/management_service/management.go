package management_service

import (
	"QuickPass/models"
	"QuickPass/pkg/util"
)

type Management struct {
	//
	Id int64
	// 代理名称
	Agency string
	// 后台账号
	UserName string
	// 密码
	Password string
	// 姓名
	FullName string
	// 角色（0：管理员，1：客服，2：财务）
	Role int
	// 权限列表(如：order,user,merchant,acceptor)
	Rules string
	// 电话号码
	PhoneNumber string
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int
}

func (a *Management) Add() error {
	management := &models.Management{
		Agency:      a.Agency,
		UserName:    a.UserName,
		Password:    a.Password,
		FullName:    a.FullName,
		Role:        a.Role,
		Rules:       a.Rules,
		PhoneNumber: a.PhoneNumber,
	}

	session := models.NewSession()
	if err := models.NewManagementModel(session).AddManagement(management); err != nil {
		return err
	}

	return nil
}

func (a *Management) UpdatePassword() error {
	session := models.NewSession()
	return models.NewManagementModel(session).EditManagement(a.Id, map[string]interface{}{
		"password": a.Password,
	})
}

func (a *Management) UpdateFullNameAndPhone() error {
	session := models.NewSession()
	return models.NewManagementModel(session).EditManagement(a.Id, map[string]interface{}{
		"full_name":    a.FullName,
		"phone_number": a.PhoneNumber,
	})
}

func (a *Management) Get() (*models.Management, error) {
	session := models.NewSession()
	management, err := models.NewManagementModel(session).GetManagement(a.Id)
	if err != nil {
		return nil, err
	}
	return management, nil
}

func (a *Management) GetByAgencyAndUsername() (*models.Management, error) {
	session := models.NewSession()
	return models.NewManagementModel(session).GetManagementByUsernameAndAgency(a.Agency, a.UserName)
}

func (a *Management) GetAll() ([]*models.Management, error) {
	session := models.NewSession()
	managements, err := models.NewManagementModel(session).GetManagements(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	return managements, nil
}

func (a *Management) Delete() error {
	session := models.NewSession()
	return models.NewManagementModel(session).DeleteManagement(a.Id)
}

func (a *Management) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewManagementModel(session).ExistManagementByID(a.Id)
}

func (a *Management) Count() (int64, error) {
	session := models.NewSession()
	return models.NewManagementModel(session).GetManagementTotal(a.getMaps())
}

func (a *Management) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.UserName != "" {
		maps["user_name"] = a.UserName
	}
	if a.Password != "" {
		maps["password"] = a.Password
	}
	if a.FullName != "" {
		maps["full_name"] = a.FullName
	}
	if a.Role != -1 {
		maps["role"] = a.Role
	}
	if a.Rules != "" {
		maps["rules"] = a.Rules
	}
	if a.PhoneNumber != "" {
		maps["phone_number"] = a.PhoneNumber
	}

	return maps
}
