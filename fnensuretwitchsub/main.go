package fnensuretwitchsub

import (
	"fmt"
	"net/http"
	"os"

	"twitch-announcer/twitch"
	"twitch-announcer/twitch/helix"
)

func respondError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func respondSuccess(w http.ResponseWriter, data []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func Handler(rw http.ResponseWriter, req *http.Request) {
	callbackURL := os.Getenv("TW_WEBHOOK_CALLBACK_URL")
	twClientID := os.Getenv("TW_CLIENT_ID")
	twClientSecret := os.Getenv("TW_CLIENT_SECRET")
	twChannelID := os.Getenv("TW_CHANNEL_ID")
	twWebhookSecret := os.Getenv("TW_WEBHOOK_SECRET")

	accessToken, err := twitch.Authenticate(req.Context(), twClientID, twClientSecret)
	if err != nil {
		respondError(rw, fmt.Errorf("authenticate twitch app: %w", err))
		return
	}

	twAPI := helix.New(twClientID)

	// Get a list of all eventsub subscriptions.
	subs, err := twAPI.GetEventsubSubs(req.Context(), accessToken)
	if err != nil {
		respondError(rw, fmt.Errorf("get subscriptions: %v", err))
		return
	}

	// Try to find out if the current server is already receiving events from Twitch.
	for _, sub := range subs {
		if sub.Transport.Callback == callbackURL && sub.Type == helix.EventsubSubTypeStreamOnline {
			// Eventsub subscription already exists. We should delete it.
			if err := twAPI.DeleteEventsubSub(req.Context(), accessToken, sub.ID); err != nil {
				respondError(rw, fmt.Errorf("delete subscription: %v", err))
				return
			}
			break
		}
	}

	// Create all required subscriptions.
	if _, err := twAPI.CreateEventsubSub(
		req.Context(),
		accessToken,
		helix.EventsubSubTypeStreamOnline,
		1,
		helix.EventsubSubCondition{BroadcasterUserID: twChannelID},
		helix.EventsubSubWebhookTransport{
			Method:   "webhook",
			Callback: callbackURL,
			Secret:   twWebhookSecret,
		},
	); err != nil {
		respondError(rw, fmt.Errorf("create subscription: %v", err))
		return
	}
	respondSuccess(rw, nil)
}
