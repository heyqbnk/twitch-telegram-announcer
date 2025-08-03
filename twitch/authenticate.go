package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Authenticate(ctx context.Context, clientID, clientSecret string) (string, error) {
	// Create request.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://id.twitch.tv/oauth2/token",
		strings.NewReader(
			url.Values{
				"client_id":     {clientID},
				"client_secret": {clientSecret},
				"grant_type":    {"client_credentials"},
			}.Encode(),
		),
	)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send request.
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", fmt.Errorf("send http request: %v", err)
	}

	// Read response body.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response. Status %d: %v", res.StatusCode, err)
	}

	// We received non-successful status code. We should try to extract error
	// sent from Twitch.
	if res.StatusCode != 200 {
		// Try to unmarshal it to expected structure.
		var response struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", fmt.Errorf("unexpected response: %v: %s", err, string(body))
		}
		return "", fmt.Errorf("request failed (status %d): %s", res.StatusCode, response.Message)
	}

	var authResult struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	// Try to unmarshal body to destination.
	if err := json.Unmarshal(body, &authResult); err != nil {
		return "", fmt.Errorf("invalid response: %v", err)
	}

	return authResult.AccessToken, nil
}
