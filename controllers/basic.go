package controllers

import (
	"context"
	"do-global.com/bee-example/services"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

type BasicController struct {
	beego.Controller

	// 请求开始时间
	startTime int64

	// 上下文
	reqCtx context.Context
}

// 统一返回值
type ResponseData struct {
	Ret        int         `json:"ret"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
	ServerTime int64       `json:"serverTime"`
}

func (t *BasicController) renderJson(d interface{}) {
	t.SetData(d)
	callbackStr := t.GetString("callback", "")
	if callbackStr != "" {
		t.ServeJSONP()
	} else {
		t.ServeJSON()
	}
}

func (t *BasicController) ok(d interface{}) {
	rd := &ResponseData{
		Ret:        200,
		Message:    "ok",
		Result:     d,
		ServerTime: time.Now().UnixNano() / 1000000,
	}
	t.renderJson(rd)
}

func (t *BasicController) Error403() {
	t.Abort("403")
}
func (t *BasicController) Error400() {
	t.Abort("400")
}
func (t *BasicController) Error500() {
	t.Abort("500")
}

func (t *BasicController) getRealIp() string {
	var ip string
	ip = t.Ctx.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = t.Ctx.Request.Header.Get("Client-Ip")
	}
	if ip == "" {
		ip = t.Ctx.Request.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = t.Ctx.Request.RemoteAddr
		lastColon := strings.LastIndex(ip, ":")
		if lastColon > -1 {
			ip = string(ip[0:lastColon - 1])
		}
	}
	return ip
}


var userService services.UserService
