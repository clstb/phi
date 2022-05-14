package client

import (
	"context"
	ory "github.com/ory/kratos-client-go"
)

func (c *AuthClient) UpdateTraits(sess *Session, traits map[string]interface{}) (*Session, error) {

	identity, _, err := c.OryClient.V0alpha2Api.AdminUpdateIdentity(context.Background(), sess.Identity.Id).AdminUpdateIdentityBody(ory.AdminUpdateIdentityBody{
		State:  *sess.Identity.State,
		Traits: traits,
	}).Execute()
	if err != nil {
		return nil, err
	}
	sess.Identity = *identity
	return sess, nil
}
