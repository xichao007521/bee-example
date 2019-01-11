package controllers

import (
	"do-global.com/bee-example/error"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"runtime"
	"time"
)

func CustomPanicRecover(ctx *context.Context) {
	if err := recover(); err != nil {
		t := ctx.Input.GetData("_____t").(*BasicController)
		if err == beego.ErrAbort {
			t.Finish()
			return
		}
		if !beego.BConfig.RecoverPanic {
			t.Finish()
			panic(err)
		}

		switch err.(type) {
		case *myError.BizError:
			bizE := err.(*myError.BizError)
			handleBizError(ctx, bizE)
			t.Ctx.Input.SetData("___status", bizE.Code)
			t.Finish()
			return
		}

		var stack string
		logs.Critical("the request url is ", ctx.Input.URL())
		logs.Critical("Handler crashed with error", err)
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		}
		t.Finish()
		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}
	}
}

func handleBizError(ctx *context.Context, bizE *myError.BizError) {
	beego.Error("[BIZE] biz exception: code", bizE.Code, " message ", bizE.Message)
	rd := &ResponseData{
		Ret:        bizE.Code,
		Message:    bizE.Message,
		ServerTime: time.Now().UnixNano() / 1000000,
	}
	hasIndent := beego.BConfig.RunMode != beego.PROD
	jsonpCallback := ctx.Request.Form.Get("callback")
	if jsonpCallback != "" {
		ctx.Output.JSONP(rd, hasIndent)
	} else {
		ctx.Output.JSON(rd, hasIndent, false)
	}

}
