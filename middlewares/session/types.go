package session

import (
	"context"
	"net/http"
)

type Store interface {
	Generate(ctx context.Context, id string) (Session, error)
	Get(ctx context.Context, id string) (Session, error)
	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
}

type Session interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
	ID() string
}

// Propagator 提取注入Session
type Propagator interface {
	Inject(id string, writer http.ResponseWriter) error
	Extract(req *http.Request) (string, error)
	Remove(writer http.ResponseWriter) error
}
