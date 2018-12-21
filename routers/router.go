package routers

import (
	"do-global.com/bee-example/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.AutoRouter(&controllers.UserController{})
    beego.Router("/testtest", &controllers.UserController{}, "*:Login")
}
