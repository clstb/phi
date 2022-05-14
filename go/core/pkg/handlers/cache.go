package handlers

import (
	"context"
	"github.com/eko/gocache/v2/store"
	"runtime/debug"
)

func (s *CoreServer) PutClientSessionTokenInCache(id string, token string) {
	err := s.UserTokenCache.Set(context.TODO(), id, token, &store.Options{Cost: 2})
	if err != nil {
		debug.PrintStack()
	}
}
