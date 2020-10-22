package logf

import (
	"QuickPass/pkg/setting"
	"fmt"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// getLogFileName get the save name of the log file
func getLogFileName(dateTime string) string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		dateTime,
		setting.AppSetting.LogFileExt,
	)
}
