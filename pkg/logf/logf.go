package logf

import (
	"QuickPass/pkg/file"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type logFileWriter struct {
	file *os.File
	// 日期
	date string
}

func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	dateTime := time.Now().Format(setting.AppSetting.TimeFormat)
	if p.file == nil {
		filePath := getLogFilePath()
		fileName := getLogFileName(dateTime)
		p.date = dateTime
		p.file, _ = file.MustOpen(fileName, filePath)
	}
	n, e := p.file.Write(data)
	//按照每天生成日志文件
	if p.date != dateTime {
		p.file.Close()

		filePath := getLogFilePath()
		fileName := getLogFileName(dateTime)
		p.file, _ = file.MustOpen(fileName, filePath)
	}
	return n, e
}

var loggerEntry *logrus.Entry

// 全局logrus实例
var logger = logrus.New()

func Setup() {
	// 为当前logrus实例设置消息的输出，同样地，
	// 可以设置logrus实例的输出到任意io.writer
	logger.Out = os.Stdout

	logger.SetOutput(io.MultiWriter(os.Stdout, &logFileWriter{}))

	// 为当前logrus实例设置消息输出格式为json格式。
	// 同样地，也可以单独为某个logrus实例设置日志级别和hook，这里不详细叙述。
	//logf.Formatter = &logrus.TextFormatter{}
	//logf.Formatter.(*logrus.TextFormatter).TimestampFormat = "2006-01-02 15:04:05"
	////logf.Formatter.(*logrus.TextFormatter).DisableSorting = true
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: util.TIME_TEMPLATE_1,
	})

	// todo:这里看是否需要hook
	//logf.AddHook(&DefaultFieldHook{})
	if setting.ServerSetting.RunMode == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	loggerEntry = logger.WithFields(logrus.Fields{})
	Info("logf SetUp Successfully")
}

func Log(level logrus.Level, args ...interface{}) {
	loggerEntry.Log(level, args)
}

func Trace(args ...interface{}) {
	loggerEntry.Trace(args)
}

func Warning(args ...interface{}) {
	loggerEntry.Warning(args)
}

func Debug(args ...interface{}) {
	loggerEntry.Debug(args)
}

func Info(args ...interface{}) {
	loggerEntry.Info(args)
}

func Error(args ...interface{}) {
	loggerEntry.Error(args)
}

func Fatal(args ...interface{}) {
	loggerEntry.Fatal(args)
}

func DebugWithFields(fields logrus.Fields) {
	loggerEntry.WithFields(fields).Debug()
}

func InfoWithFields(fields logrus.Fields) {
	loggerEntry.WithFields(fields).Info()
}

func ErrorWithFields(fields logrus.Fields) {
	loggerEntry.WithFields(fields).Error()
}

type DefaultFieldHook struct {
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["appName"] = "MyAppName"
	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

type GormLogger struct {
	name   string
	logger *logrus.Logger
}

func (l *GormLogger) Print(values ...interface{}) {
	entry := l.logger.WithField("name", l.name)
	if len(values) > 1 {

		level := values[0]
		source := values[1]
		entry = entry.WithField("source", source)
		if level == "sql" {
			duration := values[2]
			// sql
			var formattedValues []interface{}
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", string(b)))
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}
			entry.WithField("spend", fmt.Sprintf("%13v", duration)).
				WithField("return", fmt.Sprintf("[%v]", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned ")).
				Info(fmt.Sprintf(sqlRegexp.ReplaceAllString(values[3].(string), "%v"), formattedValues...))
		} else {
			entry.Error(values[2:]...)
		}
	} else {
		entry.Error(values...)
	}

}

// New Create new logger
func NewGormLogger() *GormLogger {
	return newWithName("gorm_logger")
}

// NewWithName Create new logger with custom name
func newWithName(name string) *GormLogger {
	return newWithNameAndLogger(name, logger)
}

// NewWithNameAndLogger Create new logger with custom name and logger
func newWithNameAndLogger(name string, logger *logrus.Logger) *GormLogger {
	return &GormLogger{name: name, logger: logger}
}
