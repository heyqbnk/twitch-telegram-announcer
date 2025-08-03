package helix

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Source: https://dev.twitch.tv/docs/api/reference/#get-streams

type Stream struct {
	GameName string `json:"game_name"`
	Title    string `json:"title"`
}

func (a *API) GetStream(ctx context.Context, accessToken, login string) (Stream, error) {
	var res []Stream
	if err := a.request(
		ctx,
		accessToken,
		http.MethodGet,
		"streams",
		url.Values{"user_login": {login}},
		nil,
		&res,
	); err != nil {
		return Stream{}, fmt.Errorf("request error: %w", err)
	}

	if len(res) == 0 {
		return Stream{}, errors.New("stream not found")
	}

	return res[0], nil
}
