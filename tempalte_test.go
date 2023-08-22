package xxcheng_web_framework

import (
	"html/template"
	"testing"
)

func TestGoTemplateEngine_Render(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	s := NewHTTPServer(ServerWithTemplateEngine(&GoTemplateEngine{T: tpl}))
	s.GET("/login", func(ctx *Context) {
		er := ctx.Render("login.gohtml", nil)
		if er != nil {
			t.Fatal(er)
		}
	})
	err = s.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}
