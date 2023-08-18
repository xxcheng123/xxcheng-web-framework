package xxcheng_web_framework

import "net/http"

type PathParams map[string]string
type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter ``
	PathParams PathParams
}
