package handlers

import (
	"context"
	"github.com/eko/gocache/v2/store"
)

type UserDetails struct {
	tinkId   string
	username string
}

func (s *CoreServer) PutUserInCache(sessId string, user UserDetails) {
	err := s.UserTokenCache.Set(context.TODO(), sessId, user, &store.Options{Cost: 2})
	if err != nil {
		s.Logger.Error(err)
	}
}

func (s *CoreServer) GetUserFromCache(sessId string) (*UserDetails, error) {
	user, err := s.UserTokenCache.Get(context.TODO(), sessId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	return user.(*UserDetails), nil
}
