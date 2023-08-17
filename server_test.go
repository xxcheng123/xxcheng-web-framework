package xxcheng_web_framework

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	s := &HTTPServer{}
	s.GET("/user/detail", func(ctx *Context) {

	})
	if err := s.Start(":9999"); err != nil {
		fmt.Println("err", err)
	}
}
