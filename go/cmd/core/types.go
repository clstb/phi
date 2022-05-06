package main

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	SessionId string `json:"sessionId"`
}

type SyncLedgerRequest struct {
	Username  string `json:"username"`
	SessionId string `json:"sessionId"`
}
