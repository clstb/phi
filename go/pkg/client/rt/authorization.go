package rt

import "net/http"

type AuthorizationRoundTripper struct {
	Token string
	Next  http.RoundTripper
}

func (rt AuthorizationRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+rt.Token)
	return rt.Next.RoundTrip(req)
}
