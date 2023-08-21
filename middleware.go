package xxcheng_web_framework

// MiddleWare
// 函数式的责任链模式（洋葱模式）
type MiddleWare func(next HandlerFunc) HandlerFunc
