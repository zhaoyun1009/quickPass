package models

import (
	"QuickPass/pkg/setting"
	"fmt"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
)

//权限结构
type CasbinModel struct {
	RoleName string `json:"rolename"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}

//持久化到数据库
func Casbin() *casbin.Enforcer {
	//这里最后必须设置为true,不知道为撒会抛错
	a := gormadapter.NewAdapter("mysql", fmt.Sprintf("%s:%s@tcp(%s)/blog", setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password, setting.DatabaseSetting.Host), true)
	e := casbin.NewEnforcer("conf/auth_model.conf", a, true)
	e.EnableLog(true)

	e.LoadPolicy()
	return e
}
