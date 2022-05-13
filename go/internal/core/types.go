package core

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
	Token       string `json:"token"`
	ory.Session `json:"session"`
}
