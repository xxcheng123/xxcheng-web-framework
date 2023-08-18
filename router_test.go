package xxcheng_web_framework

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

// 学习尝试TDD测试驱动开发

// TestRouter_AddRoute 注册路由
func TestRouter_AddRoute(t *testing.T) {
	//构建路由树
	list := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/detail",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/user/detail",
		},
		{
			method: http.MethodGet,
			path:   "/user/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/api/:smile",
		},
	}
	var mockHandlerFunc HandlerFunc = func(ctx *Context) {

	}
	//注册路由树
	r := NewRouter()
	for _, route := range list {
		r.AddRoute(route.method, route.path, mockHandlerFunc)
	}
	//预期路由树
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:        "/",
				handlerFunc: mockHandlerFunc,
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"detail": &node{
								path:        "detail",
								handlerFunc: mockHandlerFunc,
							},
						},
						handlerFunc: mockHandlerFunc,
						starChild: &node{
							path:        "*",
							handlerFunc: mockHandlerFunc,
						},
					},
					"api": &node{
						path: "api",
						paramChild: &node{
							path:        ":smile",
							handlerFunc: mockHandlerFunc,
						},
					},
				},
				starChild: &node{
					path:        "*",
					handlerFunc: mockHandlerFunc,
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"detail": &node{
								path:        "detail",
								handlerFunc: mockHandlerFunc,
							},
						},
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	//空路径注册
	r = NewRouter()
	assert.Panics(t, func() {
		r.AddRoute(http.MethodGet, "", mockHandlerFunc)
	})
	r = NewRouter()
	assert.Panics(t, func() {
		r.AddRoute(http.MethodGet, "user", mockHandlerFunc)
	})
	r = NewRouter()
	assert.PanicsWithValue(t, "AddRoute：路径不能 / 结尾", func() {
		r.AddRoute(http.MethodGet, "/user/", mockHandlerFunc)
	})
	r = NewRouter()
	assert.PanicsWithValue(t, "AddRoute：中间不能为连续//", func() {
		r.AddRoute(http.MethodGet, "/user////detail", mockHandlerFunc)
	})
	r = NewRouter()
	r.AddRoute(http.MethodGet, "/age", mockHandlerFunc)
	assert.PanicsWithValue(t, "handlerFunc不能重复注册", func() {
		r.AddRoute(http.MethodGet, "/age", mockHandlerFunc)
	})
	r = NewRouter()
	r.AddRoute(http.MethodGet, "/*", mockHandlerFunc)
	assert.PanicsWithValue(t, "通配符路由已注册，不允许注册参数路由", func() {
		r.AddRoute(http.MethodGet, "/:abc", mockHandlerFunc)
	})
	r = NewRouter()
	r.AddRoute(http.MethodGet, "/:abc", mockHandlerFunc)
	assert.PanicsWithValue(t, "参数路由已注册，不允许注册通配符路由", func() {
		r.AddRoute(http.MethodGet, "/*", mockHandlerFunc)
	})
}

