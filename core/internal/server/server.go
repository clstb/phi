package server

import (
	"github.com/clstb/phi/core/internal/auth"
	"github.com/clstb/phi/pkg"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"time"
)

func createCache() *cache.Cache {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	ristrettoStore := store.NewRistretto(ristrettoCache, &store.Options{Expiration: 2 * time.Hour})
	cacheManager := cache.New(ristrettoStore)
	return cacheManager
}

/*
func (s *CoreServer) putUserInCache(sessId string, user UserDetails) {
  err := s.UserTokenCache.Set(context.TODO(), sessId, user, &store.Options{Cost: 2})
  if err != nil {
    s.Logger.Error(err)
  }
}

func (s *CoreServer) getUserFromCache(sessId string) (*UserDetails, error) {
  user, err := s.UserTokenCache.Get(context.TODO(), sessId)
  if err != nil {
    s.Logger.Error(err)
    return nil, err
  }
  return user.(*UserDetails), nil
}
*/

func NewServer(authClient *auth.Client) *CoreServer {
	log := pkg.CreateLogger()
	return &CoreServer{
		AuthClient: authClient,
		Logger:     log.Named("CORE"),
		//UserTokenCache: createCache(),
	}
}
