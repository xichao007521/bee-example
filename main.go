package main

import (
	_ "do-global.com/sticker-api/error"
	_ "do-global.com/sticker-api/logger"
	_ "do-global.com/sticker-api/routers"
	"github.com/astaxie/beego"
	"os"
)

func main() {
	beego.Info("app started, pid", os.Getpid())
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