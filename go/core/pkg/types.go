package pkg

import (
	ory "github.com/ory/kratos-client-go"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionId struct {
	SessionId string `json:"sessionId"`
}

type SyncLedgerRequest struct {
	Username  string `json:"username"`
	SessionId string `json:"sessionId"`
}

type PhiSessionRequest struct {
	Token string `json:"token"`
	// mb we don't need to send whole session
	ory.Session `json:"session"`
}

type PhiClientIdResponse struct {
	TinkId string `json:"tink_id"`
}
