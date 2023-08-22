package xxcheng_web_framework

import (
	"encoding/json"
	"errors"
	"net/http"
)

type PathParams map[string]string
type Context struct {
	Req          *http.Request
	Resp         http.ResponseWriter ``
	PathParams   PathParams
	QueryParams  map[string][]string
	MatchedRoute string

	RespData       []byte
	RespStatusCode int

	tplEngine TemplateEngine
}

func (c *Context) BindJSON(jsonModel any) error {
	if jsonModel == nil {
		return errors.New("传入的模型为nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	decoder.Decode(jsonModel)
	return nil
}

func (c *Context) FormVal(key string) (string, error) {
	if err := c.Req.ParseForm(); err != nil {
		return "", err
	}
	if vs, ok := c.Req.Form[key]; !ok {
		return "", errors.New("获取" + key + "内容失败")
	} else {
		return vs[0], nil
	}

}

func (c *Context) QueryVal(key string) (string, error) {
	if c.QueryParams == nil {
		//缓存
		c.QueryParams = c.Req.URL.Query()
	}
	if vs, ok := c.QueryParams[key]; !ok {
		return "", errors.New("获取" + key + "内容失败")
	} else {
		return vs[0], nil
	}
}

func (c *Context) PathVal(key string) (string, error) {
	if vs, ok := c.PathParams[key]; !ok {
		return "", errors.New("获取" + key + "内容失败")
	} else {
		return vs, nil
	}
}

func (c *Context) RespJSON(statusCode int, jsonData any) error {
	res, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	c.RespData = res
	c.RespStatusCode = statusCode
	//c.Resp.WriteHeader(statusCode)
	//c.Resp.Write(res)
	return nil
}

func (c *Context) Render(tpl string, data any) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tpl, data)
	c.RespStatusCode = 200
	if err != nil {
		c.RespStatusCode = 500
	}
	return err
}
