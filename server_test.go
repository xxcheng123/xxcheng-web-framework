package xxcheng_web_framework

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.GET("/user/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("ok,path:/user/detail"))
	})
	s.GET("/user/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("ok,path:/user/*"))
	})
	s.GET("/api/:username", func(ctx *Context) {
		ctx.Resp.Write([]byte("ok,path:/api/" + ctx.PathParams["username"]))
	})
	s.GET("/api/:username/:age/:sex", func(ctx *Context) {
		msg := ""
		for k, v := range ctx.PathParams {
			msg = fmt.Sprintf("%s,%s=%s", msg, k, v)
		}
		ctx.Resp.Write([]byte(msg))
	})
	if err := s.Start(":9999"); err != nil {
		fmt.Println("err", err)
	}
}
