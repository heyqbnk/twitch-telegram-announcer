package helix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (a *API) request(
	ctx context.Context,
	accessToken string,
	method, path string,
	query url.Values,
	body interface{},
	dest interface{},
) error {
	var bodyBytes []byte

	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidBody, err)
		}

		bodyBytes = bodyJson
	}

	// Create request.
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("https://api.twitch.tv/helix/%s?%s", path, query.Encode()),
		nil,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Add("Client-Id", a.clientID)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	if len(bodyBytes) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.Header.Add("Content-Type", "application/json")
	}

	// Send the request.
	res, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("send http request: %v", err)
	}

	if res.StatusCode == 204 {
		if dest != nil {
			return fmt.Errorf("%w: received 204 HTTP status code, but dest specified", ErrInvalidResponse)
		}

		return nil
	}

	// Read response body.
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response. Status %d: %v", res.StatusCode, err)
	}

	// Try to unmarshal it to expected structure.
	var response struct {
		Error   string        `json:"error"`
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Data    requestResult `json:"data"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return fmt.Errorf("%w: %v", ErrUnexpectedResponse, err)
	}

	if response.Status != 0 {
		var err error

		switch response.Status {
		case 400:
			err = Err400
		case 401:
			err = ErrNotAuthorized
		case 404:
			err = Err404
		default:
			err = ErrUnknown
		}
		return fmt.Errorf("%w: %s", err, response.Message)
	}

	// In case, destination was specified, it means, caller expects some
	// data to be returned by request.
	if dest != nil {
		// Try to unmarshal body to destination.
		if err := json.Unmarshal(response.Data.bytes, dest); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidResponse, err)
		}
	}

	return nil
}

// Miscellaneous type to appropriately handle response from API. This
// type just preserves bytes passed during json unmarshalling which then
// could be used to unmarshal into another structure.
type requestResult struct {
	bytes []byte
}

func (r *requestResult) UnmarshalJSON(b []byte) error {
	r.bytes = b
	return nil
}
