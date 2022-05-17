package server

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
