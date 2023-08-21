package accesslog

import (
	"fmt"
	"net/http"
	"testing"
	"xxcheng_web_framework"
)

func TestMiddleWareBuilder_Build(t *testing.T) {
	builder := MiddleWareBuilder{}
	//配置builder的输出
	builder.LogFunc(func(log string) {
		fmt.Println(log)
	})
	server := xxcheng_web_framework.NewHTTPServer(xxcheng_web_framework.ServerWithMiddleWare(builder.Build()))
	server.GET("/user/:id", func(ctx *xxcheng_web_framework.Context) {
		fmt.Println("hello", ctx.PathParams["id"])
	})
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8888/user/jpc", nil)
	//模拟一个请求
	server.ServeHTTP(nil, req)
}
