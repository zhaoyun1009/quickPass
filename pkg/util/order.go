package util

import (
	"math/rand"
	"strconv"
	"time"
)

func GetId() (int64, error) {
	now := time.Now()
	nowTimeStamp := now.UnixNano()
	rand.Seed(nowTimeStamp)
	i := rand.Intn(9999)

	timeString := now.Format(TIME_TEMPLATE_5)

	result, err := strconv.ParseInt(timeString, 10, 64)

	if err != nil {
		return result + int64(i*10000), err
	}
	return result + int64(i*10000), nil

}
