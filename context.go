package xxcheng_web_framework

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}
