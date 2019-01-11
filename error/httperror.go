package myError

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"strconv"
)

func init() {
	for _, code := range []string{"400", "401", "402", "403", "404", "405", "500", "501", "502", "503", "504", "417", "422"} {
		beego.ErrorHandler(code, buildErrorHandler(code))
	}
}

func buildErrorHandler(errCode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := strconv.Atoi(errCode)
		if err != nil {
			panic(fmt.Sprintf("illegal err code: %s", errCode))
		}

		w.Write([]byte(http.StatusText(status)))
	}
}
