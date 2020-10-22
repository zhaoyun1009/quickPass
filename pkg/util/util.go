package util

import (
	"QuickPass/pkg/setting"
	"bytes"
	"fmt"
	"github.com/rs/xid"
	"github.com/sony/sonyflake"
	"time"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

//全局唯一雪花算法id生成器
var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

// yyyyMMdd + 15位数唯一id
func GetUniqueNo() string {
	id, _ := flake.NextID()
	return fmt.Sprintf("%s%x", time.Now().Format(TIME_TEMPLATE_5), id)
}

func GetUniqueName() string {
	id := xid.New()
	return id.String()
}

func TransHtmlJson(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u0026"), []byte(""), -1)
	data = bytes.Replace(data, []byte("\\u003c"), []byte(""), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(""), -1)
	return data
}
