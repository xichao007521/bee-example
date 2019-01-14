package controllers

import (
	"context"
	"do-global.com/bee-example/error"
	"do-global.com/bee-example/services"
	"github.com/astaxie/beego"
	"strconv"
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
			ip = string(ip[0 : lastColon-1])
		}
	}
	return ip
}

func (t *BasicController) GetStringNE(key string) string {
	v := t.GetString(key)
	if v == "" {
		panic(myError.NewBizError(400, "param:"+key+" must not be empty"))
	}
	return v
}

func (t *BasicController) GetStringsNE(key string) []string {
	v := t.GetStrings(key)
	if len(v) == 0 {
		panic(myError.NewBizError(400, "param:"+key+" must not be empty array"))
	}
	return v
}

func (t *BasicController) GetInt(key string) int {
	v, err := strconv.Atoi(t.GetStringNE(key))
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int value"))
	}
	return v
}

func (t *BasicController) GetInt8(key string) int8 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 8)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int8 value"))
	}
	return int8(i64)
}

func (t *BasicController) GetUint8(key string) uint8 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 8)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint8 value"))
	}
	return uint8(u64)
}

func (t *BasicController) GetInt16(key string) int16 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 16)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int16 value"))
	}
	return int16(i64)
}

func (t *BasicController) GetUint16(key string) uint16 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 16)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint16 value"))
	}
	return uint16(u64)
}

func (t *BasicController) GetInt32(key string) int32 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 32)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int32 value"))
	}
	return int32(i64)
}

func (t *BasicController) GetUint32(key string) uint32 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 32)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint32 value"))
	}
	return uint32(u64)
}

func (t *BasicController) GetInt64(key string) int64 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 64)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int64 value"))
	}
	return int64(i64)
}

func (t *BasicController) GetUint64(key string) uint64 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 64)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint64 value"))
	}
	return uint64(u64)
}

func (t *BasicController) GetBool(key string) bool {
	v, err := strconv.ParseBool(t.GetStringNE(key))
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a bool value"))
	}
	return v
}

func (t *BasicController) GetFloat(key string) float64 {
	v, err := strconv.ParseFloat(t.GetStringNE(key), 64)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a float value"))
	}
	return v
}

var userService services.UserService
