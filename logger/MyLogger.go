package logger

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"log"
	"sync"
)

type FileConfig struct {
	FileName string
	MaxLines int    // 每个文件保存的最大行数，默认值 1000000
	MaxSize  int    // 每个文件保存的最大尺寸, 默认值是 1 << 28
	Daily    bool   // 是否按照每天 logrotate，默认是 true
	MaxDays  int    // 文件最多保存多少天，默认保存 7 天
	Rotate   bool   // 默认是 true
	Level    int    // 默认是 Trace 级别
	Perm     string // 日志文件权限
}

type AllConfig struct {
	FileLoggers map[string]FileConfig `toml:"file"`
}

var loggerConfig AllConfig
var lock sync.Mutex

var AppConfig FileConfig

func init() {
	readConfig()
	AppConfig = loggerConfig.FileLoggers["app"]

	buildCustomLogger()
}

// 自定义日志
var ForDebugLogger *logs.BeeLogger

func buildCustomLogger() {
	ForDebugConfig := loggerConfig.FileLoggers["for_debug"]
	configContent, _ := json.Marshal(ForDebugConfig)
	ForDebugLogger = logs.NewLogger()
	ForDebugLogger.SetLogger("file", string(configContent))
	ForDebugLogger.SetLogger(logs.AdapterConsole)
	ForDebugLogger.SetLogFuncCallDepth(1)
}

// read config
func readConfig() {
	lock.Lock()
	defer lock.Unlock()
	if len(loggerConfig.FileLoggers) == 0 {
		data, err := ioutil.ReadFile("./conf/logger.toml")
		if err != nil {
			log.Fatal(err)
		}
		var loggerToml AllConfig
		if _, err := toml.Decode(string(data), &loggerToml); err != nil {
			log.Fatal(err)
		}
		loggerConfig = loggerToml
	}
}
