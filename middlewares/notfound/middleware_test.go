package notfound

import (
	"testing"
	"xxcheng_web_framework"
)

func TestNewNotFountBuilder(t *testing.T) {
	page := `
<html>
	<h1>404 NOT FOUND 我的自定义错误页面</h1>
</html>
`
	builder := NewNotFountBuilder().AddError(404, []byte(page)).Build()
	s := xxcheng_web_framework.NewHTTPServer(xxcheng_web_framework.ServerWithMiddleWare(builder))
	s.GET("/user", func(ctx *xxcheng_web_framework.Context) {
		ctx.RespData = []byte("hello, world")
	})

	s.Start(":8081")
}
