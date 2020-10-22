package token_cache

import (
	"QuickPass/pkg/constant"
	"QuickPass/pkg/gredis"
	"QuickPass/pkg/setting"
	"fmt"
)

// 缓存token
func SetCacheToken(agency, username, token string) error {
	cacheKey := formatTokenKey(agency, username)
	err := gredis.Set(cacheKey, token)
	if err != nil {
		return err
	}
	return gredis.SetKeyExpire(cacheKey, setting.AppSetting.TokenExpireTime)
}

// 获取缓存的token
func GetCacheToken(agency, username string) (string, error) {
	cacheKey := formatTokenKey(agency, username)
	return gredis.GetStringValue(cacheKey)
}

// 检查用户token存在
func CheckToken(agency, username string) (bool, error) {
	cacheKey := formatTokenKey(agency, username)
	return gredis.CheckKey(cacheKey)
}

func formatTokenKey(agency, username string) string {
	return fmt.Sprintf("%s%s-%s", constant.CacheTokenPre, agency, username)
}
