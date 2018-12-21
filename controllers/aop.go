package controllers

import (
	"github.com/astaxie/beego"
	"reflect"
	"regexp"
	"unsafe"
)

func (t *BasicController) Prepare() {
	// 访问权限检测
	checkAccess(t)
}

/***** 权限校验 *****/
// 白名单
var accessWhiteList = [] string{
	"controllers.UserController.Login",
}

func (t *BasicController) GetRequestControllerAndMethods() (reflect.Type, map[string]string, bool) {
	cInfo, isFind := beego.BeeApp.Handlers.FindRouter(t.Ctx)
	if isFind {
		controllerInfoV := reflect.ValueOf(cInfo).Elem()

		controllerTypeV := controllerInfoV.Field(1)
		controllerTypeV = reflect.NewAt(controllerTypeV.Type(), unsafe.Pointer(controllerTypeV.UnsafeAddr()))
		controllerType := controllerTypeV.Interface().(*reflect.Type)

		methodsV := controllerInfoV.Field(2)
		methodsV = reflect.NewAt(methodsV.Type(), unsafe.Pointer(methodsV.UnsafeAddr()))
		methods := methodsV.Interface().(*map[string]string)

		return *controllerType, *methods, true
	}
	return nil, nil, false
}

func checkAccess(t *BasicController) {
	needCheck := beego.AppConfig.DefaultBool("secure.control_check", true)
	if !needCheck {
		return
	}
	controllerType, methods, isFind := t.GetRequestControllerAndMethods()
	if !isFind {
		t.forbidden()
		return
	}

	controllerName := controllerType.String()
	var methodName string
	for _, v := range methods {
		methodName = v
		break
	}
	controllerAndMethod := controllerName + "." + methodName
	for _, whiteItem := range accessWhiteList {
		if controllerAndMethod == whiteItem {
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
