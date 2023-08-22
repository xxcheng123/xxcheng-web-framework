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
	middlewares []MiddleWare
	tplEngine   TemplateEngine
}

// HTTPServerOption 给HTTPServer添加额外的功能
type HTTPServerOption func(server *HTTPServer)

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	s := &HTTPServer{
		router: NewRouter(),
	}
	// 执行HTTPServerOption 执行注入功能
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ServerWithMiddleWare 返回一个HTTPServerOption的回调函数，用于注入中间件
func ServerWithMiddleWare(mdls ...MiddleWare) HTTPServerOption {
	return func(server *HTTPServer) {
		if server.middlewares == nil {
			server.middlewares = mdls
		} else {
			server.middlewares = append(server.middlewares, mdls...)
		}
	}
}
func ServerWithTemplateEngine(engine TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = engine
	}
}

// ServeHTTP 业务逻辑处理入口
func (h *HTTPServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Req:       req,
		Resp:      resp,
		tplEngine: h.tplEngine,
	}
	root := h.serve
	// ms[0]->ms[1]->...->ms[n-1]->h.serve
	// func(next HandlerFunc) HandlerFunc
	// ms[n-1]:next=>h.serve
	// ms[n-2]:next=>ms[n-1]
	// MiddleWare必须调用next，否则后面的MiddleWare以及serve就无法运行了

	for i := len(h.middlewares) - 1; i >= 0; i-- {
		root = h.middlewares[i](root)
	}

	var flashRespMiddleWare = func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}

	root = flashRespMiddleWare(root)

	root(ctx)
}

// serve 内部负责业务逻辑调度的
// 根据路由查找对应的handlerFunc
func (h *HTTPServer) serve(ctx *Context) {
	nodeInfo, ok := h.FindRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		ctx.RespJSON(404, []byte("not found,path:"+ctx.Req.URL.Path))
		//ctx.Resp.WriteHeader(404)
		//ctx.Resp.Write([]byte("not found,path:" + ctx.Req.URL.Path))
		return
	}
	ctx.PathParams = nodeInfo.pathParams
	ctx.MatchedRoute = nodeInfo.node.route
	nodeInfo.node.handlerFunc(ctx)
}

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	if ctx.RespData != nil {
		ctx.Resp.Write(ctx.RespData)
	}
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
	if h.router == nil {
		panic("Router未初始化...")
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	//生命周期hook
	return http.Serve(l, h)
}
