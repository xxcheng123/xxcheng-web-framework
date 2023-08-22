package xxcheng_web_framework

import (
	"bytes"
	"context"
	"html/template"
)

var _ TemplateEngine = &GoTemplateEngine{}

type TemplateEngine interface {
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	res := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(res, tplName, data)
	return res.Bytes(), err
}