// TestRouter_FindRoute 查找路由单元测试
func TestRouter_FindRoute(t *testing.T) {
	//构建路由树
	list := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/detail",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/*",
		},
		{
			method: http.MethodPost,
			path:   "/user/detail",
		},
		{
			method: http.MethodHead,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user/detail/:username",
		},
	}
	var mockHandlerFunc HandlerFunc = func(ctx *Context) {

	}
	//注册路由树
	r := NewRouter()
	for _, route := range list {
		r.AddRoute(route.method, route.path, mockHandlerFunc)
	}
	testCases := []struct {
		name     string
		method   string
		path     string
		found    bool
		nodeInfo *NodeInfo
	}{
		//根节点测试用例
		{
			name:   "head root",
			method: http.MethodHead,
			path:   "/",
			found:  true,
			nodeInfo: &NodeInfo{
				node: &node{
					path:        "/",
					handlerFunc: mockHandlerFunc,
				},
			},
		},
		//末尾节点
		{
			name:   "get user detail",
			method: http.MethodGet,
			path:   "/user/detail",
			found:  true,
			nodeInfo: &NodeInfo{
				node: &node{
					path:        "detail",
					handlerFunc: mockHandlerFunc,
					paramChild: &node{
						path:        ":username",
						handlerFunc: mockHandlerFunc,
					},
				},
			},
		},
		//中间节点
		{
			name:   "get user",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			nodeInfo: &NodeInfo{
				node: &node{
					path:        "user",
					handlerFunc: mockHandlerFunc,
					children: map[string]*node{
						"detail": &node{
							path:        "detail",
							handlerFunc: mockHandlerFunc,
							paramChild: &node{
								path:        ":username",
								handlerFunc: mockHandlerFunc,
							},
						},
					},
					starChild: &node{
						path:        "*",
						handlerFunc: mockHandlerFunc,
					},
				},
			},
		},
		//不存在的节点
		{
			name:   "get user info",
			method: http.MethodGet,
			path:   "/users/info",
			found:  false,
		},
		//通配符匹配
		{
			name:   "star user ",
			method: http.MethodGet,
			path:   "/user/xxx",
			found:  true,
			nodeInfo: &NodeInfo{
				node: &node{
					path:        "*",
					handlerFunc: mockHandlerFunc,
				},
			},
		},
		//参数路由
		{
			name:   "get user detail :username",
			method: http.MethodGet,
			path:   "/user/detail/xxcheng",
			found:  true,
			nodeInfo: &NodeInfo{
				node: &node{
					path:        ":username",
					handlerFunc: mockHandlerFunc,
				},
				pathParams: PathParams{
					"username": "xxcheng",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			nodeInfo, found := r.FindRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.found, found)
			if testCase.found {
				_, ok := testCase.nodeInfo.equal(nodeInfo)
				assert.True(t, ok)
			}
		})
	}
}

func (n *NodeInfo) equal(n2 *NodeInfo) (string, bool) {
	msg, ok := n.node.equal(n2.node)
	if !ok {
		return msg, ok
	}
	if n.pathParams == nil && n2.pathParams == nil {
		return "", true
	}
	msg2, ok2 := n.pathParams.equal(n2.pathParams)
	if !ok2 {
		return msg2, ok2
	}
	return "", true
}

// equal 路由对比
func (r *router) equal(r2 *router) (string, bool) {
	for method, tree := range r.trees {
		tree2, ok := r2.trees[method]
		if !ok {
			return fmt.Sprintf("没有对应的方法树：%s", method), false
		}
		if msg, eq := tree.equal(tree2); !eq {
			return msg, eq
		}
	}
	return "", true
}

// equal 节点对比
func (n *node) equal(n2 *node) (string, bool) {
	//空节点
	if n2 == nil {
		return "空节点nil", false
	}
	if n2.starChild != nil {
		msg, ok := n2.starChild.equal(n.starChild)
		if !ok {
			return msg, ok
		}
	}
	if n2.paramChild != nil {
		msg, ok := n2.paramChild.equal(n.paramChild)
		if !ok {
			return msg, ok
		}
	}
	//路径不匹配
	if n.path != n2.path {
		return fmt.Sprintf("路径不匹配：%s!=%s", n.path, n2.path), false
	}
	nHandlerFunc := reflect.ValueOf(n.handlerFunc)
	n2HandlerFunc := reflect.ValueOf(n2.handlerFunc)
	if nHandlerFunc != n2HandlerFunc {
		return fmt.Sprintf("%s节点handlerFunc 不匹配", n.path), false
	}

	//孩子长度不匹配
	if len(n.children) != len(n2.children) {
		return fmt.Sprintf("孩子长度不匹配:%s", n.path), false
	}
	//孩子内容不匹配
	for k, c := range n.children {
		c2, ok := n2.children[k]
		//获取不到对应的孩子
		if !ok {
			return "路径孩子不匹配", false
		}
		if msg, ok := c.equal(c2); !ok {
			return msg, false
		}
	}
	return "", true
}

func (p PathParams) equal(p2 PathParams) (string, bool) {
	if p2 == nil {
		return "对比数据为空", false
	}
	for k, v := range p {
		if v2, ok := p2[k]; !ok {
			return fmt.Sprintf("参数对比为空，%s!=%s", v, v2), false
		} else if v != v2 {
			return fmt.Sprintf("参数对比不正确，%s!=%s", v, v2), false
		}
	}
	return "", true
}
