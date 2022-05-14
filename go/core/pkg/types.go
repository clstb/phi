package pkg

type SessionId struct {
	SessionId string `json:"sessionId"`
}

type SyncLedgerRequest struct {
	Username  string `json:"username"`
	SessionId string `json:"sessionId"`
}
