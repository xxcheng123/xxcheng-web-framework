package notfound

import "xxcheng_web_framework"

type NotFountBuilder struct {
	resp map[int][]byte
}

func NewNotFountBuilder() *NotFountBuilder {
	return &NotFountBuilder{
		resp: map[int][]byte{},
	}
}

func (n *NotFountBuilder) AddError(code int, resp []byte) *NotFountBuilder {
	n.resp[code] = resp
	return n
}

func (n *NotFountBuilder) Build() xxcheng_web_framework.MiddleWare {
	return func(next xxcheng_web_framework.HandlerFunc) xxcheng_web_framework.HandlerFunc {
		return func(ctx *xxcheng_web_framework.Context) {
			next(ctx)
			resp, ok := n.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = resp
			}
		}
	}
}
