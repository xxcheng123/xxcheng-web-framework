package memory

import (
	"context"
	"errors"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
	"xxcheng_web_framework/middlewares/session"
)

var (
	ErrorSessionNotExist   = errors.New("store SessionNotExist")
	ErrorSessionNotMatched = errors.New("session SessionNotMatched")
)

type Store struct {
	mu         sync.RWMutex
	c          *cache.Cache
	expiration time.Duration
}

var _ session.Store = &Store{}

func NewStore(expiration time.Duration) *Store {
	store := &Store{
		c: cache.New(expiration, time.Second),
	}
	return store
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	session := &Session{
		id:   id,
		data: map[string]string{},
	}
	s.c.Set(id, session, s.expiration)
	return session, nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	session, ok := s.c.Get(id)
	if !ok {
		return nil, ErrorSessionNotExist
	}
	return session.(*Session), nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	session, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	s.c.Set(id, session, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.c.Delete(id)
	return nil
}

type Session struct {
	id   string
	data map[string]string
}

var _ session.Session = &Session{}

func (s *Session) Get(ctx context.Context, key string) (string, error) {
	val, ok := s.data[key]
	if !ok {
		return "", ErrorSessionNotMatched
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val string) error {
	s.data[key] = val
	return nil
}

func (s *Session) ID() string {
	return s.id
}
