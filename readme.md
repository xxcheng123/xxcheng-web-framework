## 深入浅出实现*HTTPServer*

## 抽象
- Server
- Router
- Context
- Middleware

## 基石
*net/http*

## 主流框架
- beego
- gin
- echo
- iris

## 实现
### Context
#### 输入解析
#### 输出处理
### Middleware
### 模板引擎
```go
type MiddleWare func(next HandlerFunc) HandlerFunc
```
*example*
```go
package accesslog

import (
	"encoding/json"
	"xxcheng_web_framework"
)

type MiddleWareBuilder struct {
	logFunc func(log string)
}

func (m *MiddleWareBuilder) LogFunc(fn func(log string)) *MiddleWareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddleWareBuilder) Build() xxcheng_web_framework.MiddleWare {
	return func(next xxcheng_web_framework.HandlerFunc) xxcheng_web_framework.HandlerFunc {
		return func(ctx *xxcheng_web_framework.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				str, _ := json.Marshal(l)
				m.logFunc(string(str))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host       string `json:"host,omitempty"`
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
```