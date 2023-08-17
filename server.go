package xxcheng_web_framework

import (
	"net"
	"net/http"
)

var _ Server = &HTTPServer{}

type HandlerFunc func(ctx *Context)

// Server 核心抽象接口保持小而美
type Server interface {
	http.Handler
	Start(addr string) error
	// AddRoute 路由注册
	AddRoute(method, path string, handlerFunc HandlerFunc)
}

// HTTPServer Server的实现
type HTTPServer struct {
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: NewRouter(),
	}
}

// ServeHTTP 业务逻辑处理入口
func (h *HTTPServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Req:  req,
		Resp: resp,
	}
	h.serve(ctx)
}

// serve 内部负责业务逻辑调度的
// 根据路由查找对应的handlerFunc
func (h *HTTPServer) serve(ctx *Context) {

}

// 根据具体HTTP方法，提供批量的请求方法

func (h *HTTPServer) GET(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodGet, path, handlerFunc)
}
func (h *HTTPServer) HEAD(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodHead, path, handlerFunc)
}
func (h *HTTPServer) POST(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodPost, path, handlerFunc)
}
func (h *HTTPServer) PUT(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodPut, path, handlerFunc)
}
func (h *HTTPServer) PATCH(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodPatch, path, handlerFunc)
}
func (h *HTTPServer) DELETE(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodDelete, path, handlerFunc)
}
func (h *HTTPServer) CONNECT(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodConnect, path, handlerFunc)
}
func (h *HTTPServer) OPTIONS(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodOptions, path, handlerFunc)
}
func (h *HTTPServer) TRACE(path string, handlerFunc HandlerFunc) {
	h.AddRoute(http.MethodTrace, path, handlerFunc)
}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	//生命周期hook
	return http.Serve(l, h)
}
