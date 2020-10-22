package models

import (
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type Management struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理名称
	Agency string `json:"agency" gorm:"column:agency"`
	// 后台账号
	UserName string `json:"user_name" gorm:"column:user_name"`
	// 密码
	Password string `json:"password" gorm:"column:password"`
	// 姓名
	FullName string `json:"full_name" gorm:"column:full_name"`
	// 角色（1：管理员，2：客服，3：财务）
	Role int `json:"role" gorm:"column:role"`
	// 权限列表(如：order,user,merchant,acceptor)
	Rules string `json:"rules" gorm:"column:rules"`
	// 电话号码
	PhoneNumber string `json:"phone_number" gorm:"column:phone_number"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置Management的表名为`management`
func (Management) TableName() string {
	return "management"
}

func NewManagementModel(session *Session) *Management {
	return &Management{Session: session}
}

// ExistManagementByID checks if an management exists based on ID
func (a *Management) ExistManagementByID(id int64) (bool, error) {
	management := new(Management)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(management).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if management.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetManagementTotal gets the total number of managements based on the constraints
func (a *Management) GetManagementTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&Management{}).Where(maps).Count(&count).Error
	return count, err
}

// GetManagements gets a list of managements based on paging constraints
func (a *Management) GetManagements(pageNum int, pageSize int, maps interface{}) ([]*Management, error) {
	var managements []*Management
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&managements).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return managements, nil
}

// GetManagement Get a single management based on ID
func (a *Management) GetManagement(id int64) (*Management, error) {
	management := new(Management)
	err := a.Session.db.Where("id = ?", id).First(management).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return management, nil
}

func (a *Management) GetManagementByUsernameAndAgency(agency string, username string) (*Management, error) {
	management := new(Management)
	err := a.Session.db.Where("user_name = ? and agency = ?", username, agency).First(management).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return management, nil
}

// EditManagement modify a single management
func (a *Management) EditManagement(id int64, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&Management{}).Where("id = ?", id).Updates(data).Error
}

// AddManagement add a single management
func (a *Management) AddManagement(management *Management) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(management).Error
}

// DeleteManagement delete a single management
func (a *Management) DeleteManagement(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(Management{}).Error
}
