package cookie

import (
	"net/http"
	"xxcheng_web_framework/middlewares/session"
)

// Propagator 使用 Cookie 实现
type Propagator struct {
	sessionName string
}

var _ session.Propagator = &Propagator{}

func NewPropagator(sessionName string) *Propagator {
	return &Propagator{
		sessionName: sessionName,
	}
}

func (p *Propagator) Inject(id string, resp http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  p.sessionName,
		Value: id,
	}
	http.SetCookie(resp, cookie)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.sessionName)
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

func (p *Propagator) Remove(resp http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   p.sessionName,
		MaxAge: -1,
	}
	http.SetCookie(resp, cookie)
	return nil
}
