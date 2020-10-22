package models

import (
	"QuickPass/pkg/util"
	"github.com/jinzhu/gorm"
)

type User struct {

	// 主键id
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 用户名
	UserName string `json:"user_name" gorm:"column:user_name"`
	// 登录密码
	Password string `json:"password" gorm:"column:password"`
	// 交易密码
	TradeKey string `json:"trade_key" gorm:"column:trade_key"`
	// 所属代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 姓名
	FullName string `json:"full_name" gorm:"column:full_name"`
	// 电话号码
	PhoneNumber string `json:"phone_number" gorm:"column:phone_number"`
	// 地址
	Address string `json:"address" gorm:"column:address"`
	// 账户类型(1:系统账户 2:普通账户)
	Type int `json:"type" gorm:"column:type"`
	// 角色（1:代理，2：承兑人、3：商家）
	Role int `json:"role" gorm:"column:role"`
	// token盐值
	SecretKey string `json:"secret_key" gorm:"column:secret_key"`
	// 账户状态（1：暂停，2：启用）
	Status int `json:"status" gorm:"column:status"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`
}

// 设置User的表名为`user`
func (User) TableName() string {
	return "user"
}

func NewUserModel(session *Session) *User {
	return &User{Session: session}
}

// ExistUserByID checks if an user exists based on ID
func (a *User) ExistUserByID(id int64) (bool, error) {
	user := new(User)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetUserTotal gets the total number of users based on the constraints
func (a *User) GetUserTotal(maps interface{}) (int64, error) {
	var count int64
	err := a.Session.db.Model(&User{}).Where(maps).Count(&count).Error
	return count, err
}

// GetUsers gets a list of users based on paging constraints
func (a *User) GetUsers(pageNum int, pageSize int, maps interface{}) ([]*User, error) {
	var users []*User
	err := a.Session.db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return users, nil
}

// GetUser Get a single user based on ID
func (a *User) GetUser(id int64) (*User, error) {
	user := new(User)
	err := a.Session.db.Where("id = ?", id).First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return user, nil
}

func (a *User) GetUserByUsernameAndAgency(agency string, username string) (*User, error) {
	user := new(User)
	err := a.Session.db.Where("user_name = ? and agency = ?", username, agency).First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return user, nil
}

func (a *User) GetAgencyGroup() ([]string, error) {
	var users []*User
	err := a.Session.db.Table(a.TableName()).Select("agency").Where("agency != ?", "").Group("agency").Find(&users).Error
	if err != nil {
		return nil, err
	}

	group := make([]string, 0, len(users))
	for _, item := range users {
		group = append(group, item.Agency)
	}

	return group, nil
}

// EditUser modify a single user
func (a *User) EditUser(id int64, data map[string]interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&User{}).Where("id = ?", id).Updates(data).Error
}

// AddUser add a single user
func (a *User) AddUser(user *User) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(user).Error
}

// DeleteUser delete a single user
func (a *User) DeleteUser(agency string, username string) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("agency = ? and user_name = ?", agency, username).Delete(User{}).Error
}
