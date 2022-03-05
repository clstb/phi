package client

import ory "github.com/ory/kratos-client-go"

type Session struct {
	ory.Session
	Token string
}
