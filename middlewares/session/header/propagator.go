package header

import (
	"errors"
	"net/http"
)

// 由前端实现

type Propagator struct {
	sessionName string
}

func NewPropagator(sessionName string) *Propagator {
	return &Propagator{
		sessionName: sessionName,
	}
}

func (p *Propagator) Inject(id string, resp http.ResponseWriter) error {
	//resp.Header().Set(p.sessionName, id)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie := req.Header.Get(p.sessionName)
	if cookie == "" {
		return "", errors.New("cookie is nil")
	}
	return cookie, nil
}

func (p *Propagator) Remove(resp http.ResponseWriter) error {
	return nil
}
