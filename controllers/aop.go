package controllers

import (
	"github.com/astaxie/beego"
	"regexp"
)

func (t *BasicController) Prepare()  {
	// 访问权限检测
	checkAccess(t)
}

/***** 权限校验 *****/
// 白名单
var accessWhiteList = [] string {
	"UserController.Login",
}

func checkAccess(t *BasicController) {
	needCheck := beego.AppConfig.DefaultBool("secure.control_check", true)
	if !needCheck {
		return
	}

	requestUri := getUriWithoutParams(t.Ctx.Request.RequestURI)
	for _, whiteItem := range accessWhiteList {
		if requestUri == beego.URLFor(whiteItem) {
			return
		}
	}

	token := t.Ctx.Request.Header.Get("x-token")
	if token == "" {
		t.forbidden()
	}

	// TODO 判断用户

}

func getUriWithoutParams(requestUri string) string {
	reg := regexp.MustCompile(`[?|#].*$`)
	return reg.ReplaceAllLiteralString(requestUri, "")
}
