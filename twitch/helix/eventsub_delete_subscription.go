package helix

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (a *API) DeleteEventsubSub(ctx context.Context, accessToken, id string) error {
	if err := a.request(
		ctx,
		accessToken,
		http.MethodDelete,
		"eventsub/subscriptions",
		url.Values{"id": {id}},
		nil,
		nil,
	); err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	return nil
}
