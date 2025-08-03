package helix

import "net/http"

// API is wrapper for Twitch API v4.
type API struct {
	client   *http.Client
	clientID string
}

func New(clientID string) *API {
	return &API{
		client:   &http.Client{},
		clientID: clientID,
	}
}
