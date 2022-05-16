package server

import (
	"github.com/clstb/phi/core/internal/auth"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CoreServer struct {
	AuthClient *auth.Client
	Logger     *zap.SugaredLogger
	//UserTokenCache *cache.Cache
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
