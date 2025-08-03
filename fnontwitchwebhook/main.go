package fnontwitchwebhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func getHmac(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyTwitchMessageSignature(
	messageID, messageTimestamp, body, signature, secret string,
) bool {
	message := messageID + messageTimestamp + body
	messageSignature := "sha256=" + getHmac(message, secret)
	return messageSignature == signature
}

func respondError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func respondSuccess(w http.ResponseWriter, data []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func Handler(rw http.ResponseWriter, req *http.Request) {
	tgChatIDStr := os.Getenv("TG_CHAT_ID")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	twWebhookSecret := os.Getenv("TW_WEBHOOK_SECRET")
	twClientID := os.Getenv("TW_CLIENT_ID")
	twClientSecret := os.Getenv("TW_CLIENT_SECRET")

	tgChatID, err := strconv.ParseInt(tgChatIDStr, 10, 64)
	if err != nil {
		respondError(rw, fmt.Errorf("chat ID invalid: %v", tgChatIDStr))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondError(rw, err)
		return
	}

	if isValid := verifyTwitchMessageSignature(
		req.Header.Get("twitch-eventsub-message-id"),
		req.Header.Get("twitch-eventsub-message-timestamp"),
		string(body),
		req.Header.Get("twitch-eventsub-message-signature"),
		twWebhookSecret,
	); !isValid {
		respondError(rw, errors.New("signature invalid"))
		return
	}

	switch req.Header.Get("twitch-eventsub-message-type") {
	// In case, we received "webhook_callback_verification" message type,
	// we just respond with the HTTP 200 status code with the challenge string
	// received from the Twitch server to make the subscription work.
	case "webhook_callback_verification":
		if challenge, err := handleVerificationMsg(body); err != nil {
			respondError(rw, fmt.Errorf("handle verification message: %w", err))
		} else {
			respondSuccess(rw, []byte(challenge))
		}
		return
	case "notification":
		if err := handleNotification(
			req.Context(),
			body,
			twClientID,
			twClientSecret,
			tgBotToken,
			tgChatID,
			&http.Client{},
		); err != nil {
			respondError(rw, fmt.Errorf("handle notification: %w", err))
		} else {
			respondSuccess(rw, nil)
		}
	default:
		respondSuccess(rw, nil)
	}
}
