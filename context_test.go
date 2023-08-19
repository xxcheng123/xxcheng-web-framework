package xxcheng_web_framework

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestContext_BindJSON(t *testing.T) {
	s := NewHTTPServer()
	wantUser := &TUser{
		Username: "jpc",
		Password: "12345678",
	}
	s.POST("/parseJSON", func(ctx *Context) {
		u := new(TUser)
		ctx.BindJSON(u)
		assert.Equal(t, wantUser, u)
	})
	go func() {
		s.Start(":9997")
	}()

	mockReq := httptest.NewRequest(http.MethodPost, "/parseJSON", strings.NewReader(`{"username":"jpc","password":"12345678"}`))
	mockResp := httptest.NewRecorder()

	s.ServeHTTP(mockResp, mockReq)

}

func TestContext_FormVal(t *testing.T) {
	s := NewHTTPServer()
	s.POST("/testFormVal", func(ctx *Context) {
		username, _ := ctx.FormVal("username")
		password, _ := ctx.FormVal("password")
		fromType, _ := ctx.FormVal("fromType")
		assert.Equal(t, "jpc", username)
		assert.Equal(t, "12345678", password)
		assert.Equal(t, "1", fromType)
	})
	go func() {
		s.Start(":9997")
	}()

	mockReq := httptest.NewRequest(http.MethodPost, "/testFormVal?username=xxcheng&fromType=1", strings.NewReader(`username=jpc&password=12345678`))
	mockReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mockResp := httptest.NewRecorder()

	s.ServeHTTP(mockResp, mockReq)
}

func TestContext_QueryVal(t *testing.T) {
	s := NewHTTPServer()
	s.POST("/testQueryVal", func(ctx *Context) {
		username, _ := ctx.QueryVal("username")
		fromType, _ := ctx.QueryVal("fromType")
		assert.Equal(t, "xxcheng", username)
		assert.Equal(t, "1", fromType)
	})
	go func() {
		s.Start(":9997")
	}()

	mockReq := httptest.NewRequest(http.MethodPost, "/testQueryVal?username=xxcheng&fromType=1", nil)
	mockResp := httptest.NewRecorder()

	s.ServeHTTP(mockResp, mockReq)
}

func TestContext_PathVal(t *testing.T) {
	s := NewHTTPServer()
	s.POST("/PathVal/:username/:fromType", func(ctx *Context) {
		username, _ := ctx.PathVal("username")
		fromType, _ := ctx.PathVal("fromType")
		assert.Equal(t, "xxcheng", username)
		assert.Equal(t, "2", fromType)
	})
	go func() {
		s.Start(":9997")
	}()

	mockReq := httptest.NewRequest(http.MethodPost, "/PathVal/xxcheng/2", nil)
	mockResp := httptest.NewRecorder()

	s.ServeHTTP(mockResp, mockReq)
}

func TestContext_RespJSON(t *testing.T) {
	s := NewHTTPServer()
	s.GET("/testRespJSON", func(ctx *Context) {
		u := &TUser{
			Username: "xxcheng",
			Password: "12345678",
		}
		ctx.RespJSON(200, u)
	})
	go func() {
		s.Start(":9997")
	}()

	mockReq := httptest.NewRequest(http.MethodGet, "/testRespJSON", nil)
	mockResp := httptest.NewRecorder()

	s.ServeHTTP(mockResp, mockReq)

	assert.Equal(t, mockResp.Result().StatusCode, 200)
	content, _ := io.ReadAll(mockResp.Result().Body)
	assert.Equal(t, []byte(`{"username":"xxcheng","password":"12345678"}`), content)
}
