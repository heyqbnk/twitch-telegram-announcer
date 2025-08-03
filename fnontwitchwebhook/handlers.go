package fnontwitchwebhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"twitch-announcer/twitch"
	"twitch-announcer/twitch/helix"
)

func escapeMarkdown(s string) string {
	return regexp.MustCompile("/([*_`~\\\\])/g").ReplaceAllString(s, "\\$1")
}

func handleVerificationMsg(bodyBytes []byte) (challenge string, err error) {
	// Represents a subscription event which contains the information about the
	// subscription required to be confirmed.
	// Source: https://dev.twitch.tv/docs/eventsub/handling-webhook-events/#responding-to-a-challenge-request
	var body struct {
		Challenge string `json:"challenge"`
	}
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return "", err
	}

	return body.Challenge, nil
}

func handleNotification(
	ctx context.Context,
	bodyBytes []byte,
	twClientID, twClientSecret, tgBotToken string,
	tgChatID int64,
	client *http.Client,
) error {
	// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#stream-subscriptions
	var body struct {
		Subscription struct {
			Type string `json:"type"`
		} `json:"subscription"`
		Event []byte `json:"event"`
	}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	// At the moment, we handle only "Stream is Online" events.
	if body.Subscription.Type != "stream.online" {
		return nil
	}

	var event struct {
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
	}
	if err := json.Unmarshal(body.Event, &event); err != nil {
		return fmt.Errorf("unmarshal event: %w", err)
	}

	accessToken, err := twitch.Authenticate(ctx, twClientID, twClientSecret)
	if err != nil {
		return fmt.Errorf("authenticate twitch app: %v", err)
	}

	stream, err := helix.New(twClientID).GetStream(ctx, accessToken, event.BroadcasterUserLogin)
	if err != nil {
		return fmt.Errorf("get stream from Twitch: %w", err)
	}

	shortStreamLink := fmt.Sprintf("twitch.tv/%s", event.BroadcasterUserLogin)
	// Reference: https://core.telegram.org/bots/api#sendmessage
	tgMessageBytes, _ := json.Marshal(
		map[string]any{
			"chat_id": tgChatID,
			"text": fmt.Sprintf(
				"Трансляция запущена!\n\n*%s*\n\n— %s\n— %s",
				escapeMarkdown(stream.Title),
				escapeMarkdown(stream.GameName),
				shortStreamLink,
			),
			// We disable the notification as long as when we start broadcasting in Telegram, it
			// generates a corresponding service message in the channel, so we don't want to have
			// several notifications about the same thing.
			"disable_notification": true,
			"parse_mode":           "Markdown",
			"link_preview_options": map[string]any{
				"url":                fmt.Sprintf("https://%s", shortStreamLink),
				"prefer_small_media": true,
			},
		},
	)
	if resp, err := client.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tgBotToken),
		"application/json",
		strings.NewReader(string(tgMessageBytes)),
	); err != nil {
		return fmt.Errorf("send stream started message: %v", err)
	} else {
		d, err := io.ReadAll(resp.Body)
		fmt.Println(string(d), err)
	}
	return nil
}
