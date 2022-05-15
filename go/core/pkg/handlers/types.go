package handlers

import (
	"github.com/clstb/phi/go/core/pkg/auth"
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

type SessionId struct {
	SessionId string `json:"sessionId"`
}
