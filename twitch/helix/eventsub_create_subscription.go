package helix

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (a *API) CreateEventsubSub(
	ctx context.Context,
	accessToken string,
	subType EventsubSubType,
	version int,
	condition EventsubSubCondition,
	transport EventsubSubWebhookTransport,
) ([]EventsubSub, error) {
	var res []EventsubSub
	if err := a.request(
		ctx,
		accessToken,
		http.MethodPost,
		"eventsub/subscriptions",
		url.Values{},
		struct {
			Type      EventsubSubType             `json:"type"`
			Version   int                         `json:"version"`
			Condition EventsubSubCondition        `json:"condition"`
			Transport EventsubSubWebhookTransport `json:"transport"`
		}{
			Type:      subType,
			Version:   version,
			Condition: condition,
			Transport: transport,
		},
		&res,
	); err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return res, nil
}
