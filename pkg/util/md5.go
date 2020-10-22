package util

import (
	"QuickPass/pkg/setting"
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	m.Write([]byte(setting.AppSetting.MD5Salt))

	return hex.EncodeToString(m.Sum(nil))
}

// 判断str的MD5加密是否为md5Str
func MD5Equals(str, md5Str string) bool {
	return EncodeMD5(str) == md5Str
}
