package helix

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (a *API) GetEventsubSubs(ctx context.Context, accessToken string) ([]EventsubSub, error) {
	var res []EventsubSub
	if err := a.request(
		ctx,
		accessToken,
		http.MethodGet,
		"eventsub/subscriptions",
		url.Values{},
		nil,
		&res,
	); err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return res, nil
}
