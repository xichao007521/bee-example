package controllers

import (
	"context"
	"do-global.com/public-server/globals"
	"do-global.com/public-server/logger"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func (t *BasicController) Prepare() {
	t.startTime = time.Now().UnixNano()

	// 设置request生命周期的参数
	ctx := context.TODO()
	ctx = globals.WithRequestID(ctx, rand.Intn(time.Now().Second() + 1))
	t.reqCtx = ctx
	// 访问权限检测
	checkAccess(t)
}

func (t *BasicController) Finish()  {
	// 删掉reqId相关资源
	globals.RemoveOrmer(t.reqCtx)
	t.reqCtx.Done()

	// access log
	spentTime := (time.Now().UnixNano() - t.startTime) / 1e6
	paramsStr, _ := json.Marshal(t.Ctx.Request.Form)
	reqPath := t.Ctx.Request.URL.Path
	now := time.Now()
	if t.Ctx.ResponseWriter.Status == 0 {
		t.Ctx.ResponseWriter.Status = http.StatusOK
	}
	accessInfo := fmt.Sprintf("%v\001%v\001%v\001%v\001%v\001%v", now.Format("20060102"), now.UnixNano() / 1e6, reqPath, string(paramsStr),
		t.Ctx.ResponseWriter.Status, spentTime)
	logger.AccessLogger.Info(accessInfo)
}


/***** 权限校验 *****/
// 白名单
var accessWhiteList = []string {
	"controllers.HealthController.Check",
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
		if strings.ToLower(controllerAndMethod) == strings.ToLower(whiteItem) {
			return
		}
	}
	// TODO
	//
	//appSecret := t.Ctx.Request.Header.Get("x-secret")
	//if appSecret == "" {
	//	t.forbidden()
	//	return
	//}
	//appId := t.GetString("app_id")
	//if appId == "" {
	//	t.forbidden()
	//	return
	//}
	//
	//productApp, err := productAppService.GetProductApp(t.reqCtx, appId, appSecret)
	//if err != nil || productApp.Id == 0 {
	//	t.forbidden()
	//	return
	//}
	//t.productApp = productApp
}

