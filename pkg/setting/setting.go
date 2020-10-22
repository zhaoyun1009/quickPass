package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
	JwtSecret     string
	PageSize      int
	PageSizeLimit int
	// md5加密盐值
	MD5Salt string
	// token过期时间(单位秒)
	TokenExpireTime int64
	// token过期之前续签时间(单位秒)
	TokenRenewalTime int64
	// 买入订单超时分钟
	OrderTimeoutMinute int64
	// 卖出订单自动放行分钟
	SellOrderAutoDischargeMinute int64
	// 买单api对接回调页面地址
	OrderOpenApiBuyUrl string

	ImageMaxSize   int
	ImageAllowExts []string

	RuntimeRootPath string
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
	Debug       bool
}

var DatabaseSetting = &Database{}

type Redis struct {
	Url string
	// 最大空闲连接数
	MaxIdle int
	// 最大空闲连接时间
	IdleTimeoutSec int
	// 密码
	Password string
	Db       int
}

var RedisSetting = &Redis{}

type Minio struct {
	BucketName string
	AccessKey  string
	SecretKey  string
	Endpoint   string
	PreUrl     string
	UseSSL     bool
}

var MinioSetting = &Minio{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("minio", MinioSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
