package handlers

import (
	"github.com/clstb/phi/go/core/auth"
	"github.com/eko/gocache/v2/cache"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CoreServer struct {
	AuthClient     *auth.AuthClient
	Logger         *zap.SugaredLogger
	UserTokenCache *cache.Cache
}

type LinkRequest struct {
	Test bool `json:"test"`
}

type AccessCodeRequest struct {
	AccessCode string `json:"access_code"`
}

type SyncRequest struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}
