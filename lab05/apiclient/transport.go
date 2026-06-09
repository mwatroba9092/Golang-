package apiclient

import "net/http"

type UserAgentTransport struct {
	Base      http.RoundTripper
	UserAgent string
}

// RoundTrip wykonuje pojedynczą transakcję HTTP, dodając nagłówek.
func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	reqClone := req.Clone(req.Context())
	
	reqClone.Header.Set("User-Agent", t.UserAgent)

	return base.RoundTrip(reqClone)
}