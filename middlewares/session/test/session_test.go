package test

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
	"xxcheng_web_framework"
	"xxcheng_web_framework/middlewares/session/cookie"
	"xxcheng_web_framework/middlewares/session/memory"
)

func TestMemoryAndCookie(t *testing.T) {

	var sessionName = "xxcheng_session"

	propagator := cookie.NewPropagator(sessionName)
	store := memory.NewStore(time.Second * 10)

	mlds := []xxcheng_web_framework.MiddleWare{
		func(next xxcheng_web_framework.HandlerFunc) xxcheng_web_framework.HandlerFunc {
			return func(ctx *xxcheng_web_framework.Context) {
				sessID, err := propagator.Extract(ctx.Req)
				if err != nil {
					id := uuid.New().String()
					session, _ := store.Generate(ctx.Req.Context(), id)
					sessID = session.ID()
				}
				ctx.Set(sessionName, sessID)

				next(ctx)
				_sessID, _ := ctx.Get(sessionName)
				sessID = _sessID.(string)
				propagator.Inject(sessID, ctx.Resp)
				store.Refresh(ctx.Req.Context(), sessID)
			}
		},
	}

	s := xxcheng_web_framework.NewHTTPServer(xxcheng_web_framework.ServerWithMiddleWare(mlds...))

	s.POST("/login", func(ctx *xxcheng_web_framework.Context) {
		_sessID, _ := ctx.Get(sessionName)
		sessID := _sessID.(string)
		session, err := store.Get(ctx.Req.Context(), sessID)
		if err != nil {
			id := uuid.New().String()
			session, _ = store.Generate(ctx.Req.Context(), id)
			sessID = session.ID()
			ctx.Set(sessionName, sessID)
		}
		session.Set(ctx.Req.Context(), "username", "xxcheng")
		session.Set(ctx.Req.Context(), "sex", "男")
		ctx.RespJSON(http.StatusOK, map[string]any{
			"msg": "登录成功" + sessID,
		})
	})
	s.GET("/info", func(ctx *xxcheng_web_framework.Context) {
		_sessID, _ := ctx.Get(sessionName)
		sessID := _sessID.(string)
		session, err := store.Get(ctx.Req.Context(), sessID)
		if err != nil {
			id := uuid.New().String()
			session, _ = store.Generate(ctx.Req.Context(), id)
			sessID = session.ID()
			ctx.Set(sessionName, sessID)
			ctx.RespJSON(http.StatusOK, map[string]any{
				"msg": "请先登录",
			})
			return
		}
		user, err := session.Get(ctx.Req.Context(), "username")
		if err == nil {
			ctx.RespJSON(http.StatusOK, map[string]any{
				"msg": fmt.Sprintf("你好,%s", user),
			})
		} else {
			ctx.RespJSON(http.StatusOK, map[string]any{
				"msg": "请先登录",
			})
		}
	})

	s.GET("/logout", func(ctx *xxcheng_web_framework.Context) {
		_sessID, _ := ctx.Get(sessionName)
		sessID := _sessID.(string)
		store.Remove(ctx.Req.Context(), sessID)
		ctx.RespJSON(http.StatusOK, map[string]any{
			"msg": "退出成功",
		})
	})

	s.Start(":9998")
}
