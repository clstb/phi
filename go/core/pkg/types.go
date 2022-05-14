package pkg

import "github.com/clstb/phi/go/core/pkg/client"

type SessionId struct {
	SessionId string `json:"sessionId"`
}

type SyncLedgerRequest struct {
	Username  string `json:"username"`
	SessionId string `json:"sessionId"`
}

type CoreServer struct {
	authClient *client.AuthClient
}
