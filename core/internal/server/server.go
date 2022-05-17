package server

import (
	"github.com/clstb/phi/core/internal/auth"
	"github.com/clstb/phi/pkg"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"go.uber.org/zap"
	"time"
)

type CoreServer struct {
	AuthClient *auth.Client
	Logger     *zap.SugaredLogger
	LedgerUri  string
	TinkGwUri  string
	//UserTokenCache *cache.Cache
}

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

func NewServer(oryUri string, tinkGwUri string, ledgerUri string) *CoreServer {
	log := pkg.CreateLogger()
	return &CoreServer{
		AuthClient: auth.NewClient(oryUri),
		Logger:     log.Named("CORE"),
		TinkGwUri:  tinkGwUri,
		LedgerUri:  ledgerUri,
		//UserTokenCache: createCache(),
	}
}
