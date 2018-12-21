package controllers

type UserController struct {
	BasicController
}

func (the *UserController) GetUser() {
	the.ok("123")
	//var tp http.RoundTripper = &http.Transport{
	//	DialContext: (&net.Dialer{
	//		Timeout:   30 * time.Second,
	//		KeepAlive: 30 * time.Second,
	//		DualStack: true,
	//	}).DialContext,
	//	MaxIdleConns:          100,
	//	IdleConnTimeout:       90 * time.Second,
	//	ExpectContinueTimeout: 1 * time.Second,
	//}
	//
	//req := httplib.Get("http://beego.me/")
	//req.Retries(1)
	//req.SetTransport(tp)
	//resp, _ := req.Response()
	//fmt.Println(resp.StatusCode)
	//
	//d := make(map[string]string)
	//d["1"] = "a"
	//d["2"] = "b"
	//the.renderJson(d)
}

func (the *UserController) Login() {
	user := userService.Login("1", "2")
	the.ok(user)
}
