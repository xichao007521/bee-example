package main

import (
	"do-global.com/bee-example/logger"
	_ "do-global.com/bee-example/routers"
	"encoding/json"
	"github.com/astaxie/beego"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	appLoggerConf := logger.AppConfig
	content, _ := json.Marshal(appLoggerConf)
	beego.SetLogger("file", string(content))
	go func() {
		beego.Run()
	}()

	beego.Info("app started, pid", os.Getpid())

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownGraceful()
}

func shutdownGraceful()  {
	// release resource, like db pool, redis pool
	beego.Info("exit graceful")
}

