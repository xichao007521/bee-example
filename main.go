package main

import (
	_ "do-global.com/bee-example/cache"
	"do-global.com/bee-example/controllers"
	"do-global.com/bee-example/logger"
	_ "do-global.com/bee-example/logger"
	_ "do-global.com/bee-example/routers"
	"encoding/json"
	"github.com/astaxie/beego"
	"os"
)

func main() {
	appLoggerConf := logger.AppConfig
	content, _ := json.Marshal(appLoggerConf)
	beego.SetLogger("file", string(content))
	beego.Info("app started, pid", os.Getpid())
	beego.BConfig.RecoverFunc = controllers.CustomPanicRecover
	prepare()
	defer shutdownGraceful()
	beego.Run()
}

func shutdownGraceful() {
	// release resource, like db pool, redis pool
	beego.Info("exit graceful")
}

// do something before http run
func prepare() {
}
