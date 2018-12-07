package controllers

import (
	"do-global.com/bee-example/services"
	"github.com/astaxie/beego"
	"time"
)

type BasicController struct {
	beego.Controller
}

// 统一返回值
type ResponseData struct {
	Ret int `json:"ret"`
	Message string `json:"message"`
	Result interface{} `json:"result"`
	ServerTime int64 `json:"serverTime"`
}

func (t *BasicController) renderJson(d interface{})  {
	t.SetData(d)
	callbackStr := t.GetString("callback", "")
	if callbackStr != "" {
		t.ServeJSONP()
	} else {
		t.ServeJSON()
	}
}

func (t *BasicController) ok(d interface{})  {
	rd := &ResponseData{
		Ret: 200,
		Message: "ok",
		Result: d,
		ServerTime: time.Now().UnixNano() / 1000000,
	}
	t.renderJson(rd)
}

func (t *BasicController) forbidden(d ...interface{})  {
	t.SetData(d)
	t.Ctx.ResponseWriter.Status = 403
	t.ServeJSON()
}

var userService services.UserService
