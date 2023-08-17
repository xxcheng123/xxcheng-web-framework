package xxcheng_web_framework

import (
	"strings"
)

type router struct {
	trees map[string]*node
}

type node struct {
	path        string
	children    map[string]*node
	handlerFunc HandlerFunc
}

func NewRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (r *router) AddRoute(method, path string, handlerFunc HandlerFunc) {
	//非法注册校验
	//路径不能为空
	if path == "" {
		panic("AddRoute：路径不能为空")
	}
	//必须/开头
	if path[0] != '/' {
		panic("AddRoute：路径必须 / 开头")
	}
	//不能/结尾
	if path != "/" && path[len(path)-1] == '/' {
		panic("AddRoute：路径不能 / 结尾")
	}

	//找到对应方法的树
	methodTree, ok := r.trees[method]
	//如果没有的话就创建一个
	if !ok {
		methodTree = &node{
			path: "/",
		}
		r.trees[method] = methodTree
	}
	//注册路由为根节点，特殊处理
	if path == "/" {
		if methodTree.handlerFunc != nil {
			panic("handlerFunc不能重复注册")
		}
		methodTree.handlerFunc = handlerFunc
		//必须使用return退出，strings.Split会分割一个空的""出来
		return
	}

	//分割path
	segs := strings.Split(path[1:], "/")
	//遍历
	root := methodTree
	for _, seg := range segs {
		//为空情况
		if seg == "" {
			panic("AddRoute：中间不能为连续//")
		}
		nextRoot, ok := root.childOf(seg)
		//获取中途节点，如果没有的话就创建一个中途节点
		if !ok {
			nextRoot = root.childCreateWithoutHandlerFunc(seg)
		}
		root = nextRoot
	}
	if root.handlerFunc != nil {
		panic("handlerFunc不能重复注册")
	}
	//注册handlerFunc
	root.handlerFunc = handlerFunc
}

func (r *router) FindRoute(method, path string) (*node, bool) {
	method = strings.Trim(method, "")
	path = strings.Trim(path, "")
	//非法校验
	if method == "" || path == "" {
		return nil, false
	}
	methodTree, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	//根节点
	if path == "/" {
		return methodTree, true
	}
	segs := strings.Split(path[1:], "/")
	currentNode := methodTree
	for _, seg := range segs {
		cn, found := currentNode.childOf(seg)
		if !found {
			return nil, false
		}
		currentNode = cn
	}
	return currentNode, currentNode.handlerFunc != nil
}

func (n *node) childOf(path string) (*node, bool) {
	cn, ok := n.children[path]
	return cn, ok
}

// childCreateWithoutHandlerFunc 创建子节点
func (n *node) childCreateWithoutHandlerFunc(path string) *node {
	//这里需要注意一个小细节
	/**
	map类型初始化默认值是nil
	此时不支持赋值，会报错
	但是可以取值，结果为nil
	参考：https://blog.csdn.net/qq_39920531/article/details/88103496
	***/
	if n.children == nil {
		n.children = map[string]*node{}
	}
	cn := &node{path: path}
	n.children[path] = cn
	return cn
}
